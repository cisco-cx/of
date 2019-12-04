#!/usr/bin/env python3
"""Convert list of APIC faults into YAML configuration for am-apic-client-go

This script scrapes the DevNet documentation webpage for APIC faults
and converts them into the YAML format expected by am-apic-client-go.
This YAML is output on STDOUT.

Usage:

    ./bin/apic_faults_to_alerts_yaml.py > am-apic-client-go/alerts.yaml
"""

import requests
import sys
import argparse
import csv
from bs4 import BeautifulSoup
from yaml import load, dump

try:
    from yaml import CLoader as Loader, CDumper as Dumper
except ImportError:
    from yaml import Loader, Dumper

parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
parser.add_argument('--csv', action='store_true', help='Output CSV instead of YAML. The CSV will in the format used by InnoEye OSS for their Alert Dictionary')
args = parser.parse_args()

# URL with list of all APIC faults and their details
APIC_FAULT_URL = "https://pubhub.devnetcloud.com/media/apic-mim-ref-411/docs/FaultMessages.html"

APIC_EMPTY_EXPLANATION = "None set."

# List of faults that we are always going to ignore
DROPPED_FAULTS = [
    'F0021',
    'F0023',
    'F0132',
    'F0413',
    'F0454',
    'F0467',
    'F0475',
    'F0603',
    'F0699',
    'F0756',
    'F0843',
    'F0844',
    'F0845',
    'F0846',
    'F0847',
    'F0848',
    'F0849',
    'F1199',
    'F1228',
    'F1296',
    'F1298',
    'F1299',
    'F1300',
    'F1313',
    'F1368',
    'F1371',
    'F1425',
    'F1432',
    'F1449',
    'F1471',
    'F1483',
    'F1545',
    'F1546',
    'F1547',
    'F1548',
    'F1549',
    'F1550',
    'F1551',
    'F1563',
    'F1564',
    'F1572',
    'F1573',
    'F1574',
    'F2168',
    'F2194',
    'F2543',
    'F2547',
    'F2773',
    'F2840',
    'F2844',
    'F2967',
    'F3019',
    'F3057',
    'F3062',
    'F100264',
    'F110473',
    'F112425',
    'F112425',
    'F119936',
    'F119936',
    'F606391',
    'F608054',
]

def apic_severity_to_syslog(severity):
    """Convert APIC severity to SYSLOG levels"""
    if severity == "critical":
        return "critical"
    elif severity == "major":
        return "error"
    elif severity == "minor":
        return "warning"
    elif severity == "warning":
        return "warning"
    elif severity == "info":
        return "info"
    else:
        raise Exception("Unknown APIC severity {}".format(severity))

def apic_severity_to_syslog_code(severity):
    """Convert APIC severity to SYSLOG levels"""
    if severity == "critical":
        return 2
    elif severity == "major":
        return 3
    elif severity == "minor":
        return 4
    elif severity == "warning":
        return 5
    elif severity == "info":
        return 6
    else:
        raise Exception("Unknown APIC severity {}".format(severity))

def apic_severity_to_oss_severity(severity):
    """Convert APIC severity to OSS Severity"""
    if severity == "critical":
        return "Critical"
    elif severity == "major":
        return "Major"
    elif severity == "minor":
        return "Minor"
    elif severity == "warning":
        return "Warning"
    elif severity == "info":
        return "Info"
    else:
        raise Exception("Unknown APIC severity {}".format(severity))

def apic_severity_to_oss_classification(severity):
    """Convert APIC severity to OSS Classification"""
    if severity == "critical":
        return "Outage"
    elif severity == "major":
        return "Deterioration"
    elif severity == "minor":
        return "Deterioration"
    elif severity == "warning":
        return "Notification"
    elif severity == "info":
        return "Normal"
    else:
        raise Exception("Unknown APIC severity {}".format(severity))

def apic_severity_to_oss_service_affected(severity):
    """Convert APIC severity to OSS Service Affected field"""
    if severity == "critical":
        return "Yes"
    elif severity == "major":
        return "Yes"
    elif severity == "minor":
        return "No"
    elif severity == "warning":
        return "No"
    elif severity == "info":
        return "No"
    else:
        raise Exception("Unknown APIC severity {}".format(severity))

def apic_severity_to_oss_incident_creation(severity):
    """Convert APIC severity to OSS Incident Creationfield"""
    if severity == "critical":
        return "Yes"
    elif severity == "major":
        return "Yes"
    elif severity == "minor":
        return "Yes"
    elif severity == "warning":
        return "Yes"
    elif severity == "info":
        return "No"
    else:
        raise Exception("Unknown APIC severity {}".format(severity))

def most_severe_fault(fault_list):
    """Given a list of faults, return the one with the highest severity (numerically lowest based on syslog code)"""
    if len(fault_list) == 1:
        return fault_list[0]
    else:
        most_severe = None
        for fault in fault_list:
            syslog_code = apic_severity_to_syslog_code(fault["severity"])
            if not most_severe or syslog_code < apic_severity_to_syslog_code(most_severe["severity"]):
                most_severe = fault

        return most_severe

def format_yaml(faults):
    """Print the YAML formatted faults"""
    # organize all of the raw data into alerts
    alerts = { "alerts": {}, "dropped_faults": {} }
    for fault_name in faults:
        am_fault_name = apic_name_to_am(fault_name)
        # pick the fault in this list with the highest severity
        most_severe = most_severe_fault(faults[fault_name])
        fault_codes = []

        # look at each fault under this name
        for fault in faults[fault_name]:
            fault_code = fault["code"]
            # if we should drop it, drop it like it's hot, else keep it
            if should_drop_fault(fault):
                alerts["dropped_faults"][fault_code] = { "fault_name": fault_name }
            else:
                fault_codes.append(fault_code)

        if len(fault_codes):
            alerts["alerts"][am_fault_name] = {
                "faults": {},
                "alert_severity": apic_severity_to_syslog(most_severe["severity"]),
            }
            for fault_code in fault_codes:
                alerts["alerts"][am_fault_name]["faults"][fault_code] = { "fault_name": fault_name }

    config = {
        "apic": {
            "alert_severity_threshold": "major",
            "drop_unknown_alerts": True,
        },
        "defaults": {
            "alert_severity": "error",
        },
        "alerts": alerts["alerts"],
        "dropped_faults": alerts["dropped_faults"],
    }
    print("# Generated by apic_faults_to_alerts_yaml.py, DO NOT MODIFY BY HAND")
    print(dump(config, Dumper=Dumper))

def should_drop_fault(fault):
    fault_code = fault["code"]
    # drop the ones in the DROPPED_FAULTS
    if fault["code"] in DROPPED_FAULTS:
        return True

    # drop the FSM faults
    if fault["name"].startswith('fsm'):
        return True

    # drop things without a description
    description = fault_description(fault)
    if not description or description == APIC_EMPTY_EXPLANATION:
        return True

    return False

def format_csv(faults):
    """Print the CSV formatted results in the format dictated by InnoEye OSS"""
    header = ["Alarm Code", "Alarm Name", "EMS", "Classification", "Service Affected", "Category", "MO", "Default Severity", "Alarm Type", "Alarming Delay (minutes)", "Estimated time to Close(Minutes)", "Incident Creation", "TT Delay (minutes)", "WO Delay (minutes)", "Clear Name", "Clear Trap", "Software Release", "Description", "Impact", "Suggestion", "Generic Alarm Name", "Probable Cause"]
    writer = csv.DictWriter(sys.stdout, header, quoting=csv.QUOTE_ALL)
    # output the header
    writer.writeheader()

    # process each fault and add it to the CSV file if necessary
    for fault_name in faults:
        nice_fault_name = apic_name_to_am(fault_name)
        # pick the fault in this list with the highest severity
        most_severe = most_severe_fault(faults[fault_name])
        fault_code = most_severe["code"]

        if not should_drop_fault(most_severe):
            severity = most_severe["severity"]

            row = {
                "Alarm Code": nice_fault_name,
                "Alarm Name": nice_fault_name,
                "EMS": "OF",
                "Classification": apic_severity_to_oss_classification(severity),
                "Service Affected": apic_severity_to_oss_service_affected(severity),
                "Category": "-",
                "MO": "-",
                "Default Severity": apic_severity_to_oss_severity(severity),
                "Alarm Type": "-",
                "Alarming Delay (minutes)": "-",
                "Estimated time to Close(Minutes)": "-",
                "Incident Creation": apic_severity_to_oss_incident_creation(severity),
                "TT Delay (minutes)": "-",
                "WO Delay (minutes)": "-",
                "Clear Name": "-",
                "Clear Trap": "Yes",
                "Software Release": "-",
                "Description": "APIC: " + fault_description(most_severe),
                "Impact": "-",
                "Suggestion": "-",
                "Generic Alarm Name": "-",
                "Probable Cause": "-",
            }
            writer.writerow(row)

def cleanup(text):
    """Cleanup the text found in the HTML page"""
    return ' '.join(text.replace(u'\u200b', '').replace('\n', ' ').split())

def apic_name_to_am(name):
    """Convert the APIC fault name into something better for AlertManager.
    Example: fltFabricSelectorIssuesConfigFailed -> apicFltFabricSelectorIssuesConfigFailed"""
    return "apic" + name[0].upper() + name[1:]

def fault_description(fault):
    description = fault["explanation"]
    if not description or description == APIC_EMPTY_EXPLANATION:
        description = fault["message"]

    return description

def main():
    req = requests.get(APIC_FAULT_URL, timeout=30)
    if req.status_code != 200:
        sys.stderr.write("Failed to retrieve APIC Faults webpage ({}), Status Code: {}\n".format(APIC_FAULT_URL, req.status_code))
        sys.exit(1)
    soup = BeautifulSoup(req.text, 'html.parser')
    fault_codes = {}

    # the document is one large table, so get all of the rows as a starting point
    rows = soup.find_all("tr")
    current_fault = None
    for row in rows:
        first_cell = cleanup(row.td.text)

        # we have a fault code which means we are starting a new fault
        if first_cell == 'Fault Code':
            fault_code = cleanup(row.find_all('td')[1].text)
            if not fault_code in fault_codes:
                fault_codes[fault_code] = { "code": fault_code, "name": "", "type": "", "cause": "", "severity": ""}
            current_fault = fault_code

        if first_cell == 'MIB Fault Name':
            fault_codes[current_fault]["name"] = cleanup(row.find_all('td')[1].text)

        if first_cell == 'Type':
            fault_codes[current_fault]["type"] = cleanup(row.find_all('td')[1].text)

        if first_cell == 'Cause':
            fault_codes[current_fault]["cause"] = cleanup(row.find_all('td')[1].text)

        if first_cell == 'Severity':
            fault_severity = cleanup(row.find_all('td')[1].text)
            fault_codes[current_fault]["severity"] = fault_severity

        if first_cell == 'Explanation':
            fault_codes[current_fault]['explanation'] = cleanup(row.find_all('td')[1].text)

        if first_cell == 'Message':
            fault_codes[current_fault]['message'] = cleanup(row.find_all('td')[1].text)

    # ACI uses the same fault name for multiple fault codes, so we need to re-organize them to be keyed off of their name as a list
    faults = {}
    for fault_code in fault_codes:
        fault = fault_codes[fault_code]
        fault_name = fault["name"]
        if not fault_name in faults:
            faults[fault_name] = []

        faults[fault_name].append(fault)

    if args.csv:
        format_csv(faults)
    else:
        format_yaml(faults)

if __name__ == "__main__":
    main()

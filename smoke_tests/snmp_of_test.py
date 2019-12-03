import json
import datetime
import requests


def get_alerts(host, silenced=False):
    """Get an alert list querying alert manager's API."""
    url = f"{host}/api/v2/alerts"
    resp = requests.get(url, params={"silenced": silenced})
    resp.raise_for_status()
    return resp.json()


def alert_in(alertname, alerts):
    """Check if the alertname is part of the alerts."""
    for alert in alerts:
        if alertname in alert["labels"].values():
            return True
    return False


def post_alert(host, data):
    """Post an alert via alert manager's API."""
    resp = requests.post(f"{host}/api/v2/events", json=data)
    resp.raise_for_status()


def silence_all_alerts(host, alertname, duration_s=60):
    """Silence all alerts `alertname` fired, during `duration_s`."""
    starts_at = datetime.datetime.utcnow()
    ends_at = starts_at + datetime.timedelta(seconds=duration_s)
    data = {
        "startsAt": starts_at.isoformat() + "Z",
        "endsAt": ends_at.isoformat() + "Z",
        "createdBy": "of tester",
        "comment": "Testing of",
    }
    for alert in get_alerts(host):
        if alert["labels"]["alertname"] != alertname:
            continue
        matchers = []
        for name, value in alert["labels"].items():
            matchers.append({"name": name, "value": value, "isRegex": False, },)
        data["matchers"] = matchers
        resp = requests.post(f"{host}/api/v2/silences", json=data)
        resp.raise_for_status()


def test_snmptrap_of_to_am_alerts(of_url, alertmanager_url, snmp_entry):
    """
    Test if the SNMPTraps taken from elasticsearch generate their corresponding alerts
    in AM.
    """
    alertname = snmp_entry["alertname"]
    silence_all_alerts(alertmanager_url, alertname)

    alerts = get_alerts(alertmanager_url)
    assert not alert_in(
        alertname, alerts
    ), f'{snmp_entry["name"]}.{alertname} found in alerts'

    data = [json.loads(snmp_entry["firing_elasticsearch_json"])["_source"]]
    post_alert(of_url, data)

    alerts = get_alerts(alertmanager_url)
    assert alert_in(
        alertname, alerts
    ), f'{snmp_entry["name"]}.{alertname} not found in alerts'
    # ToDo how does it work the clearing?

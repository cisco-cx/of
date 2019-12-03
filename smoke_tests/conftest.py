import glob
import time
import yaml
import requests


def pytest_addoption(parser):
    """Set fixture parameters from the command line"""
    parser.addoption(
        "--of-url",
        action="append",
        required=True,
        help="of url e.g. http://localhost:8080",
    )
    parser.addoption(
        "--alertmanager-url",
        action="append",
        required=True,
        help="alert manager url e.g. http://localhost:9093",
    )
    parser.addoption(
        "--wait-for-services",
        action="store_true",
        help="Wait for of and alertmanager to be ready before start",
    )


def pytest_generate_tests(metafunc):
    """
    Generate tests from the command parameters + the yaml fixtures placed in the
    fixtures folder.
    """
    if "of_url" in metafunc.fixturenames:
        metafunc.parametrize("of_url", metafunc.config.getoption("of_url"))
    if "alertmanager_url" in metafunc.fixturenames:
        metafunc.parametrize(
            "alertmanager_url", metafunc.config.getoption("alertmanager_url")
        )
    if "snmp_entry" in metafunc.fixturenames:
        cfg_entries = []
        for filepath in glob.glob("fixtures/*.yaml"):
            with open(filepath) as fh:
                cfg = yaml.load(fh, Loader=yaml.Loader)
                for device_type, snmp_entries in cfg.items():
                    for snmp_entry in snmp_entries:
                        snmp_entry["name"] = device_type
                    cfg_entries.extend(snmp_entries)
        metafunc.parametrize("snmp_entry", cfg_entries)


def request_ignore_connect_error(url):
    """Do request ignoring the connection error."""
    resp = None
    try:
        resp = requests.get(url)
    except requests.exceptions.ConnectionError:
        pass
    return resp


def pytest_configure(config):
    """
    Check if wait for services is true and check for services be ready before start
    executing the tests.
    """
    if config.getoption("wait_for_services"):
        of_ready = False
        alertmanager_ready = False
        of_url = f'{config.getoption("of_url")[0]}/api/v2/status'
        alertmanager_url = f'{config.getoption("alertmanager_url")[0]}/api/v2/status'
        print("Waiting for services (of and alertmanager) be ready", flush=True)
        while not (of_ready and alertmanager_ready):
            if not of_ready:
                resp = request_ignore_connect_error(of_url)
                if resp and "success" in resp.text:
                    of_ready = True
                    print("of is ready", flush=True)
            if not alertmanager_ready:
                resp = request_ignore_connect_error(alertmanager_url)
                if resp and resp.json().get("cluster", {}).get("status") == "ready":
                    alertmanager_ready = True
                    print("alertmanager is ready", flush=True)
            time.sleep(2)

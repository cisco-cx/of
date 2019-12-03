# Smoke tests

It performs an end to end testing of the `of`. `snmp` is the only handler supported for now.

The main idea is to apply a known input into the `of`, it can be obtained from elasticsearch or they can be built artificially, and then check if the alert was generated correctly in the `alertmanager` through its API.

There is a single test with different inputs defined under the folder fixtures

```python3
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
  ```

It sends to the `of` the snmp entry obtained from fixtures, after that it checks the alert was fired in the `alertmanager`. For defining an entry in the fixtures you just can get a dump using the tool `elasticsearch-dump` or using `kibana`, e.g.

```yaml
epc:
  - alertname: starCLISessionStart
    firing_elasticsearch_json: '{"_index":"snmptrapd-2019.11.20","_type":"doc","_id":"UYcuhm4B8tU2gtyiRIs0","_score":1,"_source":{"agent":{"type":"filebeat","version":"7.2.0","ephemeral_id":"63b833f9-3135-4de7-a544-fceeb5f0a50b","id":"a66e743a-0e1a-4faf-800c-b2cbb73bbf1b","hostname":"uhn4ttcsof101003"},"input":{"type":"stdin"},"@version":"1","@timestamp":"2019-11-20T00:20:02.489Z","host":{"name":"uhn4ttcsof101003"},"ecs":{"version":"1.0.0"},"tags":["beats_input_codec_plain_applied"],"message":"SNMPTRAP timestamp=[2019-11-20T00:20:02Z] address=[UDP/IPv6: [240b:c010:101:5483:684:2:0:101]:36693] pdu_security=[TRAP2, SNMP v3, user epc-snmp-user, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (426256873) 49 days, 8:02:48.73\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.4.1.8164.2.52\t.1.3.6.1.4.1.8164.1.19.1.1.2 = STRING: ;tool-of;\t.1.3.6.1.4.1.8164.1.19.1.1.4 = STRING: ;Security Administrator;\t.1.3.6.1.4.1.8164.1.19.1.1.3 = STRING: ;/dev/pts/2;\t.1.3.6.1.6.3.1.1.4.3.0 = OID: .1.3.6.1.4.1.8164.2]","document":{"receipts":{"filebeat":{"ecs":{"version":"1.0.0"},"agent":{"type":"filebeat","version":"7.2.0","ephemeral_id":"63b833f9-3135-4de7-a544-fceeb5f0a50b","id":"a66e743a-0e1a-4faf-800c-b2cbb73bbf1b","hostname":"uhn4ttcsof101003"},"input":{"type":"stdin"},"message":"SNMPTRAP timestamp=[2019-11-20T00:20:02Z] address=[UDP/IPv6: [240b:c010:101:5483:684:2:0:101]:36693] pdu_security=[TRAP2, SNMP v3, user epc-snmp-user, context ] vars[.1.3.6.1.2.1.1.3.0 = Timeticks: (426256873) 49 days, 8:02:48.73\t.1.3.6.1.6.3.1.1.4.1.0 = OID: .1.3.6.1.4.1.8164.2.52\t.1.3.6.1.4.1.8164.1.19.1.1.2 = STRING: ;tool-of;\t.1.3.6.1.4.1.8164.1.19.1.1.4 = STRING: ;Security Administrator;\t.1.3.6.1.4.1.8164.1.19.1.1.3 = STRING: ;/dev/pts/2;\t.1.3.6.1.6.3.1.1.4.3.0 = OID: .1.3.6.1.4.1.8164.2]","@version":"1","@timestamp":"2019-11-20T00:20:02.489Z","log":{"file":{"path":""},"offset":0},"host":{"name":"uhn4ttcsof101003"}},"snmptrapd":{"pduSecurity":"TRAP2, SNMP v3, user epc-snmp-user, context","vars":[{"oid":".1.3.6.1.2.1.1.3.0","type":"Timeticks","value":"(426256873) 49 days, 8:02:48.73"},{"oid":".1.3.6.1.6.3.1.1.4.1.0","type":"OID","value":".1.3.6.1.4.1.8164.2.52"},{"oid":".1.3.6.1.4.1.8164.1.19.1.1.2","type":"STRING","value":";tool-of;"},{"oid":".1.3.6.1.4.1.8164.1.19.1.1.4","type":"STRING","value":";Security Administrator;"},{"oid":".1.3.6.1.4.1.8164.1.19.1.1.3","type":"STRING","value":";/dev/pts/2;"},{"oid":".1.3.6.1.6.3.1.1.4.3.0","type":"OID","value":".1.3.6.1.4.1.8164.2"}],"timestamp":"2019-11-20T00:20:02Z","source":{"transportLayerProtocol":"UDP","address":"240b:c010:101:5483:684:2:0:101","port":"36693","hostname":"UHN3tte2upda0004.rmn.local","internetLayerProtocol":"IPv6"}},"logstash":{"tags":["beats_input_codec_plain_applied"]}},"kind":"SNMPTrap","apiVersion":"v1alpha1"},"log":{"file":{"path":""},"offset":0}}}'
```

You can run it as following
```sh
docker-compose -f smoke_tests/docker-compose.yaml up --build --abort-on-container-exit --exit-code-from test
```

Or if you already has alertmanager and of running, you can execute
```sh
pytest-3 smoke_tests --of-url http://localhost:8942 --alertmanager-url http://localhost:9093 --wait-for-services
```

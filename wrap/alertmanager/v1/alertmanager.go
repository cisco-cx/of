package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	http "github.com/cisco-cx/of/wrap/http/v1"
)

type AlertService struct {
	Version string
	AmURL   string
}

// Send alerts to Alertmanager
func (a *AlertService) Notify(alerts []*Alert) error {
	amAlertsURL := fmt.Sprintf("%s/api/v1/alerts", a.AmURL)

	b := []byte{}
	b, err := json.Marshal(alerts)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", amAlertsURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", a.Version)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.NewClient().Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		err := fmt.Errorf("POST to AlertManager on %q returned HTTP %d:  %s", amAlertsURL, resp.StatusCode, body)
		return err
	}

	return nil
}

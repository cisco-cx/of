package alertmanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	of "github.com/cisco-cx/of/lib/v1alpha1"
	http "github.com/cisco-cx/of/wrap/http/v1alpha1"
)

type AlertService struct {
	*of.ACIConfig
}

// Send alerts to Alertmanager
func (a *AlertService) Notify(alerts []*of.Alert) error {
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

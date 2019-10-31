package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	of "github.com/cisco-cx/of/pkg/v2"
	http "github.com/cisco-cx/of/wrap/http/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
)

type AlertService struct {
	Version   string
	AmURL     string
	Throttle  bool
	PostTime  int
	SleepTime int
	SendTime  int
	Log       *logger.Logger
}

// Send alerts to Alertmanager
func (a *AlertService) notify(alerts []of.Alert) error {
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
		a.Log.WithError(err).Errorf("")
		return err
	}

	return nil
}

// Divide alerts into smaller chunks and spread posting to Alertmanager over a.SendTime milliseconds.
func (a *AlertService) Notify(alerts *[]of.Alert) error {

	totalCount := len(*alerts)
	// Send all alerts in a single post to Alertmanager, if Throttle is disabled
	// or a.SendTime is less than time needed for a single post.
	if a.Throttle == false || a.SendTime <= a.PostTime+a.SleepTime {
		return a.notify((*alerts)[0:totalCount])
	}

	// Max num. of requests that can be send in a.SendTime.
	maxRequests := a.SendTime / (a.PostTime + a.SleepTime)
	start := 0
	if totalCount > maxRequests {
		chunkSize := totalCount / maxRequests

		end := chunkSize
		for end <= totalCount {
			err := a.notify((*alerts)[start:end])
			if err != nil {
				a.Log.WithError(err).Errorf("Failed to send alerts chunk[%d:%d]", start, end)
			}
			start = end
			end = start + chunkSize
			time.Sleep(time.Duration(a.SleepTime) * time.Millisecond)
		}
	}

	// Handle condition where totalCount is not divisible by maxRequests.
	if start < totalCount {
		err := a.notify((*alerts)[start:totalCount])
		if err != nil {
			a.Log.WithError(err).Errorf("Failed to send alerts chunk[%d:%d]", start, totalCount)
		}
		time.Sleep(time.Duration(a.SleepTime) * time.Millisecond)
	}
	return nil
}

package v1

import (
	"fmt"
	"reflect"

	"github.com/prometheus/common/model"
	of "github.com/cisco-cx/of/pkg/v1"
	strcase "github.com/cisco-cx/of/wrap/strcase/v1"
)

// Label representing alert fingerprint.
const amAlertFingerprintLabel = "alert_fingerprint"

type Alert of.Alert

// Initiate new Alertmanager alert.
func NewAlert(f of.ACIFaultRaw) *Alert {
	a := &Alert{
		Annotations: annotations(f),
		Labels:      of.LabelMap{},
	}
	return a
}

// Fingerprint alert.
func (a *Alert) Fingerprint() string {
	//a.Labels[amAlertFingerprintLabel] = model.LabelValue(a.Labels.Fingerprint().String())
	ls := make(model.LabelSet)
	for k, v := range a.Labels {
		ls[model.LabelName(k)] = model.LabelValue(v)
	}
	return ls.Fingerprint().String()
}

// Convert all alert fields to annotations.
func annotations(f of.ACIFaultRaw) map[string]string {
	// refs:
	// * https://stackoverflow.com/a/18927729
	// * https://play.golang.org/p/_zSICvw562P

	v := reflect.ValueOf(f)

	annotations := make(map[string]string, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		c := strcase.CaseString(v.Type().Field(i).Name)
		snakeCaseOldKey := c.ToSnake()
		key := fmt.Sprintf("fault_%s", snakeCaseOldKey)
		value := v.Field(i).String()
		annotations[key] = value
	}

	return annotations
}

package v1alpha1

import (
	"io"

	"gopkg.in/yaml.v2"
	"github.com/cisco-cx/of/lib/v1alpha1"
)

type Alerts v1alpha1.Alerts

func (a *Alerts) Load(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(a)
}

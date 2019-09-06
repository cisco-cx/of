package v1alpha1

import (
	"io"

	"gopkg.in/yaml.v2"
	"github.com/cisco-cx/of/lib/v1alpha1"
)

type Secrets v1alpha1.Secrets

func (s *Secrets) Load(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(s)
}

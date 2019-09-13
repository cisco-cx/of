package v1alpha1

import (
	of "github.com/cisco-cx/of/lib/v1alpha1"
)

type ACIService struct {
	*of.ACIConfig
}

func (s *ACIService) Faults() ([]of.Map, error) {
	client, err := NewACIClient(of.ACIClientConfig{Hosts: []string{s.SourceAddress}, User: s.User, Pass: s.Pass})
	if err != nil {
		return nil, err
	}

	err = client.Login()
	if err != nil {
		return nil, err
	}

	defer func() {
		err := client.Logout()
		if err != nil {
		}
	}()

	faults, err := client.Faults()
	if err != nil {
		return nil, err
	}

	return faults, nil
}

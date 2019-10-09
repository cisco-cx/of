package v1

import (
	of "github.com/cisco-cx/of/pkg/v1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
)

type ACIService struct {
	*of.ACIConfig
	*logger.Logger
}

func (s *ACIService) Faults() ([]of.Map, error) {
	client, err := NewACIClient(of.ACIClientConfig{Hosts: []string{s.SourceHostname},
		User: s.User, Pass: s.Pass}, s.Logger)
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

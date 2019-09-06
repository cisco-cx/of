package v1alpha1_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/cisco-cx/of/lib/v1alpha1"
	configv1alpha1 "github.com/cisco-cx/of/wrap/config/v1alpha1"
)

func TestSecretsLoader(t *testing.T) {

	r := strings.NewReader(`
apic:
  cluster:
    name: lab-aci`)

	expected := v1alpha1.Secrets{
		APIC: v1alpha1.SecretsConfigAPIC{
			Cluster: v1alpha1.SecretsConfigAPICCluster{
				Name: "lab-aci",
			},
		},
	}

	cfg := configv1alpha1.Secrets{}
	cfg.Load(r)
	require.EqualValues(t, cfg, expected)
}

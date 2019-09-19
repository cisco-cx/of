package v1alpha1

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	info "github.com/prometheus/common/version"
	of "github.com/cisco-cx/of/lib/v1alpha1"
)

type Version struct {
	program string
	gv      *prometheus.GaugeVec
}

// NewCollector returns a collector which exports metrics about current version information.
func NewCollector(program string) of.Versioner {
	initCollector()
	v := &Version{program: program}
	v.gv = info.NewCollector(program)
	return v
}

func (v *Version) Print() string {
	initCollector()
	return info.Print(v.program)
}

func (v *Version) BuildContext() string {
	initCollector()
	return info.BuildContext()
}

func (v *Version) Info() string {
	initCollector()
	return info.Info()
}

func initCollector() {
	info.Version = of.Version
	info.Revision = of.Revision
	info.Branch = of.Branch
	info.BuildUser = of.BuildUser
	info.BuildDate = of.BuildDate
	info.GoVersion = runtime.Version()
}

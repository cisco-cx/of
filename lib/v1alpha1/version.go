package v1alpha1

// Represents build info.
var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion string
)

type Versioner interface {
	Print() string
	BuildContext() string
	Info() string
}

package helm

import "time"

// Define release statuses
const (
	// StatusUnknown indicates that a release is in an uncertain state
	StatusUnknown Status = "unknown"
	// StatusDeployed indicates that the release has been pushed to Kubernetes
	StatusDeployed Status = "deployed"
	// StatusUninstalled indicates that a release has been uninstalled from Kubernetes
	StatusUninstalled Status = "uninstalled"
	// StatusSuperseded indicates that this release object is outdated and a newer one exists
	StatusSuperseded Status = "superseded"
	// StatusFailed indicates that the release was not successfully deployed
	StatusFailed Status = "failed"
	// StatusUninstalling indicates that a uninstall operation is underway
	StatusUninstalling Status = "uninstalling"
	// StatusPendingInstall indicates that an install operation is underway
	StatusPendingInstall Status = "pending-install"
	// StatusPendingUpgrade indicates that an upgrade operation is underway
	StatusPendingUpgrade Status = "pending-upgrade"
	// StatusPendingRollback indicates that an rollback operation is underway
	StatusPendingRollback Status = "pending-rollback"
)

// Release describes a generic chart deployment
type Release struct {
	Name      string
	Namespace string
	Chart     *Chart
	Info      *Info
	Values    map[string]interface{}
	Manifest  string
	Version   int
}

// Info holds metadata of a chart deployment
type Info struct {
	LastDeployed time.Time
	Description  string
	Status       Status
}

// Chart describes the chart for a release
type Chart struct {
	Name       string
	Version    string
	AppVersion string
	Values     Values
	Templates  []*File
}

// File represents a file as a name/value pair.
// The name is a relative path within the scope
// of the chart's base directory.
type File struct {
	Name string
	Data []byte
}

// Status holds the status of a release
type Status string

// AllowsUpgrade returns true if the status allows the release
// to be upgraded. This is currently only the case if it equals
// `StatusDeployed`.
func (s Status) AllowsUpgrade() bool {
	return s == StatusDeployed
}

// String returns the Status as a string
func (s Status) String() string {
	return string(s)
}

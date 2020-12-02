package helm

import "time"

// GetOptions holds the options available for Helm uninstall
// operations, the version implementation _must_ implement all
// fields supported by that version but can (silently) ignore
// unsupported set values.
type GetOptions struct {
	Namespace string
	Version   int
}

// UpgradeOptions holds the options available for Helm upgrade
// operations, the version implementation _must_ implement all
// fields supported by that version but can (silently) ignore
// unsupported set values.
type UpgradeOptions struct {
	Namespace    string
	Timeout      time.Duration
	Wait         bool
	Install      bool
	DisableHooks bool
	DryRun       bool
	ClientOnly   bool
	Force        bool
	ResetValues  bool
	SkipCRDs     bool
	ReuseValues  bool
	Recreate     bool
	MaxHistory   int
	Atomic       bool
}

// RollbackOptions holds the options available for Helm rollback
// operations, the version implementation _must_ implement all
// fields supported by that version but can (silently) ignore
// unsupported set values.
type RollbackOptions struct {
	Namespace    string
	Version      int
	Timeout      time.Duration
	Wait         bool
	DisableHooks bool
	DryRun       bool
	Recreate     bool
	Force        bool
}

// UninstallOptions holds the options available for Helm uninstall
// operations, the version implementation _must_ implement all
// fields supported by that version but can (silently) ignore
// unsupported set values.
type UninstallOptions struct {
	Namespace    string
	DisableHooks bool
	DryRun       bool
	KeepHistory  bool
	Timeout      time.Duration
}

// HistoryOption holds the options available for Helm history
// operations, the version implementation _must_ implement all
// fields supported by that version but can (silently) ignore
// unsupported set values.
type HistoryOptions struct {
	Namespace string
	Max       int
}

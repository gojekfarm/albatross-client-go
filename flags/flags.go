package flags

import "errors"

// CommonFlags are common to all apis
// TODO: We can maybe define setter funcs to allow setting these easily
type CommonFlags struct {
	KubeContext   string `json:"-" schema:"-"`
	KubeToken     string `json:"kube_token,omitempty" schema:"-"`
	KubeAPIServer string `json:"kube_apiserver,omitempty" schema:"-"`
	Namespace     string `json:"-" schema:"-"`
}

// InstallFlags defines flags supported by the install api
type InstallFlags struct {
	DryRun  bool   `json:"dry_run"`
	Version string `json:"version"`
	CommonFlags
}

// ListFlags defines flags supported by the List api
type ListFlags struct {
	AllNamespaces bool `schema:"-"`
	Deployed      bool `schema:"deployed,omitempty"`
	Failed        bool `schema:"failed,omitempty"`
	Pending       bool `schema:"pending,omitempty"`
	Uninstalled   bool `schema:"uninstalled,omitempty"`
	Uninstalling  bool `schema:"uninstalling,omitempty"`
	CommonFlags
}

// UpgradeFlags defines flags supported by the upgrade api
type UpgradeFlags struct {
	DryRun  bool   `json:"dry_run"`
	Version string `json:"version"`
	Install bool   `json:"install"`
	CommonFlags
}

type StatusFlags struct {
	Revision int `schema:"revision,omitempty"`
	CommonFlags
}

type UninstallFlags struct {
	DryRun       bool `schema:"dry_run,omitempty"`
	DisableHooks bool `schema:"disable_hooks,omitempty"`
	KeepHistory  bool `schema:"keep_history,omitempty"`
	Timeout      int  `schema:"timeout,omitempty"`
	CommonFlags
}

func (u *UpgradeFlags) Valid() error {
	if u.KubeContext == "" {
		return errors.New("kube context is a required parameter")
	}

	if u.Namespace == "" {
		u.Namespace = "default"
	}

	return nil
}

func (u *InstallFlags) Valid() error {
	if u.KubeContext == "" {
		return errors.New("kube context is a required parameter")
	}

	if u.Namespace == "" {
		u.Namespace = "default"
	}

	return nil
}

func (l *ListFlags) Valid() error {
	if l.KubeContext == "" {
		return errors.New("kube context is a required parameter")
	}
	if !l.AllNamespaces && l.Namespace == "" {
		l.Namespace = "default"
	}
	return nil
}

func (s *StatusFlags) Valid() error {
	if s.KubeContext == "" {
		return errors.New("kube context is a required parameter")
	}
	if s.Namespace == "" {
		s.Namespace = "default"
	}
	return nil
}

func (u *UninstallFlags) Valid() error {
	if u.KubeContext == "" {
		return errors.New("kube context is a required parameter")
	}

	if u.Namespace == "" {
		u.Namespace = "default"
	}

	return nil
}

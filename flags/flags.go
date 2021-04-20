package flags

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

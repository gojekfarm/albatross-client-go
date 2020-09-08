package flags

// CommonFlags are common to all apis
// TODO: We can maybe define setter funcs to allow setting these easily
type CommonFlags struct {
	KubeCtx       string `json:"kube_context,omitempty"`
	KubeToken     string `json:"kube_token,omitempty"`
	KubeAPIServer string `json:"kube_apiserver,omitempty"`
	Namespace     string
}

// InstallFlags defines flags supported by the install api
type InstallFlags struct {
	DryRun  bool `json:"dry_run"`
	Version string
	CommonFlags
}

// ListFlags defines flags supported by the List api
type ListFlags struct {
	AllNamespaces bool `json:"all-namespaces,omitempty"`
	Deployed      bool `json:"deployed,omitempty"`
	Failed        bool `json:"failed,omitempty"`
	Pending       bool `json:"pending,omitempty"`
	Uninstalled   bool `json:"uninstalled,omitempty"`
	Uninstalling  bool `json:"uninstalling,omitempty"`
	CommonFlags
}

// UpgradeFlags defines flags supported by the upgrade api
type UpgradeFlags struct {
	DryRun  bool `json:"dry_run"`
	Version string
	Install bool
	CommonFlags
}

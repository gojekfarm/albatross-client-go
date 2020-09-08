package release

import (
	"time"
)

// Release represents a helm release. All apis that return a release should return
// an instance of this struct to make the api response consistent for all apis
type Release struct {
	Name       string    `json:"name"`
	Namespace  string    `json:"namespace"`
	Version    int       `json:"version"`
	Updated    time.Time `json:"updated_at,omitempty"`
	Status     string    `json:"status"`
	Chart      string    `json:"chart"`
	AppVersion string    `json:"app_version"`
}

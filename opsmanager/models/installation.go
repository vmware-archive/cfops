package models

// Ops Manager installation json types
type (
	InstallationSettings struct {
		Infrastructure Infrastructure `json:"infrastructure"`
		Products       []Products     `json:"products"`
	}

	Infrastructure struct {
		Type string `json:"type"`
	}

	Products struct {
		Type string              `json:"type"`
		IPS  map[string][]string `json:"ips"`
		Jobs []Jobs              `json:"jobs"`
	}

	Jobs struct {
		Type       string       `json:"type"`
		Properties []Properties `json:"properties"`
	}

	Properties struct {
		Definition string      `json:"definition"`
		Value      interface{} `json:"value"`
	}
)

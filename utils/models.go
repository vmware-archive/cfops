package utils

type (
	DeploymentObject struct {
		Name string
	}

	VMObject struct {
		AgentId string `json:"agent_id"`
		Cid     string `json:"cid"`
		Job     string `json:"job"`
		Index   int    `json:"index"`
	}

	InstallationObject struct {
		Infrastructure Infrastructure `json:"infrastructure"`
		Products       []Products     `json:"products,components"`
	}

	Infrastructure struct {
		Type string `json:"type"`
	}

	Products struct {
		Type string              `json:"type"`
		IPS  map[string][]string `json:"ips"`
		Jobs []Jobs              `json:"jobs"`
	}

	IPS struct {
		identifier string
		value      []string
	}

	Jobs struct {
		Type       string       `json:"type"`
		Properties []Properties `json:"properties"`
	}

	Properties struct {
		Definition string `json:"definition"`
		Value      Value  `json:"value"`
	}

	Value interface {
	}

	EventObject struct {
		Id          int    `json:"id"`
		State       string `json:"state"`
		Description string `json:"description"`
		Result      string `json:"result"`
	}
)

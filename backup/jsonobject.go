package backup

type JsonObject struct {
	Infrastructure Infrastructure `json:"infrastructure"`
	Products       []Products     `json:"products"`
}

type Infrastructure struct {
	Type string `json:"type"`
}

type Products struct {
	Type string              `json:"type"`
	IPS  map[string][]string `json:"ips"`
	Jobs []Jobs              `json:"jobs"`
}

type IPS struct {
	identifier string
	value      []string
}

type Jobs struct {
	Type       string       `json:"type"`
	Properties []Properties `json:"properties"`
}

type Properties struct {
	Definition string `json:"definition"`
	Value      Value  `json:"value"`
}

type Value interface {
}

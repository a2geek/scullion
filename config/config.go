package config

type Config struct {
	Library   []LibraryDef      `json:"library"`
	Templates map[string]string `json:"templates"`
	Rules     []RuleDef         `json:"rules"`
}

type LibraryDef struct {
	Name     string   `json:"name"`
	Pipeline []string `json:"pipeline"`
	Actions  []string `json:"actions"`
}

type RuleDef struct {
	Name     string   `json:"name"`
	Schedule Schedule `json:"schedule"`
	Pipeline []string `json:"pipeline"`
	Actions  []string `json:"actions"`
}

type Schedule struct {
	Frequency string `json:"frequency"`
}

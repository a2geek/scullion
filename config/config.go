package config

type Config struct {
	Library   []LibraryDef      `json:"library" yaml:"library"`
	Templates map[string]string `json:"templates" yaml:"templates"`
	Rules     []RuleDef         `json:"rules" yaml:"rules"`
}

type LibraryDef struct {
	Name     string   `json:"name" yaml:"name"`
	Pipeline []string `json:"pipeline" yaml:"pipeline"`
	Actions  []string `json:"actions" yaml:"actions"`
}

type RuleDef struct {
	Name     string   `json:"name" yaml:"name"`
	Schedule Schedule `json:"schedule" yaml:"schedule"`
	Pipeline []string `json:"pipeline" yaml:"pipeline"`
	Actions  []string `json:"actions" yaml:"actions"`
}

type Schedule struct {
	Frequency string `json:"frequency" yaml:"frequency"`
}

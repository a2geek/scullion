package config

type TaskDef struct {
	Name     string   `json:"name"`
	Schedule Schedule `json:"schedule"`
	Filters  Filter   `json:"filters"`
}

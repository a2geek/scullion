package main

type Task struct {
	Name     string   `json:"name"`
	Schedule Schedule `json:"schedule"`
	Filters  Filter   `json:"filters"`
}
type Schedule struct {
	Frequency string `json:"frequency"`
}
type Filter struct {
	Organization string `json:"organization"`
	Space        string `json:"space"`
	Application  string `json:"application"`
	Action       string `json:"action"`
}

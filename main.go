package main

import (
	"scullion/cmd"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Enable verbose output"`

	Run      cmd.Run      `command:"run" alias:"r"`
	Validate cmd.Validate `command:"validate" alias:"v"`
}

func main() {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	parser.NamespaceDelimiter = "-"
	parser.Parse()
}

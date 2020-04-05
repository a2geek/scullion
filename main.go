package main

import (
	"scullion/cmd"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"enable verbose output"`

	Run      cmd.RunCommand      `command:"run" alias:"r"`
	Validate cmd.ValidateCommand `command:"validate" alias:"v"`
}

func main() {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	parser.NamespaceDelimiter = "-"
	parser.Parse()
}

package main

import (
	"scullion/cmd"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Run         cmd.Run         `command:"run" alias:"r" description:"Run rules continuously"`
	Validate    cmd.Validate    `command:"validate" alias:"v" description:"Validate rules"`
	Disassemble cmd.Disassemble `command:"disassemble" alias:"disasm" alias:"d" alias:"dump" description:"Disassemble rules byte code... for the curious"`
	OneTime     cmd.OneTime     `command:"one-time" alias:"once" description:"Run rules once"`
	Reference   cmd.Reference   `command:"reference" alias:"ref" description:"Generate reference information"`
}

func main() {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	parser.NamespaceDelimiter = "-"
	parser.Parse()
}

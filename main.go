package main

import (
	"scullion/cmd"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Run         cmd.Run         `command:"run" alias:"r"`
	Validate    cmd.Validate    `command:"validate" alias:"v"`
	Disassemble cmd.Disassemble `command:"disassemble" alias:"disasm" alias:"d" alias:"dump"`
	OneTime     cmd.OneTime     `command:"one-time" alias:"o" alias:"once"`
}

func main() {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	parser.NamespaceDelimiter = "-"
	parser.Parse()
}

package main

import (
	"strings"
)

// A command represents a piece of a pipeline.
type Command struct {
	// The command name to run.
	Command string

	// The arguments to give it.
	Arguments []string
}

// Parse a pipeline into a slice of command instances.
// For instance, echo foo | cat will turn into two commands: echo (foo), and cat.
func parsePipeline(cmd string) []Command {
	things := strings.Split(cmd, "|")
	cmds := []Command{}
	for _, thing := range things {
		thing = strings.Trim(thing, " ")
		vals := strings.Split(thing, " ")
		cmds = append(cmds, Command{Command: vals[0], Arguments: vals[1:]})
	}
	return cmds
}

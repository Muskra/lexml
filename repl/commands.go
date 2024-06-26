package repl

import (
	"os"
	"pmp/lexml"
)

// Command Type represent a command with it's parameters
type Command struct {
	Kind       uint8
	Keyword    string
	Parameters []string
}

// NewCommand function allocate and returns an empty Command Type
func NewCommand() Command {
	return Command{
		Kind:       0,
		Keyword:    "",
		Parameters: make([]string, 0),
	}
}

// push function 'pushes' or add a command result to a given value stored in a Repl.Values map Type
func push(repl Repl, parameters []string) (lexml.Data, error) {
	return lexml.Data{}, nil
}

// pop function 'pops' or delete a value from a Repl.Value map Type
func pop(repl Repl, parameters []string) (lexml.Data, error) {
	return lexml.Data{}, nil
}

// exit function terminates the program
func exit() {
	os.Exit(0)
}

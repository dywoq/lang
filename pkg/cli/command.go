package cli

import (
	"fmt"
	"os"
)

// Command represents the CLI command,
// with the name and execute method.
type Command interface {
	Name() string
	Execute() error
}

var commands map[string]Command

// RegisterCommand registers the command into
// the underlying map.
// Does nothing if command already exists.
func RegisterCommand(c Command) {
	if _, ok := commands[c.Name()]; ok {
		return
	}
	commands[c.Name()] = c
}

// ExecuteCommand executes the command,
// whose name is equal to os.Args[1].
//
// Otherwise, the error is returned
// if command couldn't be found.
func ExecuteCommand() error {
	cmd, ok := commands[os.Args[1]]
	if !ok {
		return fmt.Errorf("cli: unknown command: %s", os.Args[1])
	}
	return cmd.Execute()
}

// GivenCommand returns the given command,
// always equal to os.Args[1].
func GivenCommand() string {
	return os.Args[1]
}

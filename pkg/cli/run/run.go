package run

// Command implements cli.Command.
type Command struct{}

// Name returns the name of the command,
// specifically "run".
func (c *Command) Name() string {
	return "run"
}

// Execute evaluates the program,
// or multiple files with specified entry function.
func (c *Command) Execute() error {
	panic("TODO: implement me!")
}

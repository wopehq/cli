package cli

import (
	"errors"
	"fmt"
)

var ErrCommandNotFound = errors.New("command is not found in configurator's commands")
var ErrParameterNotFound = errors.New("wanted parameter is not found")

// Configurator handles all things about parsing args, stores the main command
type Configurator struct {
	Args    []string
	Command *Command
}

func NewConfigurator(args []string) *Configurator {
	return &Configurator{
		Args: args,
	}
}

// SetMainCommand sets configurator's main command
func (c *Configurator) SetMainCommand(command *Command) {
	c.Command = command
}

// FindHelp recognizes if any help command is received and returns a help
func (c *Configurator) FindHelp() (*Help, error) {
	var help *Help
	if len(c.Args) == 0 && c.Command.handler == nil {
		help = c.Command.Help()
	}

	if len(c.Args) > 0 {
		if c.Args[0] == "help" {
			if len(c.Args) == 1 {
				help = c.Command.Help()
			}

			if len(c.Args) > 1 {
				cmd, _, err := c.ParseCommand(c.Args[1:])
				if err != nil {
					return nil, err
				}
				help = cmd.Help()
			}
		}
	}

	return help, nil
}

// ParseCommand returns the current command for example
// "app run concurrent -v"
// concurrent is the current command here
func (c *Configurator) ParseCommand(args []string) (*Command, int, error) {
	currentCommand := c.Command
	var position int
	for _, arg := range args {
		if isFlagArg(arg) {
			break
		}

		cmd, exists := currentCommand.FindCommand(arg)
		if !exists {
			return currentCommand, position, fmt.Errorf("%s %w", arg, ErrCommandNotFound)
		}
		currentCommand = cmd
		position += 1
	}
	return currentCommand, position, nil
}

// Parse parses the wanted command and its flags
func (c *Configurator) Parse() (*Command, error) {
	args := c.Args
	cmd, position, err := c.ParseCommand(args)
	if err != nil {
		return nil, err
	}

	if len(c.Args) > 0 && cmd != c.Command {
		args = c.Args[position:]
	}

	err = cmd.Flagset.Parse(args)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// isFlagArg returns if the given arg is a flag or not
func isFlagArg(arg string) bool {
	return ((len(arg) >= 3 && arg[1] == '-') ||
		(len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'))
}

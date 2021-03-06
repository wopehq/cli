package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

var bold = color.New(color.Bold)
var green = color.New(color.FgGreen, color.Bold)

var ErrParameterIsNotBool = errors.New("wanted parameter is not a boolean")
var ErrParameterIsNotInt = errors.New("wanted parameter is not an integer")
var ErrParameterIsNotString = errors.New("wanted parameter is not a string")
var ErrCommandHandlerNotFound = errors.New("command's handler function is not found")
var ErrCommandNotFound = errors.New("command is not found in configurator's commands")
var ErrParameterNotFound = errors.New("wanted parameter is not found")

type Command struct {
	Name        string
	Use         string
	Description string
	Parameters  []Parameter
	Commands    Commands
	Flagset     *flag.FlagSet
	handler     func(cmd *Command)
}

type Commands []*Command

func NewCommand(name string, use string, description string) *Command {
	command := &Command{
		Name:        name,
		Use:         use,
		Description: description,
	}

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = command.Help().ShowHelp
	command.Flagset = fs

	return command
}

// Do sets the command's handler function
func (c *Command) Do(handler func(cmd *Command)) {
	c.handler = handler
}

// AddCommand adds sub commands to the command
func (c *Command) AddCommand(command *Command) *Command {
	c.Commands = append(c.Commands, command)
	return command
}

// Parse parses the args returns the wanted command
func (c *Command) Parse(args []string) (*Command, error) {
	cmd, position, err := c.getCurrentCommand(args)
	if err != nil {
		return nil, err
	}

	if len(args) > 0 && cmd != c {
		args = args[position:]
	}

	err = cmd.Flagset.Parse(args)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// getCurrentCommand returns the current command for example
// "app run concurrent -v"
// concurrent is the current command here
func (c *Command) getCurrentCommand(args []string) (*Command, int, error) {
	currentCommand := c
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

// ShowHelp prints the usage of the command
func (c *Command) Help() *Help {
	var help Help
	help.Usage = green.Sprintf("Usage: %s \n\n", c.Use)

	if len(c.Commands) > 0 {
		help.Commands += bold.Sprintf("Commands: \n")
		for _, command := range c.Commands {
			help.Commands += bold.Sprintf("%-20s %s \n", command.Name, command.Use)
		}
	}

	if len(c.Parameters) > 0 {
		help.Parameters += bold.Sprintf("Parameters: \n")
		for _, parameter := range c.Parameters {
			help.Parameters += bold.Sprintf("--%s, -%-10s %-30s %s \n", parameter.Name, parameter.Shortname, parameter.Use, parameter.Description)
		}
	}

	return &help
}

// Run runs the command's handler function
func (c *Command) Run() error {
	if c.handler == nil {
		return fmt.Errorf("%s %w", c.Name, ErrCommandHandlerNotFound)
	}

	c.handler(c)
	return nil
}

// FindHelp recognizes if any help command is received and returns a help
func (c *Command) FindHelp(args []string) (*Help, error) {
	if len(args) == 0 {
		return nil, nil
	}

	var help *Help
	if len(args) == 0 && c.handler == nil {
		help = c.Help()
	}

	if len(args) > 0 {
		if args[0] == "help" {
			if len(args) == 1 {
				help = c.Help()
			}

			if len(args) > 1 {
				cmd, _, err := c.getCurrentCommand(args[1:])
				if err != nil {
					return nil, err
				}
				help = cmd.Help()
			}
		}
	}

	return help, nil
}

// FindCommand searches the sub commands of the command
func (c *Command) FindCommand(name string) (*Command, bool) {
	for _, command := range c.Commands {
		if command.Name == name {
			return command, true
		}
	}
	return nil, false
}

func (c *Command) addParameter(parameter *Parameter) {
	c.Parameters = append(c.Parameters, *parameter)
}

// AddBoolParameter sets a bool flag in the command's flagset
func (c *Command) AddBoolParameter(parameter *Parameter, value *bool, defaultValue bool) *Parameter {
	c.addParameter(parameter)
	for _, alias := range parameter.Aliases() {
		c.Flagset.BoolVar(value, alias, defaultValue, parameter.Use)
	}

	return parameter
}

// AddIntParameter sets an int flag in the command's flagset
func (c *Command) AddIntParameter(parameter *Parameter, value *int, defaultValue int) *Parameter {
	c.addParameter(parameter)
	for _, alias := range parameter.Aliases() {
		c.Flagset.IntVar(value, alias, defaultValue, parameter.Use)
	}

	return parameter
}

// AddStringParameter sets a string flag in the command's flagset
func (c *Command) AddStringParameter(parameter *Parameter, value *string, defaultValue string) *Parameter {
	c.addParameter(parameter)
	for _, alias := range parameter.Aliases() {
		c.Flagset.StringVar(value, alias, defaultValue, parameter.Use)
	}

	return parameter
}

// GetBool gets a bool value with the given name from the command's flagset
func (c *Command) GetBool(name string) (bool, error) {
	value := c.Flagset.Lookup(name).Value.String()
	if value == "" {
		return false, ErrParameterNotFound
	}
	return strconv.ParseBool(value)
}

// GetString gets a string value with the given name from the command's flagset
func (c *Command) GetString(name string) (string, error) {
	value := c.Flagset.Lookup(name).Value.String()
	if value == "" {
		return "", ErrParameterNotFound
	}
	return value, nil
}

// GetInt gets an int value with the given name from the command's flagset
func (c *Command) GetInt(name string) (int, error) {
	value := c.Flagset.Lookup(name).Value.String()
	if value == "" {
		return 0, ErrParameterNotFound
	}
	return strconv.Atoi(value)
}

// isFlagArg returns if the given arg is a flag or not
func isFlagArg(arg string) bool {
	return ((len(arg) >= 3 && arg[1] == '-') ||
		(len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'))
}

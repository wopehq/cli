package cli

import (
	"flag"
	"fmt"
	"strings"
	"testing"
)

type RunTest struct {
	Name                string
	Args                []string
	Command             *Command
	SubCommands         []*Command
	ExpectedCommandName string
}

type GetCurrentCommandTest struct {
	Name                   string
	Command                *Command
	SubCommands            []*Command
	Args                   []string
	ExpectedCurrentCommand string
}

type SubCommandsTest struct {
	Name                string
	Command             *Command
	ExpectedSubCommands int
}

type ParameterTest struct {
	Name                   string
	Command                *Command
	Parameter              Parameter
	Args                   []string
	ExpectedParameterValue interface{}
}

type FindHelpTest struct {
	Name            string
	Args            []string
	Command         *Command
	SubCommands     []*Command
	ExpectedStrings []string
}

func TestRun(t *testing.T) {
	tests := []RunTest{
		{
			Name:                "simple do command",
			Args:                []string{"do"},
			Command:             NewCommand("seo", "seo [command]", ""),
			ExpectedCommandName: "do",
			SubCommands: []*Command{
				{
					Name: "do",
					Use:  "seo do [flags]",
					Parameters: []Parameter{
						{
							Name:      "message",
							Shortname: "m",
						},
					},
					Flagset: flag.NewFlagSet("cli", flag.ContinueOnError),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var commandName string
			for _, cmd := range test.SubCommands {
				cmd.Do(func(cmd *Command) {
					if cmd.Name == test.ExpectedCommandName {
						commandName = cmd.Name
					} else {
						t.Error("called the wrong command's handler function")
					}
				})
				test.Command.AddCommand(cmd)
			}

			cmd, err := test.Command.Parse(test.Args)
			if err != nil {
				t.Fatal(err)
			}
			cmd.Run()

			if commandName == "" {
				t.Error("no any command's handler function is called")
			}
		})
	}
}

func TestFindHelp(t *testing.T) {
	tests := []FindHelpTest{
		{
			Name:            "a simple help test",
			Args:            []string{"help", "do"},
			Command:         NewCommand("seo", "seo [command]", "seo command"),
			ExpectedStrings: []string{"message value"},
			SubCommands: []*Command{
				{
					Name: "do",
					Use:  "seo do [flags]",
					Parameters: []Parameter{
						{
							Name:        "message",
							Shortname:   "m",
							Description: "message value",
						},
					},
					Flagset: flag.NewFlagSet("cli", flag.ContinueOnError),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for _, cmd := range test.SubCommands {
				test.Command.AddCommand(cmd)
			}

			help, err := test.Command.FindHelp(test.Args)
			if err != nil {
				t.Fatal(err)
			}

			all := help.Usage + help.Parameters + help.Commands
			for _, expectedString := range test.ExpectedStrings {
				if !strings.Contains(all, expectedString) {
					t.Errorf("help string is not contains any expected string")
				}
			}
		})
	}
}

func TestGetCurrentCommand(t *testing.T) {
	tests := []GetCurrentCommandTest{
		{
			Name:                   "one simple command",
			Args:                   []string{"do"},
			Command:                NewCommand("seo", "seo [command]", ""),
			SubCommands:            []*Command{NewCommand("do", "seo do [command]", "")},
			ExpectedCurrentCommand: "do",
		},
		{
			Name:    "a subcommand",
			Args:    []string{"do", "concurrent"},
			Command: NewCommand("seo", "seo [command]", ""),
			SubCommands: []*Command{
				{
					Name: "do",
					Use:  "seo do [command]",
					Commands: []*Command{
						{
							Name: "concurrent",
							Use:  "seo do concurrent",
						},
					},
				},
			},
			ExpectedCurrentCommand: "concurrent",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for _, cmd := range test.SubCommands {
				test.Command.AddCommand(cmd)
			}

			command, _, err := test.Command.getCurrentCommand(test.Args)
			if err != nil {
				t.Fatal(err)
			}

			if command.Name != test.ExpectedCurrentCommand {
				t.Errorf("command name %s is not equal to the expected command name %s", command.Name, test.ExpectedCurrentCommand)
			}
		})
	}
}

func TestParameters(t *testing.T) {
	tests := []ParameterTest{
		{
			Name:    "a string parameter test",
			Command: NewCommand("do", "seo do", ""),
			Parameter: Parameter{
				Name:      "message",
				Shortname: "m",
			},
			Args:                   []string{"--message", "seo do rocks"},
			ExpectedParameterValue: "seo do rocks",
		},
		{
			Name:    "a boolean parameter test",
			Command: NewCommand("do", "seo do", ""),
			Parameter: Parameter{
				Name:      "verbose",
				Shortname: "v",
			},
			Args:                   []string{"--verbose"},
			ExpectedParameterValue: true,
		},
		{
			Name:    "an integer parameter test",
			Command: NewCommand("do", "seo do", ""),
			Parameter: Parameter{
				Name:      "count",
				Shortname: "c",
			},
			Args:                   []string{"--count", "16"},
			ExpectedParameterValue: 16,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			switch expectedValue := test.ExpectedParameterValue.(type) {
			case string:
				var value string
				test.Command.AddStringParameter(&test.Parameter, &value, "")
				err := test.Command.Flagset.Parse(test.Args)
				if err != nil {
					t.Fatal(err)
				}

				if value != expectedValue {
					t.Errorf("parsed value %s is not equal to the expected value %s", value, expectedValue)
				}

				value, err = test.Command.GetString(test.Parameter.Name)
				if err != nil {
					t.Fatal(err)
				}

				if value != expectedValue {
					t.Errorf("parsed value taken by GetString %s is not equal to the expected value %s", value, expectedValue)
				}
			case int:
				var value int
				test.Command.AddIntParameter(&test.Parameter, &value, 0)
				err := test.Command.Flagset.Parse(test.Args)
				if err != nil {
					t.Fatal(err)
				}

				if value != expectedValue {
					t.Errorf("parsed value %d is not equal to the expected value %d", value, expectedValue)
				}

				value, err = test.Command.GetInt(test.Parameter.Name)
				if err != nil {
					t.Fatal(err)
				}

				if value != expectedValue {
					t.Errorf("parsed value taken by GetInt %d is not equal to the expected value %d", value, expectedValue)
				}
			case bool:
				var value bool
				test.Command.AddBoolParameter(&test.Parameter, &value, false)
				err := test.Command.Flagset.Parse(test.Args)
				if err != nil {
					t.Fatal(err)
				}
				if value != expectedValue {
					t.Errorf("parsed value %t is not equal to the expected value %t", value, expectedValue)
				}

				value, err = test.Command.GetBool(test.Parameter.Name)
				if err != nil {
					t.Fatal(err)
				}

				if value != expectedValue {
					t.Errorf("parsed value taken by GetBool %t is not equal to the expected value %t", value, expectedValue)
				}
			}
		})
	}
}

func TestSubCommands(t *testing.T) {
	tests := []SubCommandsTest{
		{
			Name:                "a test with 10 sub commands",
			Command:             NewCommand("do", "seo do", ""),
			ExpectedSubCommands: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for i := 0; i < test.ExpectedSubCommands; i++ {
				test.Command.AddCommand(
					NewCommand(fmt.Sprintf("command_%d", i), "", ""),
				)
			}

			if len(test.Command.Commands) != test.ExpectedSubCommands {
				t.Errorf("command's subcommands length %d is not equal to the expected sub commands length %d", len(test.Command.Commands), test.ExpectedSubCommands)
			}

			for i := 0; i < test.ExpectedSubCommands; i++ {
				_, exists := test.Command.FindCommand(fmt.Sprintf("command_%d", i))
				if !exists {
					t.Errorf("find command method can't find the added command")
				}
			}
		})
	}
}

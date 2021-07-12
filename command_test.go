package cli

import (
	"fmt"
	"testing"
)

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
	ParameterType          string
	ExpectedParameterValue interface{}
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

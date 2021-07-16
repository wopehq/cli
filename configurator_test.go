package cli

import (
	"flag"
	"strings"
	"testing"
)

type ParseCommandTest struct {
	Name                   string
	Command                *Command
	SubCommands            []*Command
	Configurator           *Configurator
	ExpectedCurrentCommand string
}

type InitializeTest struct {
	Name                string
	Configurator        *Configurator
	Command             *Command
	SubCommands         []*Command
	ExpectedCommandName string
}

type HelpTest struct {
	Name            string
	Configurator    *Configurator
	Command         *Command
	SubCommands     []*Command
	ExpectedStrings []string
}

func TestFindHelp(t *testing.T) {
	tests := []HelpTest{
		{
			Name:            "a simple help test",
			Configurator:    NewConfigurator([]string{"help", "do"}),
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
			test.Configurator.SetMainCommand(test.Command)

			help, err := test.Configurator.FindHelp()
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

func TestInitialize(t *testing.T) {
	tests := []InitializeTest{
		{
			Name:                "simple do command",
			Configurator:        NewConfigurator([]string{"do"}),
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

			test.Configurator.SetMainCommand(test.Command)

			cmd, err := test.Configurator.Parse()
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

func TestParseCommand(t *testing.T) {
	tests := []ParseCommandTest{
		{
			Name:                   "one simple command",
			Configurator:           NewConfigurator([]string{"do"}),
			Command:                NewCommand("seo", "seo [command]", ""),
			SubCommands:            []*Command{NewCommand("do", "seo do [command]", "")},
			ExpectedCurrentCommand: "do",
		},
		{
			Name:         "a subcommand",
			Configurator: NewConfigurator([]string{"do", "concurrent"}),
			Command:      NewCommand("seo", "seo [command]", ""),
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

			test.Configurator.SetMainCommand(test.Command)
			command, _, err := test.Configurator.ParseCommand(test.Configurator.Args)
			if err != nil {
				t.Fatal(err)
			}

			if command.Name != test.ExpectedCurrentCommand {
				t.Errorf("command name %s is not equal to the expected command name %s", command.Name, test.ExpectedCurrentCommand)
			}
		})
	}
}

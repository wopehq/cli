package main

import (
	"fmt"
	"log"
	"os"

	"github.com/teamseodo/cli"
)

func main() {
	c := cli.NewConfigurator(os.Args[1:])

	mainCommand := cli.NewCommand("app", "app [command] [flags]", "description about the app")
	c.SetMainCommand(mainCommand)

	printCommand := cli.NewCommand("print", "print [flags]", "prints the values")
	mainCommand.AddCommand(printCommand)

	var messageParameter string
	printCommand.AddStringParameter(&cli.Parameter{
		Use:       "app print --message [value]",
		Name:      "message",
		Shortname: "m",
	}, &messageParameter, "hello world")

	printCommand.Do(func(cmd *cli.Command) {
		fmt.Println(messageParameter)
	})

	help, err := c.FindHelp()
	if err != nil {
		log.Fatal(err)
	}
	help.ShowHelp()

	cmd, err := c.Parse()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Run()
}

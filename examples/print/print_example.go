package main

import (
	"fmt"
	"log"
	"os"

	"github.com/teamseodo/cli"
)

func main() {
	mainCommand := cli.NewCommand("app", "app [command] [flags]", "description about the app")

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

	help, err := mainCommand.FindHelp(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	help.ShowHelp()

	cmd, err := mainCommand.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run()
	if err != nil {
		cmd.Help().ShowHelp()
	}
}

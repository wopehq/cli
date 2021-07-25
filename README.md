
<div align="center">
<pre>
 ▄████████  ▄█        ▄█  
███    ███ ███       ███  
███    █▀  ███       ███▌ 
███        ███       ███▌ 
███        ███       ███▌ 
███    █▄  ███       ███  
███    ███ ███▌    ▄ ███  
████████▀  █████▄▄██ █▀   
           ▀              
                    


a lightweight and simple cli package
</pre>
[![Go Reference](https://pkg.go.dev/badge/github.com/teamseodo/cli.svg)](https://pkg.go.dev/github.com/teamseodo/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/teamseodo/cli)](https://goreportcard.com/report/github.com/teamseodo/cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
</div>


## Contents
- [Get Started](#Get-Started)
- [Commands](#Commands)
    - [Parameters](#Parameters)
- [Help Messages](#Help-Messages)
- [Contribute](#Contribute)

## Get Started

You can use this command to install the package
```bash
go get github.com/teamseodo/cli
```

A example code for the basics of the cli
```go
package main

import "github.com/teamseodo/cli"

func main() {
    mainCommand := cli.NewCommand("app", "app [command] [flags]", "description about the app")

    say := cli.NewCommand("say", "say [flags]", "prints the values")
    mainCommand.AddCommand(say)

    var messageParameter string
    say.AddStringParameter(&cli.Parameter{
        Use: "app print --message [value]",
        Name: "message",
        Shortname: "m",
    }, &messageParameter, "hello world")

    say.Do(func(cmd *cli.Command) {
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
```

## Commands
You need to create a configurator to use commands or anything in this package. Configurator takes a main command and you can add sub commands to the main command. 

You can create a main command like that,
```go
mainCommand := cli.NewCommand("app", "app [command] [flags]", "description about the app")
```

Every command can have multiple sub commands, you can add a subcommand like that
```go
hi := cli.NewCommand("hi", "hi", "prints a hello message back")
mainCommand.AddCommand(hi)
```
Now you will get an structure in command line like that,
```shell
app hi
```

If you want to add a functionality to the command, you can use the `command.Do` method.
```go
hi.Do(func (cmd *cli.Command) {
    fmt.Println("hi")
})
```
Parse the args after all things are done and run the returned command
```go
cmd, err := mainCommand.Parse()
if err != nil {
	log.Fatal(err)
}

err = cmd.Run()
if err != nil {
	cmd.Help().ShowHelp()
}
```
Now when you type `app hi`, application will print a `"hi"` message.

**You can add a functionality to the main command either, so it will run when no command is received.**

### Parameters
Every command can take multiple parameters (also known as: flags)
You can add three types of parameters, string, int and bool.

```go
AddBoolParameter(parameter *Parameter, value *bool, defaultValue bool) *Parameter

AddIntParameter(parameter *Parameter, value *int, defaultValue int) *Parameter

AddStringParameter(parameter *Parameter, value *string, defaultValue string) *Parameter
```

If you want to add a string parameter, you can add the parameter like this
```go
var messageParameter string
printCommand.AddStringParameter(&cli.Parameter{
	Use:       "app print --message [value]",
	Name:      "message",
	Shortname: "m",
}, &messageParameter, "hello world")
```

## Help Messages
If you want to print a help message for the command you need to run the `FindHelp()` method before configurator initialization.
```go
help, err := mainCommand .FindHelp()
if err != nil {
	log.Fatal(err)
}
help.ShowHelp()
```
When running the wanted command parsing errors can be occur, so you can make an error check when you run the wanted command and print a help.
```go
err = cmd.Run()
if err != nil {
	cmd.Help().ShowHelp()
}
```
## Contribute
Pull requests are welcome. please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
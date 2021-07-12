
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
</div>

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
    c := cli.NewConfigurator(os.Args[1:])

    mainCommand := cli.NewCommand("app", "app [command] [flags]", "description about the app")
    c.SetMainCommand(mainCommand)

    printCommand := cli.NewCommand("print", "print [flags]", "prints the values")
    mainCommand.AddCommand(printCommand)

    var messageParameter string
    printCommand.AddStringParameter(&cli.Parameter{
        Use: "app print --message [value]",
        Name: "message",
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

    err := c.Initialize()
    if err != nil {
        log.Fatal(err)
    }
}
```

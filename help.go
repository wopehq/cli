package cli

import (
	"fmt"
	"os"
)

type Help struct {
	Usage      string
	Commands   string
	Parameters string
}

func (help *Help) ShowHelp() {
	if help != nil {
		fmt.Println(help.Usage + help.Commands + help.Parameters)
		os.Exit(0)
	}
}

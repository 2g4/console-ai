package main

import (
	"github.com/2g4/console-ai/cmd"
	"github.com/2g4/console-ai/data"
)

func main() {
	data.OpenDatabase()
	cmd.InitIfRequired()
	cmd.Execute()
}

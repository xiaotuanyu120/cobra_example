package main

import (
	"github.com/xiaotuanyu120/cobra_example/cmd"
)

var Version = "0.1.1"

func main() {
	cmd.Version = Version
	cmd.Execute()
}

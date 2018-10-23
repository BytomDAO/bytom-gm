package main

import (
	"runtime"

	cmd "github.com/bytom-gm/cmd/bytomcli/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}

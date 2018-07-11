package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"strings"
)

// Simply output any given arguments.
type EchoCmd struct {
}

func (this EchoCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	outChan <- lib.ShellString(strings.Join(arguments, " ") + "\n")
	close(outChan)
}

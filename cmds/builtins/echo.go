package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"strings"
)

// Simply output any given arguments.
type EchoCmd struct {
}

func (this EchoCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	outChan.Write(lib.ShellString(strings.Join(arguments, " ") + "\n"))
	outChan.Close()
}

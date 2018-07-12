package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
)

// Pretty printer.
type PpCmd struct{}

func (this PpCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		outChan.Write(lib.ShellString(fmt.Sprintf("%+v", in)))
	}
	outChan.Close()
}

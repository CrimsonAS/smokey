package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
)

// Pretty printer.
type PpCmd struct{}

func (this PpCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		outChan <- lib.ShellString(fmt.Sprintf("%+v", in))
	}
	close(outChan)
}

package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Send the data of everything from inChan to outChan.
type CatCmd struct{}

func (this CatCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		outChan <- in.Data()
	}
	close(outChan)
}

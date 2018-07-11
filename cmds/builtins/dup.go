package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Duplicate each input to output.
type DupCmd struct {
}

func (this DupCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		outChan <- in
		outChan <- in
	}

	close(outChan)
}

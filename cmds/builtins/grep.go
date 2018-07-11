package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"strings"
)

// Grep the input for an argument to filter by.
type GrepCmd struct{}

func (this GrepCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) == 0 {
		panic("no argument to grep")
	}

	searchStr := arguments[0]

	for in := range inChan {
		if strings.Contains(in.Present(), searchStr) {
			outChan <- in
		}
	}

	close(outChan)
}

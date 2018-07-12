package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"math"
)

// Provides the maximum numeric input.
type MaxCmd struct {
}

func (this MaxCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	max := lib.ShellInt(math.MinInt64)
	maxFound := false
	for in := range inChan {
		if num, ok := in.(lib.ShellInt); ok {
			if num >= max {
				max = num
				maxFound = true
			}
		}
	}
	if maxFound {
		outChan <- lib.ShellInt(max)
	} else {
		outChan <- lib.ShellInt(math.MaxInt64)
	}
	close(outChan)
}

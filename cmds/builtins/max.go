package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"math"
)

// Provides the maximum numeric input.
type MaxCmd struct {
}

func (this MaxCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	max := lib.ShellInt(math.MinInt64)
	maxFound := false
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if num, ok := in.(lib.ShellInt); ok {
			if num >= max {
				max = num
				maxFound = true
			}
		}
	}
	if maxFound {
		outChan.Write(lib.ShellInt(max))
	} else {
		outChan.Write(lib.ShellInt(math.MaxInt64))
	}
	outChan.Close()
}

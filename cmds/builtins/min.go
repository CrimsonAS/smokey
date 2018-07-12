package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"math"
)

// Provides the minimum numeric input.
type MinCmd struct {
}

func (this MinCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	min := lib.ShellInt(math.MaxInt64)
	minFound := false
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if num, ok := in.(lib.ShellInt); ok {
			if num <= min {
				min = num
				minFound = true
			}
		}
	}
	if minFound {
		outChan.Write(lib.ShellInt(min))
	} else {
		outChan.Write(lib.ShellInt(math.MaxInt64))
	}
	outChan.Close()
}

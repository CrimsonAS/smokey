package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"math"
)

// Provides the minimum numeric input.
type MinCmd struct {
}

func (this MinCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	min := lib.ShellInt(math.MaxInt64)
	minFound := false
	for in := range inChan {
		if num, ok := in.(lib.ShellInt); ok {
			if num <= min {
				min = num
				minFound = true
			}
		}
	}
	if minFound {
		outChan <- lib.ShellInt(min)
	} else {
		outChan <- lib.ShellInt(math.MinInt64)
	}
	close(outChan)
}

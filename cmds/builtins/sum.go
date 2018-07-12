package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Sums up how many numeric inputs are provided.
type SumCmd struct {
}

func (this SumCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	count := 0
	for in := range inChan {
		if num, ok := in.(lib.ShellInt); ok {
			count += int(num)
		}
	}
	outChan <- lib.ShellInt(count)
	close(outChan)
}

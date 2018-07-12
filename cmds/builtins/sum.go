package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Sums up how many numeric inputs are provided.
type SumCmd struct {
}

func (this SumCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	count := 0
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if num, ok := in.(lib.ShellInt); ok {
			count += int(num)
		}
	}
	outChan.Write(lib.ShellInt(count))
	outChan.Close()
}

package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Split input into many pieces of output.
// For strings, for instance, this will split on newlines.
type ExplodeCmd struct {
}

func (this ExplodeCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		if explodable, ok := in.(lib.ExplodableData); ok {
			for _, new := range explodable.Explode() {
				outChan <- new
			}
		} else {
			outChan <- in
		}
	}
	close(outChan)
}

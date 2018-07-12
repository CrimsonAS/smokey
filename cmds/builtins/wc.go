package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Count how many items there are in total.
// Items are exploded (if they are explodable), otherwise they are counted as a single item.
type WcCmd struct {
}

func (this WcCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	count := 0
	for in := range inChan {
		if explodable, ok := in.(lib.ExplodableData); ok {
			count += len(explodable.Explode())
		} else {
			count += 1
		}
	}
	outChan <- lib.ShellInt(count)
	close(outChan)
}

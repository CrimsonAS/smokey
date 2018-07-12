package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Count how many items there are in total.
// Items are exploded (if they are explodable), otherwise they are counted as a single item.
type WcCmd struct {
}

func (this WcCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	count := 0
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if explodable, ok := in.(lib.ExplodableData); ok {
			count += len(explodable.Explode())
		} else {
			count += 1
		}
	}
	outChan.Write(lib.ShellInt(count))
	outChan.Close()
}

package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Split input into many pieces of output.
// For strings, for instance, this will split on newlines.
type ExplodeCmd struct {
}

func (this ExplodeCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if explodable, ok := in.(lib.ExplodableData); ok {
			for _, new := range explodable.Explode() {
				outChan.Write(new)
			}
		} else {
			outChan.Write(in)
		}
	}
	outChan.Close()
}

package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Duplicate each input to output.
type DupCmd struct {
}

func (this DupCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		outChan.Write(in)
		outChan.Write(in)
	}

	outChan.Close()
}

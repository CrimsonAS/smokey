package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Send the data of everything from inChan to outChan.
type CatCmd struct{}

func (this CatCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		outChan.Write(in.Data())
	}
	outChan.Close()
}

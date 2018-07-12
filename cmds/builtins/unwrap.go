package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Unwrap a WrappedData
type UnwrapCmd struct{}

func (this UnwrapCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if wd, ok := in.(*lib.WrappedData); ok {
			outChan.Write(wd.RealData)
		} else {
			outChan.Write(in)
		}
	}
	outChan.Close()
}

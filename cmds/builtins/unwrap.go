package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Unwrap a WrappedData
type UnwrapCmd struct{}

func (this UnwrapCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		if wd, ok := in.(*lib.WrappedData); ok {
			outChan <- wd.RealData
		} else {
			outChan <- in
		}
	}
	close(outChan)
}

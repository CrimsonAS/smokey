package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Select a property of input instances that are associative (like a hash or map)
type SpCmd struct{}

func (this SpCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for _, prop := range arguments {
		for in := range inChan {
			if asd, ok := in.(lib.AssociativeShellData); ok {
				outChan <- &lib.WrappedData{RealData: in, FakeData: asd.SelectProperty(prop)}
			}
		}
	}

	close(outChan)
}

package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
)

// Select a property of input instances that are associative (like a hash or map)
type SpCmd struct{}

func (this SpCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for _, prop := range arguments {
		for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
			if asd, ok := in.(lib.AssociativeShellData); ok {
				outChan.Write(&lib.WrappedData{RealData: in, FakeData: asd.SelectProperty(prop)})
			}
		}
	}

	outChan.Close()
}

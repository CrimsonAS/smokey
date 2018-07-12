package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"strconv"
)

// Select a column of input instances that are list-like (arrays etc)
type ScCmd struct{}

func (this ScCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for _, prop := range arguments {
		propInt, err := strconv.Atoi(prop)
		if err != nil {
			panic(fmt.Sprintf("Can't parse col arg %s: %s", prop, err))
		}
		for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
			if asd, ok := in.(lib.ListyShellData); ok {
				outChan.Write(&lib.WrappedData{RealData: in, FakeData: asd.SelectColumn(propInt)})
			}
		}
	}

	outChan.Close()
}

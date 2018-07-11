package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"strconv"
)

// Select a column of input instances that are list-like (arrays etc)
type ScCmd struct{}

func (this ScCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for _, prop := range arguments {
		propInt, err := strconv.Atoi(prop)
		if err != nil {
			panic(fmt.Sprintf("Can't parse col arg %s: %s", prop, err))
		}
		for in := range inChan {
			if asd, ok := in.(lib.ListyShellData); ok {
				outChan <- &lib.WrappedData{RealData: in, FakeData: asd.SelectColumn(propInt)}
			}
		}
	}

	close(outChan)
}

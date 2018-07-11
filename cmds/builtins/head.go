package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"strconv"
)

// Pass the top n pieces of input
type HeadCmd struct{}

func (this HeadCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) == 0 {
		panic("How much do you want?")
	}

	inputLines, err := strconv.Atoi(arguments[0])
	if err != nil {
		panic(fmt.Sprintf("Can't parse head arg %s: %s", arguments[0], err))
	}

	for in := range inChan {
		outChan <- in
		inputLines--

		if inputLines == 0 {
			break
		}
	}

	close(outChan)
}

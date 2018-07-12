package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"strconv"
)

// Pass the top n pieces of input
type HeadCmd struct{}

func (this HeadCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	if len(arguments) == 0 {
		panic("How much do you want?")
	}

	inputLines, err := strconv.Atoi(arguments[0])
	if err != nil {
		panic(fmt.Sprintf("Can't parse head arg %s: %s", arguments[0], err))
	}

	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		outChan.Write(in)
		inputLines--

		if inputLines == 0 {
			break
		}
	}

	outChan.Close()
}

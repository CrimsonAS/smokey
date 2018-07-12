package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"strconv"
)

// Pass the last n pieces of input
type TailCmd struct{}

func (this TailCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	if len(arguments) == 0 {
		panic("How much do you want?")
	}

	inputLines, err := strconv.Atoi(arguments[0])
	if err != nil {
		panic(fmt.Sprintf("Can't parse head arg %s: %s", arguments[0], err))
	}
	last := make([]lib.ShellData, 0, inputLines)

	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		last = append(last, in)
		if len(last) > inputLines {
			last = last[1:]
		}
	}

	for _, in := range last {
		outChan.Write(in)
	}

	outChan.Close()
}

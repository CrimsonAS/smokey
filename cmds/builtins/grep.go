package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"strings"
)

// Grep the input for an argument to filter by.
type GrepCmd struct{}

func (this GrepCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	if len(arguments) == 0 {
		panic("no argument to grep")
	}

	searchStr := arguments[0]

	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if strings.Contains(in.Present(), searchStr) {
			outChan.Write(in)
		}
	}

	outChan.Close()
}

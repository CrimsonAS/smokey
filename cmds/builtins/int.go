package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"strconv"
)

// Creates an integer from arguments.
type IntCmd struct {
}

func (this IntCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for _, arg := range arguments {
		val, err := strconv.Atoi(arg)
		if err != nil {
			panic(fmt.Sprintf("Couldn't convert to int: %s (%s)", arg, err))
		}
		outChan.Write(lib.ShellInt(val))
	}
	outChan.Close()
}

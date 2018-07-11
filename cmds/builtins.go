package cmds

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"strconv"
	"strings"
)

// Simply output any given arguments.
type EchoCmd struct {
}

func (this EchoCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	outChan <- lib.ShellString(strings.Join(arguments, " ") + "\n")
	close(outChan)
}

type ExplodeCmd struct {
}

func (this ExplodeCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		if explodable, ok := in.(lib.ExplodableData); ok {
			for _, new := range explodable.Explode() {
				outChan <- new
			}
		} else {
			outChan <- in
		}
	}
	close(outChan)
}

// Duplicate each input to output.
type DupCmd struct {
}

func (this DupCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		outChan <- in
		outChan <- in
	}

	close(outChan)
}

// Remove duplicates from input, write the ordered (but unique) inputs to output.
type UniqCmd struct {
}

func (this UniqCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	dat := make(map[interface{}]lib.ShellData, 1024)
	orderedDat := make([]lib.ShellData, 0, 1024)

	for in := range inChan {
		_, ok := dat[in]
		if !ok {
			dat[in] = in
			orderedDat = append(orderedDat, in)
		}
	}

	for _, out := range orderedDat {
		outChan <- out
	}
	close(outChan)
}

// Send the data of everything from inChan to outChan.
type CatCmd struct{}

func (this CatCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		outChan <- in.Data()
	}
	close(outChan)
}

// Grep the input for an argument to filter by.
type GrepCmd struct{}

func (this GrepCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) == 0 {
		panic("no argument to grep")
	}

	searchStr := arguments[0]

	for in := range inChan {
		if strings.Contains(in.Present(), searchStr) {
			outChan <- in
		}
	}

	close(outChan)
}

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

// Pass the last n pieces of input
type TailCmd struct{}

func (this TailCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) == 0 {
		panic("How much do you want?")
	}

	inputLines, err := strconv.Atoi(arguments[0])
	if err != nil {
		panic(fmt.Sprintf("Can't parse head arg %s: %s", arguments[0], err))
	}
	last := make([]lib.ShellData, 0, inputLines)

	for in := range inChan {
		last = append(last, in)
		if len(last) > inputLines {
			last = last[1:]
		}
	}

	for _, in := range last {
		outChan <- in
	}

	close(outChan)
}

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

// Pretty printer.
type PpCmd struct{}

func (this PpCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for in := range inChan {
		outChan <- lib.ShellString(fmt.Sprintf("%+v", in))
	}
	close(outChan)
}

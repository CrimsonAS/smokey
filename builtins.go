package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Simply output any given arguments.
type EchoCmd struct {
}

func (this EchoCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	outChan <- shellString(strings.Join(arguments, " ") + "\n")
	close(outChan)
}

// Split input on \n
// ### not very generic
type LinesCmd struct {
}

func (this LinesCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	for in := range inChan {
		dat := in.Data()
		splitz := bytes.Split(dat, []byte{'\n'})
		for _, line := range splitz {
			outChan <- shellString(line)
		}
	}
	close(outChan)
}

// Duplicate each input to output.
type DupCmd struct {
}

func (this DupCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	for in := range inChan {
		outChan <- in
		outChan <- in
	}

	close(outChan)
}

// Remove duplicates from input, write the ordered (but unique) inputs to output.
type UniqCmd struct {
}

func (this UniqCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	dat := make(map[interface{}]shellData, 1024)
	orderedDat := make([]shellData, 0, 1024)

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

func (this CatCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	for in := range inChan {
		outChan <- in.Data()
	}
	close(outChan)
}

// Grep the input for an argument to filter by.
type GrepCmd struct{}

func (this GrepCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	if len(arguments) == 0 {
		panic("no argument to grep")
	}

	searchStr := arguments[0]

	for in := range inChan {
		if in.Grep(searchStr) {
			outChan <- in
		}
	}

	close(outChan)
}

// Pass the top n pieces of input
type HeadCmd struct{}

func (this HeadCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
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

func (this TailCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	if len(arguments) == 0 {
		panic("How much do you want?")
	}

	inputLines, err := strconv.Atoi(arguments[0])
	if err != nil {
		panic(fmt.Sprintf("Can't parse head arg %s: %s", arguments[0], err))
	}
	last := make([]shellData, 0, inputLines)

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

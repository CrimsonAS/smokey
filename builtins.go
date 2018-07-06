package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Simply output any given arguments.
type EchoCmd struct {
}

func (this EchoCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	outChan <- shellString(strings.Join(arguments, " ") + "\n")
	close(outChan)
}

// Send the data of everything from inChan to outChan.
type CatCmd struct {
}

func (this CatCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	for in := range inChan {
		outChan <- in.Data()
	}
	close(outChan)
}

// List the current directory.
type LsCmd struct {
}

func (this LsCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	if len(arguments) == 0 {
		arguments = []string{"."}
	}

	for _, dir := range arguments {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			outChan <- &shellPath{pathName: file.Name()}
		}
	}

	close(outChan)
}

// Grep the input for an argument to filter by.
type GrepCmd struct {
}

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

// Change the current working directory.
type CdCmd struct {
}

func (this CdCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	if len(arguments) == 0 {
		arguments = []string{os.Getenv("HOME")}
	}

	os.Chdir(arguments[0])
	close(outChan)
}

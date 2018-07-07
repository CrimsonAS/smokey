package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// A pathname for something on disk
type shellPath struct {
	pathName string
}

func (this *shellPath) Grep(searchable string) bool {
	return strings.Contains(this.pathName, searchable)
}

func (this *shellPath) Data() shellBuffer {
	data, err := ioutil.ReadFile(this.pathName)
	if err != nil {
		panic(fmt.Sprintf("Can't read file %s: %s", this.pathName, err))
	}
	return shellBuffer(data)
}

func (this *shellPath) Present() string {
	return fmt.Sprintf("%s\n", this.pathName)
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

// List the current directory.
type LsCmd struct {
}

func (this LsCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	if len(arguments) == 0 {
		arguments = []string{"."}
	}

	for _, path := range arguments {
		finfo, err := os.Stat(path)
		if err != nil {
			panic(fmt.Sprintf("Can't stat path %s: %s", path, err))
		}

		if finfo.Mode().IsDir() {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				panic(fmt.Sprintf("Can't read directory %s: %s", path, err))
			}

			for _, file := range files {
				outChan <- &shellPath{pathName: file.Name()}
			}
		} else {
			outChan <- &shellPath{pathName: path}
		}
	}

	close(outChan)
}

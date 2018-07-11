package cmds

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"io/ioutil"
	"os"
)

// A pathname for something on disk
type shellPath struct {
	pathName string
}

func (this *shellPath) Data() lib.ShellBuffer {
	data, err := ioutil.ReadFile(this.pathName)
	if err != nil {
		panic(fmt.Sprintf("Can't read file %s: %s", this.pathName, err))
	}
	return lib.ShellBuffer(data)
}

func (this *shellPath) Present() string {
	return fmt.Sprintf("%s\n", this.pathName)
}

// Change the current working directory.
type CdCmd struct {
}

func (this CdCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) == 0 {
		arguments = []string{os.Getenv("HOME")}
	}

	os.Chdir(arguments[0])
	close(outChan)
}

// List the current directory.
type LsCmd struct {
}

func (this LsCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
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

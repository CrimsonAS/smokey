package cmds

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"os"
)

// Change the current working directory.
type CdCmd struct {
}

func (this CdCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	if len(arguments) == 0 {
		arguments = []string{os.Getenv("HOME")}
	}

	os.Chdir(arguments[0])
	outChan.Close()
}

// List the current directory.
type LsCmd struct {
}

func (this LsCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	if len(arguments) == 0 {
		arguments = []string{"."}
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Can't find current directory: %s", err))
	}

	for _, path := range arguments {
		finfo, err := os.Stat("./" + path)
		if err != nil {
			panic(fmt.Sprintf("Can't stat path %s: %s", path, err))
		}

		sp := &lib.ShellPath{RootPath: cwd, FileName: finfo.Name(), IsDir: finfo.IsDir()}
		if sp.IsDir {
			contents := sp.Explode()

			for _, file := range contents {
				outChan.Write(file)
			}
		} else {
			outChan.Write(sp)
		}
	}

	outChan.Close()
}

package cmds

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"io/ioutil"
	"os"
)

// Data representing something on disk.
type shellPath struct {
	rootPath string
	fi       os.FileInfo
}

func (this *shellPath) fullPath() string {
	if this.rootPath != "./." {
		return fmt.Sprintf("%s/%s", this.rootPath, this.fi.Name())
	}
	return this.fi.Name()
}

func (this *shellPath) isDir() bool {
	return this.fi.IsDir()
}

func (this *shellPath) Data() lib.ShellBuffer {
	path := this.fullPath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Can't read file %s: %s", path, err))
	}
	return lib.ShellBuffer(data)
}

func (this *shellPath) Present() string {
	path := this.fullPath()
	return fmt.Sprintf("%s\n", path)
}

func (this *shellPath) Explode() []lib.ShellData {
	if this.isDir() {
		fp := this.fullPath()
		files, err := ioutil.ReadDir(fp)
		if err != nil {
			panic(fmt.Sprintf("Can't read directory %s: %s", fp, err))
		}

		ret := make([]lib.ShellData, len(files))

		for idx, fi := range files {
			ret[idx] = &shellPath{rootPath: fp, fi: fi}
		}

		return ret
	}

	return nil
}

func (this *shellPath) SelectProperty(prop string) lib.ShellData {
	if prop == "mtime" {
		return lib.ShellTime(this.fi.ModTime())
	}
	return nil
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
		sp := &shellPath{rootPath: path, fi: finfo}
		if sp.isDir() {
			contents := sp.Explode()

			for _, file := range contents {
				outChan <- file
			}
		} else {
			outChan <- sp
		}
	}

	close(outChan)
}

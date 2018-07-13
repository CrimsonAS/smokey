package lib

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Data representing something on disk.
type ShellPath struct {
	RootPath string
	FileName string
	IsDir    bool
}

func (this *ShellPath) fullPath() string {
	if this.RootPath != "./." {
		return fmt.Sprintf("%s/%s", this.RootPath, this.FileName)
	}
	return this.FileName
}

func (this *ShellPath) Data() ShellBuffer {
	path := this.fullPath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Can't read file %s: %s", path, err))
	}
	return ShellBuffer(data)
}

func (this *ShellPath) Present() string {
	path := this.fullPath() // ### unfortunate this prints an absolute path if it's in CWD...
	return fmt.Sprintf("%s\n", path)
}

func (this *ShellPath) Explode() []ShellData {
	if this.IsDir {
		fp := this.fullPath()
		files, err := ioutil.ReadDir(fp)
		if err != nil {
			panic(fmt.Sprintf("Can't read directory %s: %s", fp, err))
		}

		ret := make([]ShellData, len(files))

		for idx, fi := range files {
			ret[idx] = &ShellPath{RootPath: fp, FileName: fi.Name(), IsDir: fi.IsDir()}
		}

		return ret
	}

	return nil
}

func (this *ShellPath) SelectProperty(prop string) ShellData {
	if prop == "mtime" {
		fi, _ := os.Stat(this.fullPath()) // ### stat caching might be an idea
		return ShellTime(fi.ModTime())
	}
	return nil
}

package lib

import (
	"fmt"
	"strings"
)

// An arbitrary string
type ShellString string

func (this ShellString) Data() ShellBuffer {
	return ShellBuffer(this)
}

func (this ShellString) Present() string {
	return fmt.Sprintf("%s\n", this)
}

func (this ShellString) Explode() []ShellData {
	substrings := strings.Split(string(this), "\n")
	ret := make([]ShellData, len(substrings))
	for idx, str := range substrings {
		ret[idx] = ShellString(str)
	}
	return ret
}

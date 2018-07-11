package lib

import (
	"bytes"
	"fmt"
)

// An arbitrary series of bytes
type ShellBuffer []byte

func (this ShellBuffer) Grep(searchable string) bool {
	searchableBytes := []byte(searchable)
	return bytes.Contains(this, searchableBytes)
}

func (this ShellBuffer) Data() ShellBuffer {
	return this
}

func (this ShellBuffer) Present() string {
	return fmt.Sprintf("%s\n", this)
}

func (this ShellBuffer) Explode() []ShellData {
	substrings := bytes.Split(this, []byte{'\n'})
	ret := make([]ShellData, len(substrings))
	for idx, str := range substrings {
		ret[idx] = ShellBuffer(str)
	}
	return ret
}

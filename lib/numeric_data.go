package lib

import (
	"fmt"
)

// An integer
type ShellInt int

func (this ShellInt) Data() ShellBuffer {
	return ShellBuffer(this.Present())
}

func (this ShellInt) Present() string {
	return fmt.Sprintf("%d\n", this)
}

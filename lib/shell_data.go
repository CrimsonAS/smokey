package lib

import (
	"bytes"
	"fmt"
	"strings"
)

// Data passed in a shell pipeline.
type ShellData interface {
	// Does this data contain this string?
	Grep(searchable string) bool

	// Return the underlying data for this shell item
	// For example, for a file, return the contents
	Data() ShellBuffer

	// Present the shell object in a nice way for display
	Present() string
}

type ListyShellData interface {
	SelectColumn(col int) ShellData
}

type AssociativeShellData interface {
	SelectProperty(property string) ShellData
}

// An arbitrary string
type ShellString string

func (this ShellString) Grep(searchable string) bool {
	return strings.Contains(string(this), searchable)
}

func (this ShellString) Data() ShellBuffer {
	return ShellBuffer(this)
}

func (this ShellString) Present() string {
	return fmt.Sprintf("%s\n", this)
}

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

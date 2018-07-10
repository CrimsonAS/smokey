package main

import (
	"bytes"
	"fmt"
	"strings"
)

// Data passed in a shell pipeline.
type shellData interface {
	// Does this data contain this string?
	Grep(searchable string) bool

	// Return the underlying data for this shell item
	// For example, for a file, return the contents
	Data() shellBuffer

	// Present the shell object in a nice way for display
	Present() string
}

type listyShellData interface {
	SelectColumn(col int) shellData
}

type associativeShellData interface {
	SelectProperty(property string) shellData
}

// An arbitrary string
type shellString string

func (this shellString) Grep(searchable string) bool {
	return strings.Contains(string(this), searchable)
}

func (this shellString) Data() shellBuffer {
	return shellBuffer(this)
}

func (this shellString) Present() string {
	return fmt.Sprintf("%s\n", this)
}

// An arbitrary series of bytes
type shellBuffer []byte

func (this shellBuffer) Grep(searchable string) bool {
	searchableBytes := []byte(searchable)
	return bytes.Contains(this, searchableBytes)
}

func (this shellBuffer) Data() shellBuffer {
	return this
}

func (this shellBuffer) Present() string {
	return fmt.Sprintf("%s\n", this)
}

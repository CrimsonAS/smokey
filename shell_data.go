package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

// A pathname for a file.
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
	return fmt.Sprintf("file://%s\n", this.pathName)
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

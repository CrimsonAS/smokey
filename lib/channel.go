package lib

import (
	"fmt"
)

type Channel struct {
	ch chan ShellData
}

func (this *Channel) Handle() string {
	return fmt.Sprintf("%p", this.ch)
}

func NewChannel() *Channel {
	c := &Channel{
		ch: make(chan ShellData),
	}
	return c
}

func (this *Channel) Write(out ShellData) bool {
	this.ch <- out
	return true
}

func (this *Channel) Read() (ShellData, bool) {
	val := <-this.ch

	if val == nil {
		return nil, false
	} else {
		return val, true
	}
}

func (this *Channel) Close() {
	close(this.ch)
}

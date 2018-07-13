package lib

import (
	"fmt"
	"net/url"
)

// A description of a connection to a remote resource
type ShellUrl struct {
	*url.URL
}

func NewUrl(rawurl string) (ShellUrl, error) {
	u, p := url.Parse(rawurl)
	return ShellUrl{u}, p
}

func (this *ShellUrl) Data() ShellBuffer {
	return ShellBuffer(this.Present())
}

func (this *ShellUrl) Present() string {
	return fmt.Sprintf("%s\n", this.URL.String())
}

package web

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"io/ioutil"
	"net/http"
	"strings"
)

// A URI representing a remote resource
type shellUri struct {
	uri string
}

func (this *shellUri) Data() lib.ShellBuffer {
	resp, err := http.Get(this.uri)
	if err != nil {
		panic(fmt.Sprintf("Can't read URI %s: %s", this.uri, err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Can't read URI body %s: %s", this.uri, err))
	}

	return lib.ShellBuffer(body)
}

func (this *shellUri) Explode() []lib.ShellData {
	return this.Data().Explode()
}

func (this *shellUri) Present() string {
	return fmt.Sprintf("%s\n", this.uri)
}

// Turn arguments into URI
type FetchCmd struct{}

func (this FetchCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	if len(arguments) == 0 {
		panic("What do you want to fetch?")
	}

	for _, uri := range arguments {
		// Kind of artificial limitation. Should not modify another scheme...
		// Handle this better ###
		if !strings.HasPrefix(uri, "https://") && !strings.HasPrefix(uri, "http://") {
			uri = "https://" + uri
		}
		outChan <- &shellUri{uri}
	}

	close(outChan)
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// A URI representing a remote resource
type shellUri struct {
	uri string
}

func (this *shellUri) Grep(searchable string) bool {
	return strings.Contains(this.uri, searchable)
}

func (this *shellUri) Data() shellBuffer {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		panic(fmt.Sprintf("Can't read URI %s: %s", this.uri, err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Can't read URI body %s: %s", this.uri, err))
	}

	return shellBuffer(body)
}

func (this *shellUri) Present() string {
	return fmt.Sprintf("%s\n", this.uri)
}

// Turn arguments into URI
type FetchCmd struct{}

func (this FetchCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	if len(arguments) == 0 {
		panic("What do you want to fetch?")
	}

	for _, uri := range arguments {
		// Kind of artificial limitation. Need a factory function?
		if !strings.HasPrefix(uri, "https://") && !strings.HasPrefix(uri, "http://") {
			panic(fmt.Sprintf("URI %s not supported", uri))
		}
		log.Printf("Creating shellUri for %s", uri)
		outChan <- &shellUri{uri}
	}

	close(outChan)
}

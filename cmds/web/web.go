package web

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"io/ioutil"
	"net/http"
)

// Fetch URLs
type FetchCmd struct{}

func (this FetchCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		if url, isUrl := in.(*lib.ShellUrl); isUrl {
			switch url.Scheme {
			case "":
				panic(fmt.Sprintf("Can't assume a scheme for fetch on %s", url))
			case "http":
				fallthrough
			case "https":
				resp, err := http.Get(url.String())
				if err != nil {
					panic(fmt.Sprintf("Can't read URL %s: %s", url, err))
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(fmt.Sprintf("Can't read URL body %s: %s", url, err))
				}
				resp.Body.Close()
				outChan.Write(lib.ShellBuffer(body))
			}
		}
	}

	outChan.Close()
}

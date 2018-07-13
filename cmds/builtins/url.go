package builtins

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
)

// Create a URL
type UrlCmd struct{}

func (this UrlCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for _, arg := range arguments {
		u, err := lib.NewUrl(arg)
		if err != nil {
			panic(fmt.Sprintf("Bad URL: %s", err))
		}
		outChan.Write(&u)
	}
	outChan.Close()
}

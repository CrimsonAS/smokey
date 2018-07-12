package main

import (
	"github.com/CrimsonAS/smokey/lib"
)

func main() {
	ctx := lib.DialPlugin()
	inChan := ctx.InChan
	outChan := ctx.OutChan

	outChan.Write(lib.ShellString("Hello world, from a plugin"))

	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		outChan.Write(in)
	}

	outChan.Close()
	ctx.Wait()
}

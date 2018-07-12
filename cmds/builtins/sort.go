package builtins

import (
	"github.com/CrimsonAS/smokey/lib"
	"sort"
)

// Sort the input.
type SortCmd struct{}

type arbitraryShellData []lib.ShellData

func (this arbitraryShellData) Len() int {
	return len(this)
}
func (this arbitraryShellData) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this arbitraryShellData) Less(i, j int) bool {
	r := this[i]
	l := this[j]
	return r.Present() < l.Present()
}

func (this SortCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	sortData := arbitraryShellData{}
	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		sortData = append(sortData, in)
	}
	sort.Sort(sortData)

	for _, out := range sortData {
		outChan.Write(out)
	}
	outChan.Close()
}

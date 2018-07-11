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

func (this SortCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	sortData := arbitraryShellData{}
	for in := range inChan {
		sortData = append(sortData, in)
	}
	sort.Sort(sortData)

	for _, out := range sortData {
		outChan <- out
	}
	close(outChan)
}

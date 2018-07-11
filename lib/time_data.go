package lib

import (
	"fmt"
	"log"
	"time"
)

// A date and time
type ShellTime time.Time

func (this ShellTime) Data() ShellBuffer {
	return ShellBuffer(this.Present())
}

func (this ShellTime) Present() string {
	return fmt.Sprintf("%s\n", time.Time(this).String())
}

func (this ShellTime) SelectProperty(prop string) ShellData {
	switch prop {
	case "day":
		return ShellString(fmt.Sprintf("%d", time.Time(this).Day()))
	case "month":
		return ShellString(fmt.Sprintf("%d", time.Time(this).Month()))
	case "year":
		return ShellString(fmt.Sprintf("%d", time.Time(this).Year()))
	case "hour":
		return ShellString(fmt.Sprintf("%d", time.Time(this).Hour()))
	case "minute":
		return ShellString(fmt.Sprintf("%d", time.Time(this).Minute()))
	case "second":
		return ShellString(fmt.Sprintf("%d", time.Time(this).Second()))
	}

	log.Printf("Whoops?")
	return nil
}

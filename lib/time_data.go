package lib

import (
	"fmt"
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
		return ShellInt(time.Time(this).Day())
	case "month":
		return ShellInt(time.Time(this).Month())
	case "year":
		return ShellInt(time.Time(this).Year())
	case "hour":
		return ShellInt(time.Time(this).Hour())
	case "minute":
		return ShellInt(time.Time(this).Minute())
	case "second":
		return ShellInt(time.Time(this).Second())
	}

	return nil
}

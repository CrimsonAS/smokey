package cmds

import (
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"github.com/shirou/gopsutil/process"
	"syscall"
)

// A process object
type shellProcess struct {
	pid int32
}

func (this *shellProcess) Data() lib.ShellBuffer {
	return lib.ShellBuffer(this.Present())
}

func (this *shellProcess) Present() string {
	p, err := process.NewProcess(this.pid)
	if err != nil {
		return fmt.Sprintf("PID %d\n", this.pid)
	}
	n, err := p.Name()
	if err != nil {
		return fmt.Sprintf("PID %d\n", this.pid)
	}
	return fmt.Sprintf("PID %d (%s)\n", this.pid, n)
}

// Turn arguments into URI
type PsCmd struct{}

func (this PsCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	pidList, err := process.Pids()
	if err != nil {
		panic(fmt.Sprintf("Can't fetch PIDs: %s", err))
	}

	for _, pid := range pidList {
		outChan <- &shellProcess{pid: pid}
	}

	close(outChan)
}

// Kill a process
type KillCmd struct{}

func (this KillCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for pid := range inChan {
		proc := pid.(*shellProcess)
		syscall.Kill(int(proc.pid), syscall.SIGTERM)
	}
	close(outChan)
}

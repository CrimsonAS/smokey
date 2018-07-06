package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
)

// A command wrapping an OS process (and its stdin/stdout/stderr).
type StandardProcessCmd struct {
	// The process path to run.
	process string
}

func (this StandardProcessCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	cmd := exec.Command(this.process, arguments...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(fmt.Sprintf("Couldn't get stderr pipe: %s", err))
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("Couldn't get stdout pipe: %s", err))
	}

	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("error starting %s", err))
	}

	// ### handle stdin too

	go func() {
		errOut, err := ioutil.ReadAll(stderr)
		if err != nil {
			panic(fmt.Sprintf("stderr read failed: %s", err))
		}
		if len(errOut) > 0 {
			panic(fmt.Sprintf("process error: %s", errOut))
		}
	}()

	go func() {
		reader := bufio.NewReader(stdout)
		stdBuf := make([]byte, 4096)
		for {
			_, err := reader.Read(stdBuf)
			if err == io.EOF {
				return
			}
			if err != nil {
				panic(fmt.Sprintf("stdout read failed: %s", err))
			}
			outChan <- shellBuffer(stdBuf)
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	close(outChan)
}

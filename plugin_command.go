package main

import (
	"bufio"
	"fmt"
	"github.com/CrimsonAS/smokey/lib"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
)

// Experimental. Runs an out-of-process plugin.
type PcCmd struct {
}

const pcCmdDebug = false
const subprocDebug = false

func canIgnoreError(err error) bool {
	if strings.Contains(err.Error(), "file already closed") {
		return true
	}
	return false
}

func (this PcCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	ln, err := net.Listen("tcp", "localhost:0")
	defer ln.Close()
	if err != nil {
		panic(fmt.Sprintf("Plugin command can't communicate: %s", err))
	}

	pluginArgs := []string{ln.Addr().String()}
	if len(arguments) > 1 {
		pluginArgs = append(pluginArgs, arguments[1:]...)
	}
	pluginArgs = append(pluginArgs, pluginArgs...)
	cmd := exec.Command("plugintest/plugintest", pluginArgs...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("Couldn't get stdout pipe: %s", err))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(fmt.Sprintf("Couldn't get stderr pipe: %s", err))
	}
	consumer := func(from io.Reader) {
		breader := bufio.NewReader(from)
		for {
			oneByte, err := breader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else if !canIgnoreError(err) {
					panic(fmt.Sprintf("Stdout/stderr read failure: %s", err))
				}
			}
			if subprocDebug && len(oneByte) > 0 {
				log.Printf("STDERR: %s", oneByte)
			}
		}
	}
	go consumer(stdout)
	go consumer(stderr)

	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("error starting %s", err))
	}

	conn, err := ln.Accept()
	if err != nil {
		panic(fmt.Sprintf("Can't accept plugin connection: %s", err))
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// Send from inChan to plugin
	// Read from plugin to outChan
	go func() {
		lib.WriteFromChannel(conn, inChan)
		wg.Done()
	}()
	go func() {
		lib.ReadToChannel(conn, outChan)
		wg.Done()
	}()

	wg.Wait()
	conn.Close()

	if err := cmd.Wait(); err != nil {
		panic(fmt.Sprintf("error waiting %s", err))
	}
	if pcCmdDebug {
		log.Printf("All done")
	}

}

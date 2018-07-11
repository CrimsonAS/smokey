package main

import (
	"bufio"
	"fmt"
	"github.com/CrimsonAS/smokey/cmds"
	"github.com/CrimsonAS/smokey/cmds/influx"
	"github.com/CrimsonAS/smokey/cmds/web"
	"github.com/CrimsonAS/smokey/lib"
	"os"
	"strings"
)

// All commands implement this interface.
type commandObject interface {
	// Call the comand. The inChan and outChan are used for communication.
	// The arguments let it customize its behaviour from the command line.
	Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string)
}

// Parse and execute a given command pipeline.
func runCommandString(text string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic", r)
		}
	}()
	commands := parsePipeline(text)
	var inChan chan lib.ShellData
	var outChan chan lib.ShellData

	inChan = make(chan lib.ShellData)
	outChan = make(chan lib.ShellData)
	close(inChan) // ### not what we should do really

	for idx, cmd := range commands {
		var commandObject commandObject
		switch cmd.Command {
		case "sc":
			commandObject = cmds.ScCmd{}
		case "sp":
			commandObject = cmds.SpCmd{}
		case "unwrap":
			commandObject = cmds.UnwrapCmd{}
		case "influxConnect":
			commandObject = influx.InfluxConnect{}
		case "influxQuery":
			commandObject = influx.InfluxQuery{}
		case "ps":
			commandObject = cmds.PsCmd{}
		case "kill":
			commandObject = cmds.KillCmd{}
		case "head":
			commandObject = cmds.HeadCmd{}
		case "tail":
			commandObject = cmds.TailCmd{}
		case "echo":
			commandObject = cmds.EchoCmd{}
		case "cat":
			commandObject = cmds.CatCmd{}
		case "dup":
			commandObject = cmds.DupCmd{}
		case "uniq":
			commandObject = cmds.UniqCmd{}
		case "explode":
			commandObject = cmds.ExplodeCmd{}
		case "last":
			commandObject = LastCmd{}
		case "ls":
			commandObject = cmds.LsCmd{}
		case "cd":
			commandObject = cmds.CdCmd{}
		case "fetch":
			commandObject = web.FetchCmd{}
		case "grep":
			commandObject = cmds.GrepCmd{}
		case "pp":
			commandObject = cmds.PpCmd{}
		default:
			commandObject = cmds.StandardProcessCmd{Process: cmd.Command}
		}
		go commandObject.Call(inChan, outChan, cmd.Arguments)

		inChan = outChan
		if idx < len(commands)-1 {
			outChan = make(chan lib.ShellData)
		}
	}

	present(outChan)
}

// LastCmd just repeats whatever shell data the last command pipeline produced.
// This doesn't really belong here, but right now it hacks present(), so it's here
// for easy reference.
type LastCmd struct {
}

func (this LastCmd) Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string) {
	for _, last := range lastOut {
		outChan <- last
	}
	close(outChan)
}

// ### used by last command. ideally, it would somehow keep hold of this itself
var lastOut []lib.ShellData

// Wait for a command to finish, presenting data as it arrives.
func present(outChan chan lib.ShellData) {
	var newOut []lib.ShellData
	const presentDebug = false
	for res := range outChan {
		newOut = append(newOut, res)
		if presentDebug {
			fmt.Printf("%T: %s", res, res.Present())
		} else {
			fmt.Printf("%s", res.Present())
		}
	}
	lastOut = newOut
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("smokey the shell")
	fmt.Println("try something like echo hello my friend | cat")
	fmt.Println("---------------------")

	for {
		fmt.Print("% ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		runCommandString(text)
	}
}

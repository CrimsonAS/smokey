package main

import (
	"bufio"
	"fmt"
	"github.com/CrimsonAS/smokey/cmds"
	"github.com/CrimsonAS/smokey/cmds/builtins"
	"github.com/CrimsonAS/smokey/cmds/influx"
	"github.com/CrimsonAS/smokey/cmds/web"
	"github.com/CrimsonAS/smokey/lib"
	"os"
	"reflect"
	"strings"
	"sync"
)

// All commands implement this interface.
type commandObject interface {
	// Call the comand. The inChan and outChan are used for communication.
	// The arguments let it customize its behaviour from the command line.
	Call(inChan chan lib.ShellData, outChan chan lib.ShellData, arguments []string)
}

func runCommand(cmd Command, inChan chan lib.ShellData) chan lib.ShellData {
	var commandObject commandObject
	switch cmd.Command {
	case "sc":
		commandObject = builtins.ScCmd{}
	case "sp":
		commandObject = builtins.SpCmd{}
	case "unwrap":
		commandObject = builtins.UnwrapCmd{}
	case "influxConnect":
		commandObject = influx.InfluxConnect{}
	case "influxQuery":
		commandObject = influx.InfluxQuery{}
	case "ps":
		commandObject = cmds.PsCmd{}
	case "kill":
		commandObject = cmds.KillCmd{}
	case "head":
		commandObject = builtins.HeadCmd{}
	case "tail":
		commandObject = builtins.TailCmd{}
	case "echo":
		commandObject = builtins.EchoCmd{}
	case "cat":
		commandObject = builtins.CatCmd{}
	case "dup":
		commandObject = builtins.DupCmd{}
	case "uniq":
		commandObject = builtins.UniqCmd{}
	case "explode":
		commandObject = builtins.ExplodeCmd{}
	case "last":
		commandObject = LastCmd{}
	case "ls":
		commandObject = cmds.LsCmd{}
	case "cd":
		commandObject = cmds.CdCmd{}
	case "fetch":
		commandObject = web.FetchCmd{}
	case "grep":
		commandObject = builtins.GrepCmd{}
	case "sort":
		commandObject = builtins.SortCmd{}
	case "pp":
		commandObject = builtins.PpCmd{}
	case "wc":
		commandObject = builtins.WcCmd{}
	default:
		commandObject = cmds.StandardProcessCmd{Process: cmd.Command}
	}
	outChan := make(chan lib.ShellData)
	go commandObject.Call(inChan, outChan, cmd.Arguments)
	return outChan
}

// Parse and execute a given command pipeline.
func runCommandString(text string) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("Recovered from panic", r)
	//	}
	//}()
	nodes := parsePipeline(text)
	var inChan chan lib.ShellData

	inChan = make(chan lib.ShellData)
	close(inChan) // ### not what we should do really

	for _, node := range nodes {
		//log.Printf("Running node %+v", node)
		inChan = runNode(inChan, node)
		//log.Printf("Done running node %+v", node)
	}

	present(inChan)
}

func runUnion(node UnionNode, inChan chan lib.ShellData) chan lib.ShellData {
	leftOutChan := runNode(inChan, node.Left)
	rightOutChan := runNode(inChan, node.Right)
	unionChan := make(chan lib.ShellData)

	var wg sync.WaitGroup
	wg.Add(2)

	pipeChan := func(readFrom, writeTo chan lib.ShellData) {
		for in := range readFrom {
			writeTo <- in
		}
		wg.Done()
	}

	go pipeChan(leftOutChan, unionChan)
	go pipeChan(rightOutChan, unionChan)

	// In the normal command pipeline case, the command is expected to close
	// its output once it is finished writing. However, an operator node is
	// not quite using that path. In this case, we have two channels (for
	// lhs and rhs), and both will close normally. The union channel,
	// though, we have to close _ourselves_.
	//
	// To do this in a way that won't block further processing of the
	// pipeline, we use a WaitGroup plus an asynchronous goroutine below.
	go func() {
		wg.Wait()
		close(unionChan)
	}()

	return unionChan
}

func runDifference(node DifferenceNode, inChan chan lib.ShellData) chan lib.ShellData {
	leftOutChan := runNode(inChan, node.Left)
	rightOutChan := runNode(inChan, node.Right)

	leftData := []lib.ShellData{}
	rightData := []lib.ShellData{}

	var wg sync.WaitGroup
	wg.Add(2)

	pipeChan := func(outChan chan lib.ShellData, data *[]lib.ShellData) {
		for newData := range outChan {
			*data = append(*data, newData)
		}
		wg.Done()
	}

	go pipeChan(leftOutChan, &leftData)
	go pipeChan(rightOutChan, &rightData)

	// We can't do this asynchronously, because we can't produce data until we
	// have lhs and rhs data.
	wg.Wait()

	// But we must stream our data out asynchronously.
	differenceChan := make(chan lib.ShellData)
	go func() {
		// ### this is O(n^n) and it really shouldn't be.
		for _, left := range leftData {
			found := false
			for _, right := range rightData {
				if reflect.DeepEqual(left, right) {
					found = true
					break
				}
			}

			if !found {
				differenceChan <- left
			}
		}
		close(differenceChan)
	}()

	return differenceChan
}

func runNode(inChan chan lib.ShellData, node PipelineNode) chan lib.ShellData {
	switch typedNode := node.(type) {
	case Command:
		return runCommand(typedNode, inChan)
	case UnionNode:
		return runUnion(typedNode, inChan)
	case DifferenceNode:
		return runDifference(typedNode, inChan)
	default:
		panic("Unsupported node")
	}
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

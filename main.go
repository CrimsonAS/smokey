package main

import (
	"bufio"
	"fmt"
	"github.com/CrimsonAS/smokey/cmds"
	"github.com/CrimsonAS/smokey/cmds/builtins"
	"github.com/CrimsonAS/smokey/cmds/influx"
	"github.com/CrimsonAS/smokey/cmds/web"
	"github.com/CrimsonAS/smokey/lib"
	"log"
	"os"
	"reflect"
	"runtime/trace"
	"strings"
	"sync"
)

type commandObject interface {
	// Call the comand. The inChan and outChan are used for communication.
	// The arguments let it customize its behaviour from the command line.
	Call(inChan, outChan *lib.Channel, arguments []string)
}

func runCommand(cmd Command, inChan *lib.Channel) *lib.Channel {
	var commandObject commandObject
	switch cmd.Command {
	case "sc":
		commandObject = builtins.ScCmd{}
	case "sp":
		commandObject = builtins.SpCmd{}
	case "unwrap":
		commandObject = builtins.UnwrapCmd{}
	case "url":
		commandObject = builtins.UrlCmd{}
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
	case "int":
		commandObject = builtins.IntCmd{}
	case "sum":
		commandObject = builtins.SumCmd{}
	case "min":
		commandObject = builtins.MinCmd{}
	case "max":
		commandObject = builtins.MaxCmd{}
	case "exit":
		exited = true
		commandObject = builtins.EchoCmd{}
	case "pc":
		commandObject = PcCmd{}
	default:
		commandObject = cmds.StandardProcessCmd{Process: cmd.Command}
	}
	outChan := lib.NewChannel()
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
	inChan := lib.NewChannel()
	inChan.Close() // ### not what we should do really

	for _, node := range nodes {
		//log.Printf("Running node %+v", node)
		inChan = runNode(inChan, node)
		//log.Printf("Done running node %+v", node)
	}

	present(inChan)
}

func runUnion(node UnionNode, inChan *lib.Channel) *lib.Channel {
	leftOutChan := runNode(inChan, node.Left)
	rightOutChan := runNode(inChan, node.Right)
	unionChan := lib.NewChannel()

	var wg sync.WaitGroup
	wg.Add(2)

	pipeChan := func(readFrom, writeTo *lib.Channel) {
		for in, ok := readFrom.Read(); ok; in, ok = readFrom.Read() {
			writeTo.Write(in)
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
		unionChan.Close()
	}()

	return unionChan
}

func runDifference(node DifferenceNode, inChan *lib.Channel) *lib.Channel {
	leftOutChan := runNode(inChan, node.Left)
	rightOutChan := runNode(inChan, node.Right)

	leftData := []lib.ShellData{}
	rightData := []lib.ShellData{}

	var wg sync.WaitGroup
	wg.Add(2)

	pipeChan := func(inChan *lib.Channel, data *[]lib.ShellData) {
		for newData, ok := inChan.Read(); ok; newData, ok = inChan.Read() {
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
	differenceChan := lib.NewChannel()
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
				differenceChan.Write(left)
			}
		}
		differenceChan.Close()
	}()

	return differenceChan
}

func runNode(inChan *lib.Channel, node PipelineNode) *lib.Channel {
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

func (this LastCmd) Call(inChan, outChan *lib.Channel, arguments []string) {
	for _, last := range lastOut {
		outChan.Write(last)
	}
	outChan.Close()
}

// ### used by last command. ideally, it would somehow keep hold of this itself
var lastOut []lib.ShellData

// Wait for a command to finish, presenting data as it arrives.
func present(outChan *lib.Channel) {
	var newOut []lib.ShellData
	const presentDebug = false
	for in, ok := outChan.Read(); ok; in, ok = outChan.Read() {
		newOut = append(newOut, in)
		if presentDebug {
			fmt.Printf("%T: %s", in, in.Present())
		} else {
			fmt.Printf("%s", in.Present())
		}
	}
	lastOut = newOut
}

var exited = false

const shouldTrace = false

func main() {
	if shouldTrace {
		log.Printf("Tracing")
		f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		err = trace.Start(f)
		if err != nil {
			panic(err)
		}
		defer func() {
			log.Printf("Stopping trace")
			trace.Stop()
			f.Close()
		}()
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("smokey the shell")
	fmt.Println("try something like echo hello my friend | cat")
	fmt.Println("---------------------")

	for !exited {
		fmt.Print("% ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		runCommandString(text)
	}
}

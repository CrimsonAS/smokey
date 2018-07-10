package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// All commands implement this interface.
type commandObject interface {
	// Call the comand. The inChan and outChan are used for communication.
	// The arguments let it customize its behaviour from the command line.
	Call(inChan chan shellData, outChan chan shellData, arguments []string)
}

// Parse and execute a given command pipeline.
func runCommandString(text string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic", r)
		}
	}()
	commands := parsePipeline(text)
	var inChan chan shellData
	var outChan chan shellData

	inChan = make(chan shellData)
	outChan = make(chan shellData)
	close(inChan) // ### not what we should do really

	for idx, cmd := range commands {
		var commandObject commandObject
		switch cmd.Command {
		case "influxConnect":
			commandObject = InfluxConnect{}
		case "influxQuery":
			commandObject = InfluxQuery{}
		case "ps":
			commandObject = PsCmd{}
		case "kill":
			commandObject = KillCmd{}
		case "head":
			commandObject = HeadCmd{}
		case "tail":
			commandObject = TailCmd{}
		case "echo":
			commandObject = EchoCmd{}
		case "cat":
			commandObject = CatCmd{}
		case "dup":
			commandObject = DupCmd{}
		case "uniq":
			commandObject = UniqCmd{}
		case "lines":
			commandObject = LinesCmd{}
		case "last":
			commandObject = LastCmd{}
		case "ls":
			commandObject = LsCmd{}
		case "cd":
			commandObject = CdCmd{}
		case "fetch":
			commandObject = FetchCmd{}
		case "grep":
			commandObject = GrepCmd{}
		default:
			commandObject = StandardProcessCmd{process: cmd.Command}
		}
		go commandObject.Call(inChan, outChan, cmd.Arguments)

		inChan = outChan
		if idx < len(commands)-1 {
			outChan = make(chan shellData)
		}
	}

	present(outChan)
}

// LastCmd just repeats whatever shell data the last command pipeline produced.
// This doesn't really belong here, but right now it hacks present(), so it's here
// for easy reference.
type LastCmd struct {
}

func (this LastCmd) Call(inChan chan shellData, outChan chan shellData, arguments []string) {
	for _, last := range lastOut {
		outChan <- last
	}
	close(outChan)
}

// ### used by last command. ideally, it would somehow keep hold of this itself
var lastOut []shellData

// Wait for a command to finish, presenting data as it arrives.
func present(outChan chan shellData) {
	var newOut []shellData
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

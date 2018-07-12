package lib

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// Register the types we'll be sending.
func init() {
	var ss ShellString
	gob.Register(ss)
	var sb ShellBuffer
	gob.Register(sb)
	gob.Register(eofMessage{})
	gob.Register(&ShellPath{})
}

const gobDebug = false

// Used to mark end of messages.
type eofMessage struct{}

// Decode messages from a conn, and push them to outChan.
func ReadToChannel(conn io.Reader, outChan *Channel) {
	dec := gob.NewDecoder(conn)
	for gobFinished := false; !gobFinished; {
		var p interface{}
		err := dec.Decode(&p)
		if err != nil {
			panic(fmt.Sprintf("Plugin decode failed: %s", err))
		}
		if gobDebug {
			log.Printf("Read from GOB: %T %+v, %+v", p, p, err)
		}

		switch sd := p.(type) {
		case ShellData:
			outChan.Write(sd)
		case eofMessage:
			if gobDebug {
				log.Printf("Reading GOB got EOF")
			}
			gobFinished = true
		default:
			panic(fmt.Sprintf("unknown data %T %+v", p, p))
		}
	}
	if gobDebug {
		log.Printf("Reading GOB done")
	}
	outChan.Close()
}

// Write messages from inChan to a conn.
func WriteFromChannel(conn io.Writer, inChan *Channel) {
	enc := gob.NewEncoder(conn)

	gobWrite := func(val interface{}) {
		err := enc.Encode(&val)
		if gobDebug {
			log.Printf("Wrote %T %+v from channel", val, val)
		}
		if err != nil {
			panic(fmt.Sprintf("Failed writing in pipeline: %s", err))
		}
	}

	for in, ok := inChan.Read(); ok; in, ok = inChan.Read() {
		gobWrite(in)
	}

	var eof interface{}
	eof = eofMessage{}
	gobWrite(eof)
	if gobDebug {
		log.Printf("EOF sent, Writing done")
	}
}

// Establishes a connection to the host from a plugin.
func DialPlugin() *PluginContext {
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("Can't connect: %s", err))
	}

	if gobDebug {
		log.Printf("Connected")
	}

	inChan := NewChannel()
	outChan := NewChannel()

	// Decode from main to inChan
	// Write from outChan to main
	// We must not finish until we have finished with i/o!
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if gobDebug {
			log.Printf("PLUGIN: Reading to channel...")
		}
		ReadToChannel(conn, inChan)
		if gobDebug {
			log.Printf("PLUGIN: Reading to channel done")
		}
		wg.Done()
		if gobDebug {
			log.Printf("PLUGIN: Reading to channel WG done")
		}
	}()

	go func() {
		if gobDebug {
			log.Printf("PLUGIN: Writing to channel...")
		}
		WriteFromChannel(conn, outChan)
		if gobDebug {
			log.Printf("PLUGIN: Writing to channel done")
		}
		wg.Done()
		if gobDebug {
			log.Printf("PLUGIN: Writing to channel WG done")
		}
	}()

	return &PluginContext{
		InChan:  inChan,
		OutChan: outChan,
		conn:    conn,
		wg:      &wg,
	}
}

// Used to provide communication with the host from a plugin.
type PluginContext struct {
	// Messages to the plugin
	InChan *Channel

	// Messages from the plugin
	OutChan *Channel

	conn net.Conn
	wg   *sync.WaitGroup
}

// Wait for all sending/recieving to finish, and exit cleanly.
func (this *PluginContext) Wait() {
	if gobDebug {
		log.Printf("PLUGIN: Waiting to close. Waiting WG...")
	}
	this.wg.Wait()
	if gobDebug {
		log.Printf("PLUGIN: Closing conn...")
	}
	if gobDebug {
		log.Printf("PLUGIN: Waiting DONE!")
	}

}

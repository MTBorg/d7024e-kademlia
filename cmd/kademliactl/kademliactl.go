package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func reader(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()

	//TODO: Don't hardcode buffer size to 1024 bytes
	buf := make([]byte, 1024)
	_, err := r.Read(buf[:])
	if err != nil {
		return
	}
}

func readCmdLineArgs() []byte {
	var data strings.Builder
	for index, arg := range os.Args[1:] { //Skip first argument (program name)
		data.WriteString(arg)
		if index < len(os.Args)-2 { // Add space if not last word
			data.WriteString(" ")
		}
	}
	return []byte(data.String())
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if len(os.Args) > 1 { // If a command was specified
		c, err := net.Dial("unix", "/tmp/echo.sock")
		if err != nil {
			panic(err)
		}
		defer c.Close()

		var wg sync.WaitGroup
		wg.Add(1)

		// Make sure reader is set up before writing
		go reader(&wg, c)

		data := readCmdLineArgs()

		_, err = c.Write(data)
		if err != nil {
			log.Error().Msgf("Failed to write to socket: %s", err.Error())
		}
		wg.Wait()
	} else {
		//TODO: Print usage
		log.Print("Usage: To be done...")
	}
}

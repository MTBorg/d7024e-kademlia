package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"kademlia/internal/utils/arrays"
	"net"
	"os"
	"sync"
	"time"
)

func reader(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()

	//TODO: Don't hardcode buffer size to 1024 bytes
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		return
	}
	log.Info().Msgf("Received response: %s", string(buf[:n]))
}

func sendMessage(msg *[]byte) {
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	// Make sure reader is set up before writing
	go reader(&wg, c)

	_, err = c.Write(*msg)
	if err != nil {
		log.Error().Msgf("Failed to write to socket: %s", err.Error())
	}
	wg.Wait()
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if len(os.Args) > 1 { // If a command was specified
		msg := arrays.StrArrayToByteArray(os.Args[1:])
		sendMessage(&msg)
	} else {
		//TODO: Print usage
		log.Print("Usage: To be done...")
	}
}

package main

import (
	"io"
	"kademlia/internal/logger"
	"kademlia/internal/utils/arrays"
	"net"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
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
	logger.InitLogger(os.Getenv("LOG_LEVEL"))

	if len(os.Args) > 1 { // If a command was specified
		msg := arrays.StrArrayToByteArray(os.Args[1:])
		sendMessage(&msg)
	} else {
		//TODO: Print usage
		log.Print("Usage: To be done...")
	}
}

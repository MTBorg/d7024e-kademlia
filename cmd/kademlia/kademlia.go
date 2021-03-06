package main

import (
	"kademlia/internal/command/listener"
	"kademlia/internal/http/listener"
	"kademlia/internal/logger"
	"kademlia/internal/node"
	"kademlia/internal/udplistener"
	"net"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func getHostIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Error().Msgf("Failed to get container interface addresses: %s", err)
	}
	for _, address := range addresses {

		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if err := logger.InitLogger(logLevel); err == nil {
		log.Info().Str("Level", logLevel).Msg("Log level set")
	} else {
		log.Error().Str("Level", logLevel).Msg("Failed to parse log level, defaulting to info level...")
	}

	lport, err := strconv.Atoi(os.Getenv("LISTEN_PORT"))
	if err != nil {
		log.Error().Msgf("Failed to convert env variable LISTEN_PORT from string to int: %s", err)
	}

	host, err := os.Hostname()
	ip := getHostIP()
	if err != nil {
		log.Error().Str("Host", host).Msgf("Failed to get container host: %s", err)
	}
	log.Info().Str("Hostname", host).Str("IP", ip).Msg("Starting node...")

	node := node.Node{}

	go cmdlistener.Listen(&node)
	go httplistener.Listen(&node)
	udplistener.Listen(ip, lport, &node)
}

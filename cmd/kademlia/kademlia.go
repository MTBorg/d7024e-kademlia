package main

import (
	"kademlia/internal/command/listener"
	"kademlia/internal/logger"
	"kademlia/internal/udplistener"
	"net"
	"os"

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
	logger.InitLogger(os.Getenv("LOG_LEVEL"))
	host, err := os.Hostname()
	ip := getHostIP()
	if err != nil {
		log.Error().Str("Host", host).Msgf("Failed to get container host: %s", err)
	}
	log.Info().Str("Hostname", host).Str("IP", ip).Msg("Starting node...")

	go cmdlistener.Listen()
	udplistener.Listen(ip, 1776)
}

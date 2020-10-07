package address

import (
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"strconv"
)

type Address struct {
	host string
	port string
}

func New(address string) *Address {
	lport := os.Getenv("LISTEN_PORT")
	if lport == "" { // if the env var was not defined
		lport = "1776"
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil && port == "" {
		parsedHost := net.ParseIP(address)
		if parsedHost == nil {
			log.Error().Msgf("Given address is not valid: %s", err)
			return &Address{}
		}
		host = parsedHost.String()
	}

	return &Address{
		host: host,
		port: lport,
	}
}

func (address *Address) String() string {
	return net.JoinHostPort(address.host, address.port)
}

func (address *Address) GetHost() string {
	return address.host
}

func (address *Address) GetPortAsInt() (int, error) {
	return strconv.Atoi(address.port)
}

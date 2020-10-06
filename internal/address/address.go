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
	host, port, err := net.SplitHostPort(address)
	if err != nil && port == "" {
		host = net.ParseIP(address).String()
	}

	if host == "" {
		log.Error().Msgf("Given address is not valid: %s", err)
		return &Address{}

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

package address

import (
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
)

type Address struct {
	host string
	port string
}

func New(address string) *Address {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		if port == "" { // Assume address is correct if port is missing
			host = address
		} else {
			log.Error().Str("Address", address).Msgf("Failed to parse given address, error: %s", err)
		}
	}
	return &Address{
		host: host,
		port: "1776", // TODO: Don't hardcore, maybe use env var?
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

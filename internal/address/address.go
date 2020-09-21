package address

import ()
import (
	// "github.com/rs/zerolog/log"
	"net"
)

type Address struct {
	host string
	port string
}

func New(address string) Address {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address // TODO: This is dumb, need to check if error is of "missing address port"
		// log.Error().Str("Address", host).Msgf("Failed to parse given address, error: %s", err)

	}

	return Address{
		host: host,
		port: "1776", // TODO: Don't hardcore, maybe use env var?
	}
}

func (address *Address) String() string {
	return net.JoinHostPort(address.host, address.port)
}

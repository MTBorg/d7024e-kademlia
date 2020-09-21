package udpsender

import (
	"github.com/rs/zerolog/log"
	"kademlia/internal/address"
	"net"
)

type UDPSender struct {
	target *net.UDPAddr
}

func New(target *address.Address) UDPSender {

	port, err := target.GetPortAsInt()
	if err != nil {
		log.Error().Str("Address", target.String()).Msgf("Failed to parse given string port to int: %s", err)
	}
	return UDPSender{target: &net.UDPAddr{
		IP:   net.ParseIP(target.GetHost()),
		Port: port,
	}}
}

func (udp UDPSender) Send(data string) error {
	conn, err := net.DialUDP("udp4", nil, udp.target)
	if err != nil {
		log.Error().Msgf("Failed to dial to UDP address: %s", err)
	}

	_, err = conn.Write([]byte(data))
	return err
}

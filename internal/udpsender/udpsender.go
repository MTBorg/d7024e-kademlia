package udpsender

import (
	"net"

	"github.com/rs/zerolog/log"
)

type UDPSender struct {
	target string
}

func New(target string) UDPSender {
	return UDPSender{target: target}
}

func (udp UDPSender) Send(data string) error {
	dest, err := net.ResolveUDPAddr("udp4", udp.target)
	if err != nil {
		log.Error().Msgf("Failed to resolve UDP address: %s", err)
	}

	conn, err := net.DialUDP("udp4", nil, dest)
	if err != nil {
		log.Error().Msgf("Failed to dial to UDP address: %s", err)
	}

	_, err = conn.Write([]byte(data))
	return err
}

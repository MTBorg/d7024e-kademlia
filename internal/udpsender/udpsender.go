package udpsender

import (
	"net"
	"strconv"

	"github.com/rs/zerolog/log"
)

type UDPSender struct {
	target *net.UDPAddr
}

func New(target string) UDPSender {
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		log.Error().Str("Target", target).Msgf("Failed to parse given target address: %s", err)
	}
	iport, err := strconv.Atoi(port)
	if err != nil {
		log.Error().Str("Port", port).Msgf("Failed to parse given string port to int: %s", err)
	}
	return UDPSender{target: &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: iport,
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

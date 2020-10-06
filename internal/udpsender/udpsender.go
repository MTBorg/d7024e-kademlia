package udpsender

import (
	"github.com/rs/zerolog/log"
	"kademlia/internal/address"
	"net"
	"os"
)

// UDPSender holds the target UDP address and the sender address
type UDPSender struct {
	target *net.UDPAddr
	sender *net.UDPAddr
}

// New Creates a UDPSender from given target address and resolves the sender address
func New(target *address.Address) UDPSender {

	port, err := target.GetPortAsInt()
	if err != nil {
		log.Error().Str("Address", target.String()).Msgf("Failed to parse given string port to int: %s", err)
	}

	sport := ":" + os.Getenv("SEND_PORT")           // Get the sender port from env
	laddr, err := net.ResolveUDPAddr("udp4", sport) // Set the sender address, since only port is given it only resolves the ip

	return UDPSender{target: &net.UDPAddr{IP: net.ParseIP(target.GetHost()), Port: port}, sender: laddr}
}

// Send sends the udp packet to the udp target from given port
func (udp UDPSender) Send(data string) error {

	conn, err := net.DialUDP("udp4", udp.sender, udp.target)
	defer conn.Close() // Close the connection to avoid error bind: address already in use

	if err != nil {
		log.Error().Msgf("Failed to dial to UDP address: %s", err)
	}

	_, err = conn.Write([]byte(data))
	return err
}

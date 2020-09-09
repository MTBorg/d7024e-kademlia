package network

import (
	"fmt"
	. "kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"net"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// TODO: Pref, env
var sendingPort = 1775

// Net holds the cluster network
var Net = new(Network)

// A Network consists of local address, remote address and connection
type Network struct {
	laddr *net.UDPAddr
	raddr *net.UDPAddr
	conn  *net.UDPConn
}

func (network *Network) replyPingMessage(id string) (string, error) {
	network.laddr.Port = sendingPort
	conn, err := net.DialUDP("udp", network.laddr, network.raddr)
	if err != nil {
		log.Error().Msgf("Failed to dial to UDP Address: %s", err)
		return "", err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("PONG %s", id)))
	if err != nil {
		log.Error().Msgf("Failed to write PONG to UDP Address: %s", err)
		return "", err
	}
	log.Info().Str("Address", network.raddr.String()).Msg("PONG replied to address")
	return fmt.Sprintf("PONG replied! to Address: %s", network.raddr.String()), nil
}

func (network *Network) parsePacket(data string) {
	fields := strings.Fields(data)
	if len(fields) < 1 {
		log.Error().Msgf("Packet is empty!")
	}

	switch packet := fields[0]; packet {
	case "PING":
		// TODO: Bucket AddContact (update bucket)
		network.replyPingMessage(fields[1])
	case "PONG":
		// TODO: Bucket AddContact (update bucket)
		log.Info().Str("Id", fields[1]).Msg("PONG received with id")
	default:
		log.Error().Str("packet", packet).Msg("Received unknown packet")
	}

}

// Listen initiates UDP Packet listenening on given port (UDP server)
func Listen(port int) {
	sp := strconv.Itoa(port)
	laddr, err := net.ResolveUDPAddr("udp", ":"+sp)
	if err != nil {
		log.Error().Msgf("Failed to resolve Address: %s", err)
	}
	Net.laddr = laddr
	conn, err := net.ListenUDP("udp", Net.laddr)
	if err != nil {
		log.Error().Msgf("Failed to listen on Address: %s", err)
	}
	Net.conn = conn
	log.Info().Str("Address", Net.laddr.String()).Msg("Listening on UDP packets on address")
	defer conn.Close()

	for {

		buf := make([]byte, 512)

		nr, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		Net.raddr = remoteAddr

		data := string(buf[0:nr])
		log.Info().Str("Content", data).Str("From", remoteAddr.String()).Msg("Received message from and with content,")

		Net.parsePacket(data)

	}

}

// SendPingMessage handles the client sending a PING message to a remote address
func (network *Network) SendPingMessage(contact *Contact) (string, error) {
	id := fmt.Sprint(kademliaid.NewRandomKademliaID())
	log.Info().Str("Id", id).Msg("Random Kademlia id generated")
	raddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Error().Msgf("Failed to resolve remote UDP Address: %s", err)
		return "", err
	}
	network.laddr.Port = sendingPort
	conn, err := net.DialUDP("udp", network.laddr, raddr)
	if err != nil {
		log.Error().Msgf("Failed to dial to UDP Address: %s", err)
		return "", err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("PING %s", id)))
	if err != nil {
		log.Error().Msgf("Failed to write PING to UDP Address: %s", err)
		return "", err
	}
	log.Info().Str("Address", contact.Address).Msg("PING sent to address")
	return fmt.Sprintf("PING SENT! to Address: %s", contact.Address), nil
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

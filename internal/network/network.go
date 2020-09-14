package network

import (
	"github.com/rs/zerolog/log"
	"kademlia/internal/contact"
	"kademlia/internal/kademliaid"
	"kademlia/internal/rpc"
	"kademlia/internal/udpsender"
	"net"
	"strconv"
	"strings"
)

var Net Network

// A Network consists of local address, remote address and connection
type Network struct {
	listenPort string
}

func (network *Network) parsePacket(sender *contact.Contact, rpcID *kademliaid.KademliaID, data string) {
	log.Debug().Str("String", data).Msg("Parsing data string")
	fields := strings.Fields(data)
	if len(fields) < 1 {
		log.Error().Msgf("Packet is empty!")
	}

	switch packet := fields[0]; packet {
	case "PING":
		// TODO: Bucket AddContact (update bucket)
		log.Info().Str("Id", rpcID.String()).Msg("PING received with RPC id")
		network.SendPongMessage(sender.Address, rpcID)

	case "PONG":
		// TODO: Bucket AddContact (update bucket)
		log.Info().Str("Id", rpcID.String()).Msg("PONG received with RPC id")
	default:
		log.Error().Str("packet", packet).Msg("Received packet with unkown command")

	}

}

func (network *Network) received(c *net.UDPConn) {
	for {
		buf := make([]byte, 512)
		nr, addr, err := c.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		data := buf[0:nr]
		rpc, err := rpc.Deserialize(string(data))
		if err == nil {
			log.Info().Str("Content", rpc.Content).Str("SenderId", rpc.SenderId.String()).Msg("Received message")

			c := contact.NewContact(rpc.SenderId, addr.String())
			network.parsePacket(&c, rpc.RPCId, rpc.Content)

			log.Debug().Str("Id", c.ID.String()).Str("Address", c.Address).Msg("Updating bucket")
			// TODO: Add to routing table
		} else {
			log.Warn().Str("Error", err.Error()).Msg("Failed to deserialize message")
		}
	}
}

// Listen initiates a UDP server
func Listen(ip string, port int) {
	Net.listenPort = strconv.Itoa(port)

	addr := net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	ln, err := net.ListenUDP("udp4", &addr)
	if err != nil {
		log.Error().Msgf("Failed to listen on UDP Address: %s", err)
	}
	log.Info().Str("Address", addr.String()).Msg("Listening on UDP packets on address")
	defer ln.Close()

	Net.received(ln)
}

// SendPongMessage replies a "PONG" message to the remote "pinger" address
func (network *Network) SendPongMessage(target string, id *kademliaid.KademliaID) {
	host, _, err := net.SplitHostPort(target)
	if err != nil {
		log.Error().Str("Target", host).Msgf("Failed to parse given target address: %s", err)
	}
	target = net.JoinHostPort(host, Net.listenPort)
	log.Debug().Str("Address", target).Msg("Sending PONG to address")
	rpc := rpc.New("PONG", target)
	rpc.RPCId = id
	udpSender := udpsender.New(target)
	err = rpc.Send(udpSender)

	if err != nil {
		log.Error().Msgf("Failed to write RPC PING message to UDP: %s", err.Error())
		log.Info().Str("Address", target).Str("Content", "PING").Msg("Message sent to address")
	}
}

// SendPingMessage sends a "PING" message to a remote address
func (network *Network) SendPingMessage(target string) {

	rpc := rpc.New("PING", target)
	udpSender := udpsender.New(target)
	err := rpc.Send(udpSender)

	if err != nil {
		log.Error().Msgf("Failed to write RPC PING message to UDP: %s", err.Error())
		log.Info().Str("Address", target).Str("Content", "PING").Msg("Message sent to address")
	}
}

func (network *Network) SendFindContactMessage(contact *contact.Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

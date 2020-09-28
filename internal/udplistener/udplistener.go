package udplistener

import (
	"kademlia/internal/address"
	"kademlia/internal/contact"
	"kademlia/internal/node"
	"kademlia/internal/rpc"
	"kademlia/internal/rpc/parser"
	"net"
	"strings"

	"github.com/rs/zerolog/log"
)

// Listen initiates a UDP server
func Listen(ip string, port int, node *node.Node) {
	addr := net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	ln, err := net.ListenUDP("udp4", &addr)
	defer ln.Close()
	if err != nil {
		log.Error().Msgf("Failed to listen on UDP Address: %s", err)
	}
	log.Info().Str("Address", addr.String()).Msg("Listening on UDP packets on address")

	waitForMessages(ln, node)
}

func waitForMessages(c *net.UDPConn, node *node.Node) {
	for {
		buf := make([]byte, 512)
		nr, addr, err := c.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		data := buf[0:nr]
		adr := address.New(addr.String())
		rpcMsg, err := rpc.Deserialize(string(data))
		if err == nil {
			c := contact.NewContact(rpcMsg.SenderId, adr)
			node.RoutingTable.AddContact(c)

			cmd, err := rpcparser.ParseRPC(&c, &rpcMsg)
			if err != nil {
				log.Warn().Str("Error", err.Error()).Msg("Failed to parse RPC")
				continue
			}

			options := strings.Split(rpcMsg.Content, " ")[1:]
			if err = cmd.ParseOptions(&options); err == nil {
				cmd.Execute(node)
			} else {
				log.Warn().
					Str("Error", err.Error()).
					Msg("Failed to parse RPC options")
			}

			log.Trace().Str("NodeID", c.ID.String()).Str("Address", c.Address.String()).Msg("Inserting new node to bucket")
		} else {
			log.Warn().Str("Error", err.Error()).Msg("Failed to deserialize message in UDPListener")
		}
	}
}

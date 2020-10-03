package udpsender

import (
	"kademlia/internal/address"
	"net"
	"os"

	"github.com/rs/zerolog/log"
)

// UDPSender holds the target UDP address and the sender address
type UDPSender struct {
	conn *net.UDPConn
}

func New() (*UDPSender, error) {
	sport := ":" + os.Getenv("SEND_PORT") // Get the sender port from env
	laddr, err := net.ResolveUDPAddr("udp4", sport)
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}
	return &UDPSender{conn: conn}, nil
}

func (udp *UDPSender) Send(data string, target *address.Address) error {
	adr, err := net.ResolveUDPAddr("udp", target.String())
	if err != nil {
		return err
	}

	var n int
	n, err = udp.conn.WriteTo([]byte(data), adr)
	log.Trace().Int("BytesSent", n).
		Str("LAddr", udp.conn.LocalAddr().String()).
		Str("RAddr", target.String()).
		Msg("Send UDP packet")
	return err
}

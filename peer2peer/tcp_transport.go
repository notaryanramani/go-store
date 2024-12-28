package peer2peer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type TCPTransport struct {
	Ops     TCPTransportOps
	Listner net.Listener
	msgch   chan Message
}

func NewTCPTransport(ops TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		Ops:   ops,
		msgch: make(chan Message),
	}
}

// TCPPeer represents remote node
// over TCP connection
type TCPPeer struct {
	// Conn is connection of peers
	Conn net.Conn

	// Outbound indicates if the peer is outbound or inbound
	// If true, the peer is outbound (dial)
	// If false, the peer is inbound (accept)
	Outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		Outbound: outbound,
	}
}

type TCPTransportOps struct {
	ListnAddr  string
	Handshaker Handshaker
	Decoder    Decoder
	OnPeer     func(Peer) error
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	listner, err := net.Listen("tcp", t.Ops.ListnAddr)
	if err != nil {
		return err
	}

	t.Listner = listner
	go t.accept()
	log.Printf("Listening on %s\n", t.Ops.ListnAddr)

	return nil
}

func (t *TCPTransport) accept() {
	for {
		conn, err := t.Listner.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
		}
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	var err error

	defer func() {
		fmt.Println("Dropping peer connection: ", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err = t.Ops.Handshaker(peer); err != nil {
		return
	}

	if t.Ops.OnPeer != nil {
		if err = t.Ops.OnPeer(peer); err != nil {
			return
		}
	}

	msg := &Message{}
	for {
		if err = t.Ops.Decoder.Decode(conn, msg); err != nil {
			switch {
			case err == io.EOF:
				return
			case errors.Is(err, net.ErrClosed):
				return
			default:
				fmt.Println("TCP Read Error: \n", err)
			}
		}
		msg.From = conn.RemoteAddr()
		t.msgch <- *msg
	}
}

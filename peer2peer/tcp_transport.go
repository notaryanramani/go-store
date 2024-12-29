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

// Send implements the Peer interface and
// sends the message to the peer
func (p *TCPPeer) Send(msg []byte) error {
	_, err := p.Conn.Write(msg)
	return err
}

// RemoteAddr implements the Peer interface and
// returns the remote address of the peer
// of its underlying connection
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.Conn.RemoteAddr()
}

// Close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.Conn.Close()
}

type TCPTransportOps struct {
	ListnAddr  string
	Handshaker Handshaker
	Decoder    Decoder
	OnPeer     func(Peer) error
}

// ListenAndAccept implements the Transport interface and listens
// for incoming connections from other peers in the network
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

// Consume implements the Transport interface and
// only returns read only channel
// for reading messages received from another peer
// in the network
func (t *TCPTransport) Consume() <-chan Message {
	return t.msgch
}

// Dial implements the Transport interface and dials
// another peer in the network
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Printf("Dialed %s\n", addr)
	go t.handleConnection(conn, true)

	return nil
}

// Close implements the Transport interface and closes the listener
func (t *TCPTransport) Close() error {
	return t.Listner.Close()
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
		go t.handleConnection(conn, false)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Println("Dropping peer connection: ", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)
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

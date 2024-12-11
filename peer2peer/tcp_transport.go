package peer2peer

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransport struct {
	Ops 	TCPTransportOps
	Listner net.Listener

	Peers 	map[net.Addr]Peer
	Mu 		sync.RWMutex
}

func NewTCPTransport(ops TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		Ops: ops,
		Peers: make(map[net.Addr]Peer),
	}
}

// TCPPeer represents remote node 
// over TCP connection
type TCPPeer struct {
	// Conn is connection of peers
	Conn 		net.Conn

	// Outbound indicates if the peer is outbound or inbound
	// If true, the peer is outbound (dial)
	// If false, the peer is inbound (accept)
	Outbound 	bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn: conn,
		Outbound: outbound,
	}
}

type TCPTransportOps struct {
	ListnAddr 	string
	Handshaker 	Handshaker
	Decoder 	Decoder
}
 
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	
	listner, err := net.Listen("tcp", t.Ops.ListnAddr)
	if err != nil {
		return err
	}

	t.Listner = listner
	go t.accept()
	return nil
}

func (t *TCPTransport) accept () {
	for {
		conn, err := t.Listner.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
		}
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn)  {
	peer := NewTCPPeer(conn, true)
	if err := t.Ops.Handshaker(peer); err != nil {
		conn.Close()
		fmt.Println("Error handshaking with peer: \n", err)
		return
	}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error decoding message: \n", err)
			return
		}
		fmt.Println("Received message: \n", buffer[:n - 2])
	}
}
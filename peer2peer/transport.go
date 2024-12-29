package peer2peer

import "net"

// Peer is representation of nodes in the network
type Peer interface {
	net.Conn
	Close() error
}

// Transport is an interface that handles
// communication between peers
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan Message
	Close() error
	Dial(addr string) error
	ListenAddr() string
}

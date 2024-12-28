package peer2peer

// Peer is representation of nodes in the network
type Peer interface {
	Close() error
}

// Close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.Conn.Close()
}

// Transport is an interface that handles
// communication between peers
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan Message
	Close() error
}

// Consume only returns read only channel
// for reading messages received from another peer
// in the network
func (t *TCPTransport) Consume() <-chan Message {
	return t.msgch
}

// Close closes the listener
func (t *TCPTransport) Close() error {
	return t.Listner.Close()
}

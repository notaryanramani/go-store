package peer2peer


// Peer is representation of nodes in the network
type Peer interface {}

// Transport is an interface that handles
// communication between peers
type Transport interface {
	ListenAndAccept() error
}
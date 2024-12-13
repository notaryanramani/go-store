package peer2peer

import "net"

// Message holds any data which is being
// sent over the transport between peers
// in the network
type Message struct {

	// From is the address of the sender (peer)
	From net.Addr

	// Payload is the data stored in byte buffer
	Payload []byte
}

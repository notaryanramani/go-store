package peer2peer

// Message holds any data which is being
// sent over the transport between peers 
// in the network
type Message struct {
	Payload []byte
}
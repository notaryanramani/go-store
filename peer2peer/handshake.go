package peer2peer

type Handshaker func(Peer) error

func NoHandshake(Peer) error {
	return nil
}
package peer2peer

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTCPTransport (t *testing.T) {
	trOps := TCPTransportOps{
		ListnAddr: ":8080",
		Handshaker: NoHandshake,
		Decoder: GOBDecoder{},
	}

	transport := NewTCPTransport(trOps)
	assert.Equal(t, ":8080", transport.Ops.ListnAddr)
	assert.Nil(t, transport.ListenAndAccept())
}
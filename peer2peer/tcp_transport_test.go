package peer2peer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	trOps := TCPTransportOps{
		ListnAddr:  ":3000",
		Handshaker: NoHandshake,
		Decoder:    DefaultDecoder{},
	}

	transport := NewTCPTransport(trOps)
	assert.Equal(t, ":3000", transport.Ops.ListnAddr)
	assert.Nil(t, transport.ListenAndAccept())
}

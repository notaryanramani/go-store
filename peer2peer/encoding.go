package peer2peer

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	// Decode reads the data from the reader
	// and stores it in a byte buffer
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *Message) error {
	return gob.NewDecoder(r).Decode(msg)
}

type Encoder interface {
	Encode(io.Writer) error
}

type GOBEncoder struct{}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *Message) error {
	buf := make([]byte, 1024*5)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil
}

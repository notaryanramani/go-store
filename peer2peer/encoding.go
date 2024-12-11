package peer2peer

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode (io.Reader, any) (error)
}

type GOBDecoder struct {}

func (dec GOBDecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

type Encoder interface {
	Encode (io.Writer) (error)
}

type GOBEncoder struct {}

type NoDecoder struct {}

func (dec NoDecoder) Decode(r io.Reader, v any) error {
	buf := make([]byte, 1024)
	_, err := r.Read(buf)
	if err != nil {
		return err
	}
	return nil
}


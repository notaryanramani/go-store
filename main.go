package main

import (
	"log"
	"time"

	"github.com/notaryanramani/go-store/peer2peer"
)

func onPeer(p peer2peer.Peer) error {
	p.Close()
	return nil
}

func main() {
	log.Println("Starting file server...")
	trOps := peer2peer.TCPTransportOps{
		ListnAddr:  ":3000",
		Handshaker: peer2peer.NoHandshake,
		Decoder:    peer2peer.DefaultDecoder{},
		OnPeer:     onPeer,
	}
	transport := peer2peer.NewTCPTransport(trOps)

	fileServerOpts := NewServerOpts(":3000", "3000_tmp", CASPathTransformFunc, transport)
	fileServer := NewServer(fileServerOpts)

	go func() {
		time.Sleep(2 * time.Second)
		fileServer.Stop()
	}()

	if err := fileServer.Start(); err != nil {
		log.Fatal(err)
	}
}

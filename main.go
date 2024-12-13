package main

import (
	"fmt"
	"log"

	"github.com/notaryanramani/go-store/peer2peer"
)

func onPeer(p peer2peer.Peer) error {
	p.Close()
	return nil
}

func main() {
	fmt.Println("executing main.go...")
	tcp_opts := peer2peer.TCPTransportOps{
		ListnAddr:  ":8080",
		Handshaker: peer2peer.NoHandshake,
		Decoder:    peer2peer.DefaultDecoder{},
		OnPeer:     onPeer,
	}
	transport := peer2peer.NewTCPTransport(tcp_opts)

	go func() {
		for {
			msg := <-transport.Consume()
			fmt.Printf("Received message: %s \n", msg.Payload)
		}
	}()

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatalf("Error listening and accepting connections: %v", err)
	}
	select {}
}

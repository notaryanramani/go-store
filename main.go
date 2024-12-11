package main

import (
	"fmt"
	"github.com/notaryanramani/go-store/peer2peer"
	"log"
)


func main() {
	fmt.Println("Hello, World!")
	tcp_opts := peer2peer.TCPTransportOps{
		ListnAddr: ":8080",
		Handshaker: peer2peer.NoHandshake,
		Decoder: peer2peer.GOBDecoder{},
	}
	transport := peer2peer.NewTCPTransport(tcp_opts)
	if err := transport.ListenAndAccept(); err != nil {
		log.Fatalf("Error listening and accepting connections: %v", err)
	}
	select {}
}

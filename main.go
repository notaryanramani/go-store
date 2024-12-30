package main

import (
	"bytes"
	"log"
	"time"

	"github.com/notaryanramani/go-store/peer2peer"
)

func makeFileServer(l string, r string, n ...string) *FileServer {
	trOps := peer2peer.TCPTransportOps{
		ListnAddr:  l,
		Handshaker: peer2peer.NoHandshake,
		Decoder:    peer2peer.DefaultDecoder{},
	}
	transport := peer2peer.NewTCPTransport(trOps)

	fileServerOpts := NewServerOpts(l, r, CASPathTransformFunc, transport)
	fileServer := NewServer(fileServerOpts)
	fileServer.Opts.Nodes = n

	transport.Ops.OnPeer = fileServer.OnPeer

	return fileServer
}

func main() {
	log.Println("Starting file server...")

	fs1 := makeFileServer(":3000", "fs1", "")
	fs2 := makeFileServer(":3001", "fs2", ":3000")

	go func() {
		log.Fatal(fs1.Start())
	}()

	go func() {
		fs2.Start()
	}()

	time.Sleep(1 * time.Second)

	data := bytes.NewReader([]byte("Other Data!"))
	fs2.StoreWrite("myprivatekey", data)

	select {}
}

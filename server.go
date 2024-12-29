package main

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/notaryanramani/go-store/peer2peer"
)

type FileServerOpts struct {
	ListenAddr          string
	StoreRoot           string
	PathTransformedFunc getPathTransformed
	Transport           peer2peer.Transport
	Nodes               []string
}

func NewServerOpts(l string, r string, p getPathTransformed, t peer2peer.Transport) *FileServerOpts {
	return &FileServerOpts{
		ListenAddr:          l,
		StoreRoot:           r,
		PathTransformedFunc: p,
		Transport:           t,
	}
}

type FileServer struct {
	Opts  *FileServerOpts
	Store *Store

	peerLock sync.Mutex
	peers    map[string]peer2peer.Peer

	quitch chan struct{}
}

func NewServer(Opts *FileServerOpts) *FileServer {
	storeOpts := NewStoreOpts(Opts.PathTransformedFunc, Opts.StoreRoot)

	return &FileServer{
		Opts:   Opts,
		Store:  NewStore(storeOpts),
		quitch: make(chan struct{}),
		peers:  make(map[string]peer2peer.Peer),
	}
}

// Stop stops the file server by closing the quitch channel
func (fs *FileServer) Stop() {
	close(fs.quitch)
}

// OnPeer
func (fs *FileServer) OnPeer(p peer2peer.Peer) error {
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()

	remoteAddrString := p.RemoteAddr().String()

	log.Printf("New peer connected: %s\n", remoteAddrString)
	fs.peers[remoteAddrString] = p

	return nil
}

// Loop listens for incoming messages
func (fs *FileServer) Loop() {
	defer func() {
		log.Println("Shutting down file server...")
		fs.Opts.Transport.Close()
	}()

	for {
		select {
		case msg := <-fs.Opts.Transport.Consume():
			fmt.Println("Received message: ", msg)
		case <-fs.quitch:
			return
		}
	}
}

// Bootstrap Network with nodes
func (fs *FileServer) Bootstrap() error {
	for _, nodeAddr := range fs.Opts.Nodes {

		if len(nodeAddr) == 0 {
			continue
		}

		go func(addr string) {
			if err := fs.Opts.Transport.Dial(addr); err != nil {
				log.Printf("Error dialing %s: %v\n", addr, err)
			}
		}(nodeAddr)
	}
	return nil
}

// Start starts the file server
func (fs *FileServer) Start() error {
	if err := fs.Opts.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(fs.Opts.Nodes) > 0 {
		fs.Bootstrap()
	}
	fs.Loop()

	return nil
}

// StoreRead reads the file from the store
func (fs *FileServer) StoreWrite(key string, r io.Reader) error {
	return fs.Store.Write(key, r)
}

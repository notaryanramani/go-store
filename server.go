package main

import (
	"fmt"
	"io"
	"log"

	"github.com/notaryanramani/go-store/peer2peer"
)

type FileServerOpts struct {
	ListenAddr          string
	StoreRoot           string
	PathTransformedFunc getPathTransformed
	Transport           peer2peer.Transport
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
	Opts   *FileServerOpts
	Store  *Store
	quitch chan struct{}
}

func NewServer(Opts *FileServerOpts) *FileServer {
	storeOpts := NewStoreOpts(Opts.PathTransformedFunc, Opts.StoreRoot)

	return &FileServer{
		Opts:   Opts,
		Store:  NewStore(storeOpts),
		quitch: make(chan struct{}),
	}
}

// Stop stops the file server by closing the quitch channel
func (fs *FileServer) Stop() {
	close(fs.quitch)
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

// Start starts the file server
func (fs *FileServer) Start() error {
	if err := fs.Opts.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.Loop()

	return nil
}

// StoreRead reads the file from the store
func (fs *FileServer) StoreWrite(key string, r io.Reader) error {
	return fs.Store.Write(key, r)
}

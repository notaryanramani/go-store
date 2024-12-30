package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

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
			var pMsg MessagePayload
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&pMsg); err != nil {
				log.Printf("Boo: %v\n", err)
			}
			if err := fs.handleMessagePayload(msg.From.String(), &pMsg); err != nil {
				log.Printf("Error handling message: %v\n", err)
			}

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

// StoreWrite writes the content of the reader to the store
func (fs *FileServer) StoreWrite(key string, r io.Reader) error {
	fBuf := new(bytes.Buffer)
	tee := io.TeeReader(r, fBuf)

	n, err := fs.Store.Write(key, tee)
	if err != nil {
		return err
	}

	msg := &MessagePayload{
		From: fs.Opts.Transport.ListenAddr(),
		Payload: &MessageStoreFile{
			Key:  key,
			Size: n,
		},
	}

	msgBuf := new(bytes.Buffer)
	if err := gob.NewEncoder(msgBuf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range fs.peers {
		err := peer.Send(msgBuf.Bytes())
		if err != nil {
			return err
		}
	}

	time.Sleep(2 * time.Second)

	for _, peer := range fs.peers {
		_, err := io.Copy(peer, fBuf)
		if err != nil {
			return err
		}
	}

	return nil
}

type MessagePayload struct {
	From    string
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func (fs *FileServer) broadcast(mp *MessagePayload) error {
	peers := []io.Writer{}

	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(mp)
}

func init() {
	gob.Register(&MessagePayload{})
	gob.Register(&MessageStoreFile{})
}

func (fs *FileServer) handleMessagePayload(from string, m *MessagePayload) error {
	switch v := m.Payload.(type) {
	case *MessageStoreFile:
		fmt.Printf("Received payload: %+v\n", v)
		return fs.handleMesssageStoreFile(from, v)
	}
	return nil
}

func (fs *FileServer) handleMesssageStoreFile(from string, m *MessageStoreFile) error {
	peer, ok := fs.peers[from]
	if !ok {
		return fmt.Errorf("peer [%s] not found", from)
	}

	n, err := fs.Store.Write(m.Key, io.LimitReader(peer, m.Size)) 
	if err != nil {
		return err
	}

	log.Printf("Received %d bytes from %s\n", n, from)

	peer.(*peer2peer.TCPPeer).Wg.Done()

	return nil
}

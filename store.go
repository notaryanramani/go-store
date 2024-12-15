package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLength := len(hashStr) / blockSize // folder depth

	paths := make([]string, sliceLength)
	for i := 0; i < sliceLength; i++ {
		paths[i] = hashStr[i*blockSize : (i+1)*blockSize]
	}

	PathName := strings.Join(paths, "/")
	return PathKey{
		PathName: PathName,
		FileName: hashStr,
	}
}

type getPathTransformed func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPathName() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

var defaultGetPathTransformed = func(path string) string {
	return path
}

type StoreOpts struct {
	getPathTransformed getPathTransformed
}

func NewStoreOpts(f getPathTransformed) StoreOpts {
	return StoreOpts{
		getPathTransformed: f,
	}
}

type Store struct {
	StOps StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StOps: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.StOps.getPathTransformed(key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := pathKey.FullPathName()
	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Wrote %d bytes to %s", n, pathAndFileName)

	return nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.StOps.getPathTransformed(key)
	pathAndFileName := pathKey.FullPathName()
	return os.Open(pathAndFileName)
}

func (s *Store) Read(key string) (io.Reader, error) {
	r, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, r)

	return buf, err
}

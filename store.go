package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const DefaultRootFolderName = "root"

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

func DefaultPathKey(key string) PathKey {
	return PathKey{
		PathName: "default",
		FileName: key,
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

func (p PathKey) getFirstFolder() string {
	return strings.Split(p.PathName, "/")[0]
}

type StoreOpts struct {
	// Root contains folder name of the root folder,
	// containing all the files of the system.
	Root               string
	getPathTransformed getPathTransformed
}

func NewStoreOpts(f getPathTransformed, r string) StoreOpts {
	return StoreOpts{
		Root:               r,
		getPathTransformed: f,
	}
}

type Store struct {
	StOps StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.getPathTransformed == nil {
		opts.getPathTransformed = DefaultPathKey
	}
	if len(opts.Root) == 0 {
		opts.Root = DefaultRootFolderName
	}
	return &Store{
		StOps: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.StOps.getPathTransformed(key)

	// Creating the folders
	folderPathWithRoot := s.StOps.Root + "/" + pathKey.PathName
	if err := os.MkdirAll(folderPathWithRoot, os.ModePerm); err != nil {
		return 0, err
	}

	// Creating the file
	fullPathWithRoot := s.StOps.Root + "/" + pathKey.FullPathName()
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return 0,err
	}

	// Writing the file
	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.StOps.getPathTransformed(key)
	fullPathWithRoot := s.StOps.Root + "/" + pathKey.FullPathName()
	return os.Open(fullPathWithRoot)
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

func (s *Store) Delete(key string) error {
	pathKey := s.StOps.getPathTransformed(key)

	defer func() {
		log.Printf("Removed [%s] from disk", pathKey.FileName)
	}()

	firstFolder := pathKey.getFirstFolder()
	firstFolderWithRoot := s.StOps.Root + "/" + firstFolder
	return os.RemoveAll(firstFolderWithRoot)
}

func (s *Store) Has(key string) bool {
	pathKey := s.StOps.getPathTransformed(key)
	fullPathWithRoot := s.StOps.Root + "/" + pathKey.FullPathName()

	_, err := os.Stat(fullPathWithRoot)
	return err == fs.ErrNotExist
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.StOps.Root)
}

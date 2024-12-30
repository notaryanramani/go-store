package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	store := getNewTestStore()
	defer clear(t, store)

	assert.NotNil(t, store)

	for i := 0; i < 10; i++ {
		key := generateRandomKey()
		data := generateRandomData()
		writeData := bytes.NewReader(data)

		_, err := store.writeStream(key, writeData)
		if err != nil {
			t.Fatalf("Error writing stream: %v", err)
		}

		if ok := store.Has(key); ok {
			t.Fatalf("Expected to have %s folder", key)
		}

		r, err := store.Read(key)
		if err != nil {
			t.Fatalf("Error reading stream: %v", err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Fatalf("Expected %s, got %s", data, b)
		}
	}
}

func TestPathTransformFunc(t *testing.T) {
	path := "SpecialPictures"
	pathKey := CASPathTransformFunc(path)
	expectedPathName := "0471c/0f32e/e0383/21c2b/50fbd/5e832/ede43/6c7ca"
	expectedFileName := "0471c0f32ee038321c2b50fbd5e832ede436c7ca"
	if pathKey.PathName != expectedPathName {
		t.Fatalf("Expected %s, got %s", expectedPathName, pathKey.PathName)
	}
	if pathKey.FileName != expectedFileName {
		t.Fatalf("Expected %s, got %s", expectedFileName, pathKey.FileName)
	}
}

// helpers
func getNewTestStore() *Store {
	stOpts := NewStoreOpts(
		CASPathTransformFunc,
		"root",
	)
	return NewStore(stOpts)
}

func clear(t *testing.T, store *Store) {
	if err := store.Clear(); err != nil {
		t.Fatalf("Error clearing store: %v", err)
	}
}

func generateRandomKey() string {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)
}

func generateRandomData() []byte {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return []byte(base32.StdEncoding.EncodeToString(randomBytes))
}

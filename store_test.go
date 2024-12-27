package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	stOpts := NewStoreOpts(
		CASPathTransformFunc,
	)
	store := NewStore(stOpts)

	assert.NotNil(t, store)

	data := []byte("SomeBytes")
	writeData := bytes.NewReader(data)
	err := store.writeStream("SpecialPictures", writeData)
	if err != nil {
		t.Fatalf("Error writing stream: %v", err)
	}

	if ok := store.Has("SpecialPictures"); !ok {
		t.Fatalf("Expected to have SpecialPictures Folder")
	}

	r, err := store.Read("SpecialPictures")
	if err != nil {
		t.Fatalf("Error reading stream: %v", err)
	}

	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Fatalf("Expected %s, got %s", data, b)
	}
}

func TestStoreDelete(t *testing.T) {
	stOps := NewStoreOpts(
		CASPathTransformFunc,
	)
	store := NewStore(stOps)

	data := []byte("SomeBytes")
	writeData := bytes.NewReader(data)
	err := store.writeStream("SpecialPictures", writeData)
	if err != nil {
		t.Fatalf("Error writing stream: %v", err)
	}
	err = store.Delete("SpecialPictures")
	if err != nil {
		t.Fatalf("Error deleting stream: %v", err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	path := "SpecialPictures"
	pathKey := CASPathTransformFunc(path)
	expectedPathName := "0471c/0f32e/e0383/21c2b/50fbd/5e832/ede43/6c7ca"
	expectedOriginalKey := "0471c0f32ee038321c2b50fbd5e832ede436c7ca"
	if pathKey.PathName != expectedPathName {
		t.Fatalf("Expected %s, got %s", expectedPathName, pathKey.PathName)
	}
	if pathKey.FileName != expectedOriginalKey {
		t.Fatalf("Expected %s, got %s", expectedOriginalKey, pathKey.FileName)
	}
}

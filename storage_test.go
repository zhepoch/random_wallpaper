package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadPhotoQueryKey(t *testing.T) {
	hostDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(filepath.Join(hostDir, ".rwallpaper", "config"))
	if err != nil {
		t.Fatal(err)
	}

	key, err := ReadPhotoQueryKey()
	if err == nil {
		t.Fatal(key, err)
	}

	t.Log("if not file, must return error")
}

func TestSavePhotoQueryKey(t *testing.T) {
	err := SavePhotoQueryKey("urban")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("save key successful")
}

func TestReadPhotoQueryKeyAfterSave(t *testing.T) {
	key, err := ReadPhotoQueryKey()
	if err != nil {
		t.Fatal(key, err)
	}

	if key == "urban" {
		t.Log("read key successsful, ", key)
	} else {
		t.Fatal("get failed key: ", key)
	}

}

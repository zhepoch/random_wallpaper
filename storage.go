package main

import (
	"os"
	"path/filepath"

	"bitbucket.org/zhepoch/utilGo/fileUtil"
)

func SavePhotoQueryKey(key string) error {
	log.Debugf("Save photo query key to file: %v", key)
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	file, err := fileUtil.OpenFile(filepath.Join(homePath, ".rwallpaper", "config"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.WriteString(key)
	return err
}

func ReadPhotoQueryKey() (key string, err error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	file, err := fileUtil.OpenFile(filepath.Join(homePath, ".rwallpaper", "config"), os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 128)
	n, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:n]), nil
}

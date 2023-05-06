package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"bitbucket.org/zhepoch/utilGo/fileUtil"
)

func getConfigPath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homePath, ".rwallpaper", "config")
	return configPath, nil
}

func SavePhotoQueryKey(key string) error {
	log.Debugf("Save photo query key to file: %v", key)
	configPath, err := getConfigPath()
	file, err := fileUtil.OpenFile(configPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
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
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	fileBody, err := ioutil.ReadFile(configPath)
	return string(fileBody), err
}

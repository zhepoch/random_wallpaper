package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	UnsplashAPI = "https://api.unsplash.com"
	GetRandomPhotoAPI = "/photos/random"
)

type PhotoUrl struct {
	Full string `json:"full"`
}

type PhotoInfo struct {
	Id string `json:"id"`
	Urls PhotoUrl `json:"urls"`
}

func GetRandomPhoto() (string, error) {
	url := fmt.Sprintf("%s%s", UnsplashAPI, GetRandomPhotoAPI)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	body, err := Done(req)
	if err != nil {
		return "", err
	}

	var resp PhotoInfo
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	return DownloadPhoto(resp)
}

func DownloadPhoto(photoInfo PhotoInfo) (string, error) {
	req, err := http.NewRequest(http.MethodGet, photoInfo.Urls.Full, nil)
	if err != nil {
		return "", err
	}

	body, err := Done(req)

	var filePath = *FilePath
	if filePath[len(filePath) - 1] == '/' {
		filePath = filePath[:len(filePath) - 1]
	}
	fileName := fmt.Sprintf("%s/%s", filePath, photoInfo.Id)
	err = ioutil.WriteFile(fileName, body, os.FileMode(0644))
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func Done(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", *AccessKey))

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
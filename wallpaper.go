package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func GetDesktopCount() int {
	countStr, err := OsaScript(GetDesktopCountCommand)
	if err != nil {
		log.Debugln("get desktop error:", err)
		return 0
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Debugln("parse int error:", err)
		return 0
	}
	return count
}

func ApplyWallpaper(picturePath string, desktopIndex int) error {
	script := fmt.Sprintf(ApplyDesktopCommand, desktopIndex+1, picturePath)
	_, err := OsaScript(script)
	if err != nil {
		log.Debugf("Apply wallpaper got error: %v", err)
	}

	log.Debugln("run osasccript: ", script)
	return err
}

func ChangeWallPaper(photoQueryKey string) {
	desktopCount := GetDesktopCount()

	photoInfoList, err := GetRandomPhoto(desktopCount, photoQueryKey)
	if err != nil {
		log.Errorf("Get random photo got error: %v", err)
		return
	}

	log.Debugln("got photo number:", len(photoInfoList))

	var wg sync.WaitGroup
	lastWallPaper = []string{}
	for i := 0; i < len(photoInfoList); i++ {
		lastWallPaper = append(lastWallPaper,
			fmt.Sprintf("%s.jpg", photoInfoList[i].Id))

		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			log.Debugln("work", index, "starting...")
			photoPath, err := DownloadPhoto(photoInfoList[index])
			if err != nil {
				return
			}

			log.Infof("Get photo: %s", photoPath)
			err = ApplyWallpaper(photoPath, index)
			if err != nil {
				_ = RemoveFile(photoPath)
				log.Errorf("Apply wallpaper got error: %v", err)
			}
			log.Debugln("work", index, "finished.")
		}(i)
	}

	wg.Wait()
}

func RemoveExtraFile() {
	dir, err := os.Open(*FilePath)
	if err != nil {
		log.Errorf("open file, got error: %v", err)
		return
	}

	fileInfoList, err := dir.Readdir(-1)
	_ = dir.Close()
	if err != nil {
		log.Errorf("read dir, got error: %v", err)
		return
	}

	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			if err := os.RemoveAll(fmt.Sprintf("%s/%s", *FilePath, fileInfo.Name())); err != nil {
				log.Errorf("remove dir error: %v", err)
			}
		} else {
			if !Contains(fileInfo.Name(), lastWallPaper) {
				if err := os.Remove(fmt.Sprintf("%s/%s", *FilePath, fileInfo.Name())); err != nil {
					log.Errorf("remove file error: %v", err)
				}
			}
		}
	}
	return
}

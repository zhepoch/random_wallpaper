package main

import (
	"fmt"
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
	return err
}

func ChangeWallPaper() {
	desktopCount := GetDesktopCount()

	photoInfoList, err := GetRandomPhoto(desktopCount)
	if err != nil {
		log.Errorf("Get random photo got error: %v", err)
		return
	}

	log.Debugln("got photo number:", len(photoInfoList))

	var wg sync.WaitGroup
	for i := 0; i < len(photoInfoList); i++ {
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

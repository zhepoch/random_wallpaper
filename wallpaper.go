package main

import (
	"fmt"
	"strconv"
)

func GetDesktopCount() int {
	countStr, err := OsaScript(GetDesktopCountCommand)
	if err != nil {
		fmt.Println("get desktop error:", err)
		return 0
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		fmt.Println("parse int error:", err)
		return 0
	}
	return count
}


func ApplyWallpaper(picturePath string, desktopIndex int) error {
	script := fmt.Sprintf(ApplyDesktopCommand, desktopIndex+1, picturePath)
	_, err := OsaScript(script)
	return err
}


func ChangeWallPaper(index int) error {
	picturePath, err := GetRandomPhoto()
	if err != nil {
		return err
	}

	fmt.Println("get picture:", picturePath)
	return ApplyWallpaper(picturePath, index)
}
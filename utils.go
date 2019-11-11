package main

import (
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func Cmd(name string, args ...string) (string, error) {
	c := exec.Command(name, args...)
	out, err := c.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func OsaScript(command string) (string, error) {
	command = strings.TrimSpace(command)

	out, err := Cmd("osascript", "-e", command)
	if err != nil {
		return "", err
	}

	return out, nil
}

func Mkdir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RemoveFile(path string) error {
	err := os.Remove(path)
	return err
}

func Contains(val interface{}, arr interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		s := arr.([]string)
		for _, item := range s {
			if reflect.DeepEqual(val, item) {
				return true
			}
		}
	}
	return false
}

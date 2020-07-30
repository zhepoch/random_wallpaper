package main

import (
	"strings"

	"bitbucket.org/zhepoch/utilGo/processUtil"
)

func OsaScript(command string) (string, error) {
	command = strings.TrimSpace(command)

	out, err := processUtil.Cmd("osascript", "-e", command)
	if err != nil {
		return "", err
	}

	return out, nil
}

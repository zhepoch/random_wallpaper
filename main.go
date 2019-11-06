package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	AccessKey   = pflag.StringP("access_key", "a", "", "Access key of unsplash.")
	ReplaceTime = pflag.IntP("replace_time", "t", 5, "Change wallpaper every few minutes.")
	FilePath    = pflag.StringP("file_path", "p", "/tmp/random_wallpaper/", "save download wallpaper path.")
	LogLevel    = pflag.UintP("log_level", "v", 4, "debug level 0-5, 0:panic, 1:Fatal, 2:Error, 3:Warn, 4:Info 5:debug")
)

var (
	log = logrus.New()
)

func Init() {
	log.SetLevel(logrus.Level(*LogLevel))
	log.Formatter = &logrus.TextFormatter{FullTimestamp: true}
}

func main() {
	pflag.Parse()

	Init()

	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		log.Println("Please ctrl+c to stop!")
		<-c
		cancel()
	}(cancel)

	go func() {
		ChangeWallPaper()
	}()

	ticker := time.NewTicker(time.Minute * time.Duration(*ReplaceTime))
	for {
		select {
		case <-ticker.C:
			ChangeWallPaper()
		case <-ctx.Done():
			log.Println("quiting...")
			ticker.Stop()
			return
		}
	}
}

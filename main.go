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
	AccessKey     = pflag.StringP("access_key", "a", "", "Access key of unsplash.")
	ReplaceTime   = pflag.IntP("replace_time", "t", 5, "Change wallpaper every few minutes.")
	FilePath      = pflag.StringP("file_path", "p", "/tmp/random_wallpaper/", "save download wallpaper path.")
	LogLevel      = pflag.UintP("log_level", "v", 4, "debug level 0-5, 0:panic, 1:Fatal, 2:Error, 3:Warn, 4:Info 5:debug")
	PhotoQueryKey = pflag.StringP("phtot_query_key", "q", "", "Limit selection to photos matching a search term.")
)

var (
	lastWallPaper []string
	log           = logrus.New()
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
		RemoveExtraFile()
	}()

	removeExtraFileTicker := time.NewTicker(time.Hour * 5)
	go func(ctx context.Context) {
		for {
			select {
			case <-removeExtraFileTicker.C:
				log.Println("strating remove extra...")
				RemoveExtraFile()
			case <-ctx.Done():
				log.Println("remove extra file work quiting...")
				removeExtraFileTicker.Stop()
				return
			}
		}
	}(ctx)

	changeWallPaperTicker := time.NewTicker(time.Minute * time.Duration(*ReplaceTime))
	for {
		select {
		case <-changeWallPaperTicker.C:
			ChangeWallPaper()
		case <-ctx.Done():
			log.Println("change wallPaper work quiting...")
			changeWallPaperTicker.Stop()
			return
		}
	}
}

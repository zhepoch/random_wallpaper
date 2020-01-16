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
	FilePath      = pflag.StringP("file_path", "f", "/tmp/random_wallpaper/", "save download wallpaper path.")
	LogLevel      = pflag.UintP("log_level", "v", 4, "debug level 0-5, 0:panic, 1:Fatal, 2:Error, 3:Warn, 4:Info 5:debug")
	PhotoQueryKey = pflag.StringP("photo_query_key", "q", "", "Limit selection to photos matching a search term.")
	ListenPort    = pflag.IntP("listen_port", "p", 16606, "change get unsplash query key http server listen port.")
)

var (
	lastWallPaper []string
	log           = logrus.New()
	UserChange    = make(chan struct{}, 1)
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
		RemoveExtraFile()
		ChangeWallPaper()
	}()

	removeExtraFileTicker := time.NewTicker(time.Hour * 5)
	go func(ctx context.Context) {
		for {
			select {
			case <-removeExtraFileTicker.C:
				log.Println("Starting remove extra...")
				RemoveExtraFile()
			case <-ctx.Done():
				log.Println("Remove extra file work quiting...")
				removeExtraFileTicker.Stop()
				return
			}
		}
	}(ctx)

	go func() {
		log.Println("Starting listen addr: 127.0.0.1", *ListenPort)
		log.Errorln(ListenHttpServer())
	}()

	changeWallPaperTicker := time.NewTicker(time.Minute * time.Duration(*ReplaceTime))
	for {
		select {
		case <-UserChange:
			log.Println("Initiative Change Wallpaper...")
			ChangeWallPaper()
			changeWallPaperTicker.Stop()
			changeWallPaperTicker = time.NewTicker(time.Minute * time.Duration(*ReplaceTime))
		case <-changeWallPaperTicker.C:
			ChangeWallPaper()
		case <-ctx.Done():
			log.Println("Change wallPaper work quiting...")
			changeWallPaperTicker.Stop()
			return
		}
	}
}

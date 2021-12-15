package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	Version       = pflag.BoolP("version", "v", false, "show version info.")
	AccessKey     = pflag.StringP("access_key", "a", "", "Access key of unsplash.")
	ReplaceTime   = pflag.IntP("replace_time", "t", 5, "Change wallpaper every few minutes.")
	FilePath      = pflag.StringP("file_path", "f", "/tmp/random_wallpaper/", "Save download wallpaper path.")
	LogLevel      = pflag.UintP("log_level", "g", 4, "Debug level 0-5, 0:panic, 1:Fatal, 2:Error, 3:Warn, 4:Info 5:debug.")
	ProxyString   = pflag.StringP("proxy_url", "p", "", "Used in proxy to request unsplash API.")
	PhotoQueryKey = pflag.StringP("photo_query_key", "q", "", "Limit selection to photos matching a search term.")
	ListenPort    = pflag.IntP("listen_port", "l", 16606, "Change get unsplash query key http server listen port.")
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

	if *Version {
		marshaled, err := json.MarshalIndent(GetVersion(), "", " ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	Init()

	if *PhotoQueryKey == "" {
		key, err := ReadPhotoQueryKey()
		log.Debugf("ReadPhotoQueryKey: %v, %v\n", key, err)
		if err != nil {
			log.Errorf("Read Query key from file got error: %v\n", err)
		} else {
			*PhotoQueryKey = key
		}
	} else {
		err := SavePhotoQueryKey(*PhotoQueryKey)
		if err != nil {
			log.Errorf("Save Query key to file got error: %v", err)
		}
	}

	log.Infof("QueryKey: %s", *PhotoQueryKey)

	ctx, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		log.Println("Please ctrl+c to stop!")
		<-c
		cancel()
		<-c
		os.Exit(0)
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

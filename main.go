package main

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	AccessKey = pflag.StringP("access_key", "a", "", "Access key of unsplash.")
	ReplaceTime = pflag.IntP("replace_time", "t", 5, "Change wallpaper every few minutes.")
	FilePath = pflag.StringP("file_path", "p", "/tmp/random_wallpaper/", "save download wallpaper path.")
)

func main() {
	pflag.Parse()

	err := Mkdir(*FilePath)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	desktopCount := GetDesktopCount()

	var wg sync.WaitGroup
	for i := 0; i < desktopCount; i++ {
		wg.Add(1)
		go func(ctx context.Context, index int) {
			defer wg.Done()

			fmt.Printf("work %d starting...\n", index)

			err := ChangeWallPaper(index)
			if err != nil {
				fmt.Println("change wallpaper get error:", err)
			}

			ticker := time.NewTicker(time.Minute * time.Duration(*ReplaceTime))

			select {
			case <- ticker.C:
				err := ChangeWallPaper(index)
				if err != nil {
					fmt.Println("change wallpaper get error:", err)
				}
			case <- ctx.Done():
				fmt.Println("quiting...")
				ticker.Stop()
				return
			}
		}(ctx, i)
	}

	wg.Add(1)
	go func(cancel context.CancelFunc) {
		defer wg.Done()

		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		fmt.Println("Ctrl + C to exit....")
		<- c
		cancel()
	}(cancel)

	wg.Wait()
}
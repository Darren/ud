package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func download(video string) error {
	cmd := exec.Command("youtube-dl", video)
	var err error

	out, err := cmd.StdoutPipe()
	errout, err := cmd.StderrPipe()

	go func() {
		io.Copy(os.Stdout, out)
	}()

	go func() {
		io.Copy(os.Stderr, errout)
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

// DB is all video files in current directory
var DB map[string]string

var videoFileRegex = regexp.MustCompile(".*?-([a-zA-Z0-9_-]*).mp4")

var videoURLRegex = regexp.MustCompile(".*?v=([a-zA-Z0-9_-]*)")

func loadVideos() error {
	DB = make(map[string]string)

	walk := func(path string, info os.FileInfo, err error) error {
		r := videoFileRegex.FindStringSubmatch(path)
		if len(r) == 2 {
			videoID := r[1]
			DB[videoID] = path
		}
		return nil
	}
	return filepath.Walk(".", walk)
}

func usage() {
	fmt.Printf(" Usage: %s URL \n", os.Args[0])
	os.Exit(1)

}
func main() {
	if len(os.Args) < 2 {
		usage()
	}

	err := loadVideos()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Current video count: %d", len(DB))

	um := videoURLRegex.FindStringSubmatch(os.Args[1])
	if len(um) == 0 {
		log.Printf("Bad url: %s", os.Args[1])
		usage()
	}

	videoID := um[1]
	log.Printf("VideoID: %s", videoID)

	if file, ok := DB[videoID]; ok {
		log.Fatalf("Video [%s] exists: %s", videoID, file)
	}

	err = download(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

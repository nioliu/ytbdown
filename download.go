package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	util "github.com/nioliu/alfred"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const formatDefault = "default"
const typeDefault = "mp4"
const qualityDefault = "medium"

func downloadSingle() {
	if *video == "" {
		util.PrintError(errors.New("lack video and format"))
	}
	v := new(youtube.Video)
	if err := json.Unmarshal([]byte(*video), v); err != nil {
		util.PrintError(errors.New(fmt.Sprintf("err: %s, args: %s", err.Error(), *video)))
	}

	f := new(youtube.Format)
	if *format != formatDefault && *format != "" {
		if err := json.Unmarshal([]byte(*format), f); err != nil {
			util.PrintError(err)
		}
	} else if *format == formatDefault {
		v.Formats = v.Formats.Quality(qualityDefault)
		v.Formats = v.Formats.Type(typeDefault)
		v.Formats.Sort()
		f = &v.Formats[0]
	} else {
		util.PrintError(errors.New("lack format"))
	}

	doDownload(v, f, *downloadDir, *videoName)
}

func downloadAll() {
	videos := new(youtube.Video)
	if *video == "" {
		videos = listInfo()
	} else {
		if err := json.Unmarshal([]byte(*video), videos); err != nil {
			util.PrintError(errors.New(fmt.Sprintf("err: %s, args: %s", err.Error(), *video)))
		}
	}

	group := sync.WaitGroup{}
	downloadDir := filepath.Join(*downloadDir, videos.Title) + "/"
	createDir(downloadDir)

	videos.Formats = videos.Formats.Type(typeDefault)
	videos.Formats.Sort()

	for i, format := range videos.Formats {
		group.Add(1)
		go func(f youtube.Format, i int) {
			doDownload(videos, &f, downloadDir, fmt.Sprintf("%s_%d", videos.Title, i))
			group.Done()
		}(format, i)
	}
	group.Wait()
	printRes(videos.Title, downloadDir, "")
}

func createDir(name string) {
	if err := os.MkdirAll(name, 0777); err != nil {
		util.PrintError(err)
	}
}

func doDownload(video *youtube.Video, format *youtube.Format, dir, name string) {
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		util.PrintError(err)
	}
	if name == "" {
		name = video.Title
	}
	f, err := CreateTempVideoFile(name, dir)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		util.PrintError(err)
	}
	_, err = io.Copy(f, stream)
	if err != nil {
		util.PrintError(err)
	}

	printRes(video.Title, dir, name)
}

func printRes(title string, dir string, name string) {
	res := util.AlfredWorkflowRsp{AlfredWorkflow: &util.AlfredWorkflow{
		Arg: fmt.Sprintf("video[%s] has downloaded into [%s] ", title, dir+name),
		Variables: map[string]string{
			"file_path":  dir + name,
			"video_name": title,
		},
	}}

	bytes, err := json.Marshal(res)
	if err != nil {
		util.PrintError(err)
	}
	fmt.Println(string(bytes))
}

func CreateTempVideoFile(name string, dir string) (*os.File, error) {
	if dir == "" {
		dir = "./"
	}
	return os.Create(dir + name + ".mp4")
}

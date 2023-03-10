package main

import (
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/pflag"
)

var videoUrl = new(string)
var downloadDir = new(string)
var client = youtube.Client{}
var videoName = new(string)
var video = new(string)
var format = new(string)

func main() {
	p := pflag.StringP("method", "m", "info", "Specify a func to run, has: info, downloadSingle, downloadAll")
	pflag.StringVarP(videoUrl, "video_url", "u", "", "Specify the video url")
	pflag.StringVarP(downloadDir, "download_dir", "d", "", "Specify the video downloadSingle dir")
	pflag.StringVarP(videoName, "video_name", "n", "", "Specify the video downloadSingle video name")
	pflag.StringVarP(video, "video", "v", "", "Specify the video downloadSingle")
	pflag.StringVarP(format, "format", "f", "", "Specify the video format")

	pflag.Parse()

	switch *p {
	case "info":
		listInfo()
	case "downloadSingle":
		downloadSingle()
	case "downloadAll":
		downloadAll()
	}
}

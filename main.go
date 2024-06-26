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

func init() {
	// 设置代理
	//proxyURL, err := url.Parse("http://your-proxy-server:port")
	//if err != nil {
	//	log.Fatal("error parsing proxy URL: %w", err)
	//}
	//client.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyURL()}

	client.ChunkSize = youtube.Size1Mb * 5
}
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

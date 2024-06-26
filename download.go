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
	"sort"
	"strings"
	"sync"
)

const formatDefault = "default"
const typeDefault = "mp4"
const qualityDefault = "medium"

func downloadSingle() {
	if *video == "" {
		util.PrintError(errors.New("lack video and format"))
		return
	}
	v := new(youtube.Video)
	if err := json.Unmarshal([]byte(*video), v); err != nil {
		util.PrintError(fmt.Errorf("err: %s, args: %s", err.Error(), *video))
		return
	}
	videoName = &v.Title

	if *format != formatDefault && *format != "" {
		f := new(youtube.Format)
		if err := json.Unmarshal([]byte(*format), f); err != nil {
			util.PrintError(err)
		}
		videoFileName := fmt.Sprintf("%s_video.%s", *videoName, getExtensionFromMimeType(f.MimeType))
		doDownload(v, f, *downloadDir, videoFileName)
	} else {
		// 按照视频分辨率从高到低排序
		sortFormatsByQuality(v.Formats)

		// 选择分辨率最高的格式
		videoFormat := &v.Formats[0]

		// 检查是否有音频
		if videoFormat.AudioChannels == 0 {
			group := sync.WaitGroup{}
			// 下载视频
			videoFileName := fmt.Sprintf("%s_video.%s", *videoName, getExtensionFromMimeType(videoFormat.MimeType))
			group.Add(1)
			go func() {
				defer group.Done()
				doDownload(v, videoFormat, *downloadDir, videoFileName)
			}()

			// 过滤出具有音频的格式
			formatsWithAudio := filterFormatsWithAudio(v.Formats)
			if len(formatsWithAudio) == 0 {
				util.PrintError(errors.New("no audio formats found"))
				return
			}

			// 选择音频质量最高的格式
			audioFormat := selectHighestQualityFormat(formatsWithAudio)
			audioFileName := fmt.Sprintf("%s_audio.%s", *videoName, getExtensionFromMimeType(audioFormat.MimeType))
			group.Add(1)
			go func() {
				defer group.Done()
				doDownload(v, audioFormat, *downloadDir, audioFileName)
			}()

			group.Wait()
		} else {
			// 直接下载带有音频的视频
			fileName := fmt.Sprintf("%s.%s", *videoName, getExtensionFromMimeType(videoFormat.MimeType))
			doDownload(v, videoFormat, *downloadDir, fileName)
		}
	}
}

func selectHighestQualityFormat(formats []youtube.Format) *youtube.Format {
	sort.SliceStable(formats, func(i, j int) bool {
		return formats[i].Bitrate > formats[j].Bitrate
	})
	return &formats[0]
}

func filterFormatsWithAudio(formats []youtube.Format) []youtube.Format {
	var filtered []youtube.Format
	for _, format := range formats {
		if format.AudioChannels > 0 {
			filtered = append(filtered, format)
		}
	}
	return filtered
}

func sortFormatsByQuality(formats []youtube.Format) {
	sort.SliceStable(formats, func(i, j int) bool {
		return qualityValue(formats[i].Quality) > qualityValue(formats[j].Quality)
	})
}

func qualityValue(quality string) int {
	qualityMap := map[string]int{
		"hd2160": 2160,
		"hd1440": 1440,
		"hd1080": 1080,
		"hd720":  720,
		"large":  480,
		"medium": 360,
		"small":  240,
		"tiny":   144,
	}
	return qualityMap[quality]
}

func getExtensionFromMimeType(mimeType string) string {
	mimeType = strings.Split(mimeType, ";")[0]
	switch mimeType {
	case "video/webm":
		return "webm"
	case "video/mp4":
		return "mp4"
	case "audio/webm":
		return "webm"
	case "audio/mp4":
		return "m4a"
	default:
		return "bin"
	}
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
	return os.Create(dir + name)
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	util "github.com/nioliu/alfred"
	"strconv"
)

func listInfo() *youtube.Video {
	if *videoUrl == "" {
		util.PrintError(errors.New("lack video url"))
	}

	video, err := client.GetVideo(*videoUrl)
	if err != nil {
		util.PrintError(err)
	}
	videoJson, _ := json.Marshal(video)

	u := new(util.Result)
	u.Variables = map[string]string{
		"base_url": *videoUrl,
	}

	u.Items = append(u.Items, &util.Item{
		Uid:   video.ID + "-random",
		Title: video.Title + "-" + "Download recommended video",
		Icon:  nil,
		Arg:   "continue...",
		Variables: map[string]string{
			"video":  string(videoJson),
			"format": "default",
		},
	}, &util.Item{
		Uid:   video.ID + "-all",
		Title: video.Title + "-" + "Download all the videos",
		Icon:  nil,
		Arg:   "continue...",
		Variables: map[string]string{
			"video":  string(videoJson),
			"format": "all",
		},
	})

	for _, format := range video.Formats {
		formatJson, _ := json.Marshal(format)

		u.Items = append(u.Items, &util.Item{
			Uid:   video.ID + strconv.Itoa(format.ItagNo),
			Type:  format.MimeType,
			Title: video.Title + "-" + video.Description,
			Subtitle: fmt.Sprintf("Audio: %d, Quality: %s, MimeType: %s, Duration: %sMs",
				format.AudioChannels, format.Quality, format.MimeType, format.ApproxDurationMs),
			Arg:          string(formatJson),
			Autocomplete: "",
			Icon:         nil,
			Variables: map[string]string{
				"video":         string(videoJson),
				"format":        string(formatJson),
				"quality_label": format.QualityLabel,
			},
			QuickLookUrl: format.URL,
			Text: &util.Text{
				Copy: format.URL,
			},
		})
	}

	bytes, err := json.Marshal(u)
	if err != nil {
		util.PrintError(err)
	}

	fmt.Print(string(bytes))
	return video
}

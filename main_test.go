package main

import (
	"encoding/json"
	"testing"
)

func Test_listInfo(t *testing.T) {
	*videoUrl = "https://www.youtube.com/watch?v=8yzqumXb3QA"
	info := listInfo()
	videoJson, _ := json.Marshal(info)
	*video = string(videoJson)
	*format = formatDefault
	downloadSingle()
}

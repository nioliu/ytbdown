package main

import (
	"encoding/json"
	"github.com/kkdai/youtube/v2"
	"testing"
)

func Test_createDir(t *testing.T) {
	createDir("/Users/nioliu/Movies/youtube/nioliu/")
}

func TestUnmarhsal(t *testing.T) {
	raw := `{"ID":"mgaJYCIIj7A","Title":"A380"}`
	y := new(youtube.Video)
	err := json.Unmarshal([]byte(raw), y)
	t.Log(y)
	t.Log(err)
}

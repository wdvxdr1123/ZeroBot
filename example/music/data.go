package music

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

func queryNeteaseMusic(musicName string) int64 {
	client := http.Client{}
	req, err := http.NewRequest("GET", "http://music.163.com/api/search/get?type=1&s="+url.QueryEscape(musicName), nil)
	if err != nil {
		return 0
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.66")
	res, err := client.Do(req)
	if err != nil {
		return 0
	}
	data, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return 0
	}
	return gjson.ParseBytes(data).Get("result.songs.0.id").Int()
}

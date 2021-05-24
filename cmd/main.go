package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/suconghou/mediaindex"
	"github.com/suconghou/youtubevideoparser"
)

var client = &http.Client{
	Timeout: time.Minute,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func main() {
	testWebm()
}

func parseMp4() error {
	bs, err := ioutil.ReadFile("/Users/admin/Downloads/140.mp4")
	if err != nil {
		return err
	}
	// indexRange: {start: "632", end: "1023"}
	data := bs[632:1023]
	// 包含632 但不包含1023,所以391字节
	fmt.Println(len(data))
	fmt.Printf("%x\n", md5sum(data))
	p, err := mediaindex.ParseMp4(data)
	fmt.Println(p, err)
	return nil
}

func testWebm() {
	arr := []string{
		"HzOjwL7IP_o",
	}
	for _, id := range arr {
		url := fmt.Sprintf("https://stream.pull.workers.dev/video/%s.json", id)
		err := parseWebm(url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("==============")
		}
	}

}

func parseWebm(url string) error {
	data, indexEndOffset, totalSize, err := Parse(url, "243")
	if err != nil {
		return err
	}
	fmt.Println(url, len(data))
	fmt.Printf("%x\n", md5sum(data))
	p, err := mediaindex.ParseWebM(data, indexEndOffset, totalSize)
	fmt.Println(p, err)

	data, indexEndOffset, totalSize, err = Parse(url, "249")
	if err != nil {
		return err
	}
	fmt.Println(url, len(data))
	fmt.Printf("%x\n", md5sum(data))
	p, err = mediaindex.ParseWebM(data, indexEndOffset, totalSize)
	fmt.Println(p, err)

	return nil
}

func md5sum(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func Parse(url string, itag string) ([]byte, uint64, uint64, error) {
	bs, err := Get(url)
	if err != nil {
		return nil, 0, 0, err
	}
	var data youtubevideoparser.VideoInfo
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return nil, 0, 0, err
	}
	item := data.Streams[itag]
	if item == nil || item.IndexRange == nil || item.ContentLength == "" {
		return nil, 0, 0, fmt.Errorf("error info %s %s", url, itag)
	}
	var indexEndOffset uint64
	var totalSize uint64
	indexEndOffset, err = strconv.ParseUint(item.IndexRange.End, 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}
	totalSize, err = strconv.ParseUint(item.ContentLength, 10, 64)
	if err != nil {
		return nil, 0, 0, err
	}
	uri := fmt.Sprintf("https://stream.pull.workers.dev/video/%s/%s/%s-%s.ts", data.ID, itag, item.IndexRange.Start, item.IndexRange.End)
	index, err := Get(uri)
	if err != nil {
		return nil, 0, 0, err
	}
	return index, indexEndOffset, totalSize, nil
}

// Get http data, the return value should be readonly
func Get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s:%s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

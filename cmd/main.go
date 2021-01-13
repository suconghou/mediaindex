package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"

	"github.com/suconghou/mediaindex"
)

func main() {
	err := parseWebm()
	if err != nil {
		fmt.Println(err)
	}
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
	var p = mediaindex.ParseMp4(data)
	fmt.Println(p)
	return nil
}

func parseWebm() error {
	bs, err := ioutil.ReadFile("/Users/admin/Downloads/249.webm")
	if err != nil {
		return err
	}
	// indexRange: {start: "219", end: "1228"}
	data := bs[259:764]
	// 包含632 但不包含1023,所以391字节
	fmt.Println(len(data))
	fmt.Printf("%x\n", md5sum(data))
	var p = mediaindex.ParseWebM(data, 764, 1727516)
	fmt.Println(p)
	return nil
}

func md5sum(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

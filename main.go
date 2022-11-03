package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	test()
}

func test() {
	const url = "https://www.taiwan.net.tw/m1.aspx?sNo=0012076"

	imgs := crawlerImgUrl(url)

	if imgs == nil {
		return
	}

	for _, imgUrl := range imgs {
		if strings.Contains(imgUrl, "http") {
			if err := saveImgContentByHttp(imgUrl); err != nil {
				fmt.Println(err.Error())
				continue
			}
		} else if strings.Contains(imgUrl, "data:image") {
			if err := saveImgContentByData(imgUrl); err != nil {
				fmt.Println(err.Error())
				continue
			}

		}
	}

	//saveImgContentByHttp("https://a.deviantart.net/avatars-big/w/l/wlop.jpg?10")
	//saveImgContentByData("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=")
	return
}
func crawlerImgUrl(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	data, err := io.ReadAll(resp.Body)
	content := string(data)

	r, err := regexp.Compile("<img.*?src=[\"|'](.*?)[\"|']")

	if err != nil {
		return nil
	}

	s := r.FindAllStringSubmatch(content, -1)

	imgMap := make(map[string]bool)
	var imgs []string
	for _, target := range s {
		if imgMap[target[1]] {
			continue
		}
		imgMap[target[1]] = true
		imgs = append(imgs, target[1])
	}
	return imgs
}

func saveImgContentByHttp(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	contentType := resp.Header["Content-Type"][0]
	return saveImage(data, url[len(url)-5:len(url)], contentType)
}

func saveImgContentByData(data string) error {
	r, err := regexp.Compile("data:(image/.*);base64,(.*)")
	if err != nil {
		return err
	}

	s := r.FindAllStringSubmatch(data, -1)
	ext := s[0][1]
	content := s[0][2]
	ss, err := base64.StdEncoding.DecodeString(content)

	if err != nil {
		return err
	}

	return saveImage(ss, content[0:5], ext)
}

func saveImage(data []byte, fileName string, contentType string) error {
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = "jpeg"
	case "image/png":
		ext = "png"
	default:
		ext = ""
	}
	if ext == "" {
		return errors.New(fmt.Sprintf("not support this content type: %s", contentType))
	}
	err := os.MkdirAll("imgs", 0755)
	if err != nil {
		return errors.New("mkdir imgs failed")
	}
	fullName := fmt.Sprintf("%s/%s.%s", "imgs", strings.Replace(fileName, "?", "", -1), ext)

	if err := os.WriteFile(fullName, data, 0666); err != nil {
		return err
	}
	return nil
}

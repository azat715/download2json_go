package main

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"sync"
)

type Album struct {
	Userid int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
}

type Photo struct {
	Albumid      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	Thumbnailurl string `json:"thumbnailUrl"`
}

type PhotoFile struct {
	Photo
	path string
}

var albumsc = make(chan map[int]Album)
var photosc = make(chan []Photo)
var photo_file = make(chan PhotoFile)

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}

func get_albums() {
	body, _ := get("https://jsonplaceholder.typicode.com/albums/")
	albums := []Album{}
	json.Unmarshal(body, &albums)
	var m map[int]Album
	m = make(map[int]Album)
	for _, v := range albums {
		m[v.Id] = v
	}
	albumsc <- m
}

func get_photos() {
	body, _ := get("https://jsonplaceholder.typicode.com/albums/")
	photos := []Photo{}
	json.Unmarshal(body, &photos)
	photosc <- photos
}

func prepare_photo_files() {
	albums := <-albumsc
	photos := <-photosc
	for _, v := range photos {
		album := albums[v.Albumid]
		c := PhotoFile{
			Photo: v,
			// /<album.title>/<photo.title>.png
			path: filepath.Join(album.Title, v.Title, ".png"),
		}
		photo_file <- c
	}
}

func worker_poll() {
	var wg sync.WaitGroup
	workerPoolSize := 100

}

func main() {
	go get_albums()
	go get_photos()
	go prepare_photo_files()
	// serialize_tese()
}

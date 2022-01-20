package main

import (
	"download2json/internal/models"
	"download2json/internal/utils"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

const FOLDER string = "albums"
const ALBUM_URL string = "https://jsonplaceholder.typicode.com/albums/"
const PHOTOS_URL string = "https://jsonplaceholder.typicode.com/photos/"

var albumsc = make(chan map[int]models.Album)
var photosc = make(chan []models.Photo)

func get_albums() {
	body, _ := utils.Get(ALBUM_URL)
	var album models.Album
	albums, err := album.Serialize(body)
	if err != nil {
		log.Fatal("произошла ошибка сериализации albums")
	}
	var m map[int]models.Album
	m = make(map[int]models.Album)
	for _, v := range albums {
		m[v.Id] = v
	}
	albumsc <- m
}

func get_photos() {
	body, _ := utils.Get(PHOTOS_URL)
	var photo models.Photo
	photos, err := photo.Serialize(body)
	if err != nil {
		log.Fatal("произошла ошибка сериализации photos")
	}
	photosc <- photos
}

func worker_poll() {
	var wg sync.WaitGroup
	workerPoolSize := 100

	jobCh := make(chan models.PhotoFile, workerPoolSize)

	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range jobCh {
				file, errDownload := task.Download()
				if errDownload != nil {
					log.Fatal("Ошибка при скачивании фотографии")
				}
				album_path := filepath.Join(FOLDER, task.Album.Title)
				utils.Create_folder(album_path)
				photo_path := filepath.Join(album_path, task.Photo.Title+".png")
				errWrite := ioutil.WriteFile(photo_path, file, 0644)
				if errWrite != nil {
					log.Fatal("Ошибка при записи файла")
				}
			}
		}()
	}
	albums := <-albumsc
	photos := <-photosc
	for _, v := range photos {
		c := models.PhotoFile{
			Photo: v,
			Album: albums[v.Albumid],
		}
		jobCh <- c
	}
	close(jobCh)
	wg.Wait()

}

func main() {
	go get_albums()
	go get_photos()
	utils.Create_folder(FOLDER)
	worker_poll()
}

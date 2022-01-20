package core

import (
	l "download2json/internal/logger"
	"download2json/internal/models"
	"download2json/internal/utils"
	"io/ioutil"
	"path/filepath"
	"sync"
)

var albumsc = make(chan map[int]models.Album)
var photosc = make(chan []models.Photo)

func Get_albums(url string) {
	body, _ := utils.Get(url)
	var album models.Album
	albums, err := album.Serialize(body)
	if err != nil {
		l.ErrorLogger.Fatalln("произошла ошибка сериализации albums")
	}
	var m map[int]models.Album
	m = make(map[int]models.Album)
	for _, v := range albums {
		m[v.Id] = v
	}
	albumsc <- m
}

func Get_photos(url string) {
	body, _ := utils.Get(url)
	var photo models.Photo
	photos, err := photo.Serialize(body)
	if err != nil {
		l.ErrorLogger.Fatalln("произошла ошибка сериализации photos")
	}
	photosc <- photos
}

func Worker_poll(folder string) {
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
					l.ErrorLogger.Fatalln("Ошибка при скачивании фотографии")
				}
				album_path := filepath.Join(folder, task.Album.Title)
				utils.Create_folder(album_path)
				photo_path := filepath.Join(album_path, task.Photo.Title+".png")
				errWrite := ioutil.WriteFile(photo_path, file, 0644)
				if errWrite != nil {
					l.ErrorLogger.Fatalln("Ошибка при записи файла")
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

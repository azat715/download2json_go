package core

import (
	l "download2json/internal/logger"
	"download2json/internal/models"
	"download2json/internal/utils"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type Result struct {
	albums map[int]models.Album
	photos []models.Photo
	err    error
}

type TaskPhoto struct {
	photo_path        string
	album_folder_name string
	url               string
}

type Executor interface {
	Execute(url string, result chan Result)
}

type AlbumStrategy struct{} // конкретная стратегия альбомы

func (AlbumStrategy) Execute(url string, result chan Result) {
	body, _ := utils.Get(url)
	var album models.Album
	albums, err := album.Serialize(body)
	if err != nil {
		result <- Result{err: err}
	}
	var m map[int]models.Album
	m = make(map[int]models.Album)
	for _, i := range albums {
		m[i.Id] = i
	}
	result <- Result{albums: m}
}

type PhotoStrategy struct{} // конкретная стратегия photos

func (PhotoStrategy) Execute(url string, result chan Result) {
	body, _ := utils.Get(url)
	var photo models.Photo
	photos, err := photo.Serialize(body)
	if err != nil {
		result <- Result{err: err}
	}
	result <- Result{photos: photos}
	l.GeneralLogger.Println("Закончилосб скачивание альбомов")
}

type PhotoBinStrategy struct{}

func (PhotoBinStrategy) Execute(task TaskPhoto, result chan Result) {

}

type Context struct {
	Executor
}

func (c *Context) Download(url string, result chan Result) {
	c.Execute(url, result)
}

func Get_albums(url string, res chan Result) {
	l.GeneralLogger.Println("Началось скачивание альбомов")
	c := Context{AlbumStrategy{}}
	c.Download(url, res)
}

func Get_photos(url string, res chan Result) {
	l.GeneralLogger.Println("Началось скачивание фоток")
	c := Context{PhotoStrategy{}}
	c.Download(url, res)

}

func DownloadAll(album_url string, photos_url string, folder string) {

	var albumsc = make(chan Result)
	var photosc = make(chan Result)
	go Get_albums(album_url, albumsc)
	go Get_photos(photos_url, photosc)

	workerPoolSize := 100
	jobCh := make(chan TaskPhoto, workerPoolSize)

	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range jobCh {
				l.GeneralLogger.Println("Началось скачивание фото")
				body, errDownload := utils.Get(task.url)
				if errDownload != nil {
					l.GeneralLogger.Printf("Ошибка %v\n", errDownload)
				}
				utils.Create_folder(task.album_folder_name)
				errWrite := ioutil.WriteFile(task.photo_path, body, 0644)
				if errWrite != nil {
					l.GeneralLogger.Printf("Ошибка %v\n", errWrite)
				}
			}
		}()
	}
	albums := <-albumsc
	photos := <-photosc
	for _, photo := range photos.photos {
		album := albums.albums[photo.Albumid]
		album_folder_name := filepath.Join(folder, album.Title)
		task := TaskPhoto{
			album_folder_name: album_folder_name,
			photo_path:        filepath.Join(album_folder_name, photo.Title+".png"),
			url:               photo.Url,
		}
		jobCh <- task
	}

	wg.Wait()
	close(jobCh)
}

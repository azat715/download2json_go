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

type Task struct {
	photo_path        string
	album_folder_name string
	url               string
}

type Executor interface {
	Execute(task Task, result chan []Result)
}

type AlbumStrategy struct{} // конкретная стратегия альбомы

func (AlbumStrategy) Execute(task Task, result chan []Result) {
	var res []Result
	body, _ := utils.Get(task.url)
	var album models.Album
	albums, err := album.Serialize(body)
	if err != nil {
		result <- append(res, Result{err: err})
	}
	var m map[int]models.Album
	m = make(map[int]models.Album)
	for _, i := range albums {
		m[i.Id] = i
	}
	result <- append(res, Result{albums: m})
}

type PhotoStrategy struct{} // конкретная стратегия photos

func (PhotoStrategy) Execute(task Task, result chan []Result) {
	var res []Result
	body, _ := utils.Get(task.url)
	var photo models.Photo
	photos, err := photo.Serialize(body)
	if err != nil {
		result <- append(res, Result{err: err})
	}
	result <- append(res, Result{photos: photos})
	l.GeneralLogger.Println("Закончилосб скачивание альбомов")
}

type PhotoBinStrategy struct{}

func (PhotoBinStrategy) Execute(task Task, result chan []Result) {
	var res []Result
	body, errDownload := utils.Get(task.url)
	if errDownload != nil {
		result <- append(res, Result{err: errDownload})
		l.GeneralLogger.Printf("Ошибка %v\n", errDownload)
	}
	utils.Create_folder(task.album_folder_name)
	errWrite := ioutil.WriteFile(task.photo_path, body, 0644)
	if errWrite != nil {
		result <- append(res, Result{err: errWrite})
		l.GeneralLogger.Printf("Ошибка %v\n", errWrite)
	}
}

type Context struct {
	Executor
}

func (c *Context) Download(task Task, result chan []Result) {
	c.Execute(task, result)
}

func Get_albums(url string, res chan []Result) {
	l.GeneralLogger.Println("Началось скачивание альбомов")
	c := Context{AlbumStrategy{}}
	task := Task{url: url}
	c.Download(task, res)
}

func Get_photos(url string, res chan []Result) {
	l.GeneralLogger.Println("Началось скачивание фоток")
	c := Context{PhotoStrategy{}}
	task := Task{url: url}
	c.Download(task, res)

}

func DownloadAll(album_url string, photos_url string, folder string) {

	var albumsCh = make(chan []Result)
	var photosCh = make(chan []Result)
	go Get_albums(album_url, albumsCh)
	go Get_photos(photos_url, photosCh)

	workerPoolSize := 100
	jobCh := make(chan Task, workerPoolSize)
	results := make(chan []Result, workerPoolSize)

	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range jobCh {
				l.GeneralLogger.Println("Началось скачивание фото")
				c := Context{PhotoBinStrategy{}}
				c.Download(task, results)
			}
		}()
	}
	albumsRes := <-albumsCh
	photosRes := <-photosCh
	photos := albumsRes[0]
	albums := photosRes[0]
	for _, photo := range photos.photos {
		album := albums.albums[photo.Albumid]
		album_folder_name := filepath.Join(folder, album.Title)
		task := Task{
			album_folder_name: album_folder_name,
			photo_path:        filepath.Join(album_folder_name, photo.Title+".png"),
			url:               photo.Url,
		}
		jobCh <- task
	}
	results2 := <-results
	for res := range results2 {
		l.ErrorLogger.Printf("Ошибка %v\n", res)
	}
	wg.Wait()
	close(jobCh)
	close(results)
}

package core

import (
	"context"
	l "download2json/internal/logger"
	"download2json/internal/models"
	"download2json/internal/utils"
	"download2json/internal/wpool"
	"fmt"
	"io/ioutil"
)

const workerCount int = 100

type Result struct {
	data interface{}
}

type Problem interface {
	Get_Task() Task
}

type Task struct {
	photo_path        string
	album_folder_name string
	url               string
}

type Executor interface {
	Execute(task Task) (Result, error)
}

type AlbumStrategy struct{} // конкретная стратегия альбомы

func (AlbumStrategy) Execute(task Task) (Result, error) {
	body, _ := utils.Get(task.url)
	var album models.Album
	albums, err := album.Serialize(body)
	if err != nil {
		return Result{}, err
	}
	var m map[int]models.Album
	m = make(map[int]models.Album)
	for _, i := range albums {
		m[i.Id] = i
	}
	return Result{data: m}, nil
}

type PhotoStrategy struct{} // конкретная стратегия photos

func (PhotoStrategy) Execute(task Task) (Result, error) {
	body, _ := utils.Get(task.url)
	var photo models.Photo
	photos, err := photo.Serialize(body)
	if err != nil {
		return Result{}, err
	}
	return Result{data: photos}, nil
}

type PhotoBinStrategy struct{}

func (PhotoBinStrategy) Execute(task Task) (Result, error) {
	body, errDownload := utils.Get(task.url)
	if errDownload != nil {
		return Result{}, errDownload
	}
	utils.Create_folder(task.album_folder_name)
	errWrite := ioutil.WriteFile(task.photo_path, body, 0644)
	if errWrite != nil {
		return Result{}, errDownload
	}
	return Result{data: "Загрузка завершена"}, nil
}

type Context struct {
	Executor
}

func (c *Context) Download(task Task) (Result, error) {
	return c.Execute(task)
}

func DownloadAll(album_url string, photos_url string, folder string) {
	wp := wpool.New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.GenerateFrom(testJobs())

	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			fmt.Println(r)
		case <-wp.Done:
			return
		default:
		}
	}

}

var TestValue wpool.ExecutionFn = func(ctx context.Context, args interface{}) (interface{}, error) {
	l.GeneralLogger.Println("Началось скачивание альбомов")
	c := Context{AlbumStrategy{}}
	task := Task{url: args}
	return c.Download(task)
}

func Get_album(ctx context.Context, args string) (Result, error) {
	l.GeneralLogger.Println("Началось скачивание альбомов")
	c := Context{AlbumStrategy{}}
	task := Task{url: args}
	return c.Download(task)
}

func Get_photos(ctx context.Context, args string) (Result, error) {
	l.GeneralLogger.Println("Началось скачивание фоток")
	c := Context{PhotoStrategy{}}
	task := Task{url: args}
	return c.Download(task)
}

func testJobs() []wpool.Job {
	jobs := make([]wpool.Job, 2)
	ctx := context.Background()

	jobs[0] = wpool.Job{
		Descriptor: wpool.JobDescriptor{
			ID:       wpool.JobID(fmt.Sprintf("%v", 0)),
			JType:    "Album",
			Metadata: nil,
		},
		ExecFn: wpool.ExecutionFn(Get_album),
		Args:   "https://jsonplaceholder.typicode.com/albums/",
	}
	jobs[1] = wpool.Job{
		Descriptor: wpool.JobDescriptor{
			ID:       wpool.JobID(fmt.Sprintf("%v", 0)),
			JType:    "Photo",
			Metadata: nil,
		},
		ExecFn: Get_photos,
		Args:   "https://jsonplaceholder.typicode.com/photos/",
	}
	return jobs
}

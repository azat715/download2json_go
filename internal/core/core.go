package core

import (
	"context"
	l "download2json/internal/logger"
	"download2json/internal/models"
	"download2json/internal/wpool"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Result struct {
	albums []models.Album
	photos []models.Photo
}

func Create_folder(folder string) {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folder, 0755)
		if errDir != nil {
			l.ErrorLogger.Fatalln(err)
		}
	}
}

func Get(url string) ([]byte, error) {
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

const workerCount int = 100
const jobsCount int = 100

func DownloadJson(album_url string, photo_url string) Result {
	wp := wpool.New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	jobs := make([]wpool.Job, 2)
	jobs[0] = wpool.Job{
		Descriptor: wpool.JobDescriptor{
			ID:       wpool.JobID(fmt.Sprintf("%v", 0)),
			JType:    "Album",
			Metadata: nil,
		},
		ExecFn: Get,
		Args:   album_url,
	}
	jobs[1] = wpool.Job{
		Descriptor: wpool.JobDescriptor{
			ID:       wpool.JobID(fmt.Sprintf("%v", 1)),
			JType:    "Photo",
			Metadata: nil,
		},
		ExecFn: Get,
		Args:   photo_url,
	}

	go wp.GenerateFrom(jobs)

	go wp.Run(ctx)

	res := new(Result)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			if r.Descriptor.JType == "Album" {
				var album models.Album
				albums, err := album.Serialize(r.Value)
				if err != nil {
					l.ErrorLogger.Panicln("Ошибка при сериализации альбомов")
				}
				l.GeneralLogger.Print("закончилось скачивание альбомов")
				res.albums = albums
			}
			if r.Descriptor.JType == "Photo" {
				var photo models.Photo
				photos, err := photo.Serialize(r.Value)
				if err != nil {
					l.ErrorLogger.Panicln("Ошибка при сериализации фотографий")
				}
				l.GeneralLogger.Print("закончилось скачивание фотографий")
				res.photos = photos
			}
		case <-wp.Done:
			return *res
		default:
		}
	}
}

func DownloadAll(folder string, album_url string, photo_url string) {
	res := DownloadJson(album_url, photo_url)
	wp := wpool.New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.GenerateFrom(photosJobs(folder, &res))

	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			i, err := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if err != nil {
				l.ErrorLogger.Fatalf("unexpected error: %v", err)
			}
			Create_folder(r.Descriptor.Metadata["albumPath"])
			errWrite := ioutil.WriteFile(r.Descriptor.Metadata["photoPath"], r.Value, 0644)
			if errWrite != nil {
				l.ErrorLogger.Printf("Произошла ошибка сохранения %s", errWrite)
			}
			l.GeneralLogger.Printf("Выполенение задачи %v", i)
		case <-wp.Done:
			return
		default:
		}
	}
}

func photosJobs(folder string, res *Result) []wpool.Job {
	jobs := make([]wpool.Job, 5000)

	var albumsDict map[int]models.Album
	albumsDict = make(map[int]models.Album)
	for _, i := range res.albums {
		albumsDict[i.Id] = i
	}
	for i, photo := range res.photos {
		fmt.Print(photo)
		albumTitle := albumsDict[photo.Albumid].Title
		albumPath := filepath.Join(folder, albumTitle)
		photoPath := filepath.Join(albumPath, photo.Title+".png")
		jobs[i] = wpool.Job{
			Descriptor: wpool.JobDescriptor{
				ID:    wpool.JobID(fmt.Sprintf("%v", i)),
				JType: "anyType",
				Metadata: map[string]string{
					"albumPath": albumPath,
					"photoPath": photoPath,
				},
			},
			ExecFn: Get,
			Args:   photo.Url,
		}

	}
	return jobs
}

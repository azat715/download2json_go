package core

import (
	"context"
	l "download2json/internal/logger"
	"download2json/internal/models"
	"download2json/internal/wpool"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const workerCount int = 1000
const jobsCount int = 100

type Photo struct {
	url  string
	path models.FilePath
}

func Create_folder(folder string) error {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folder, 0755)
		if errDir != nil {
			return errDir
		}
	} else {
		return err
	}
	return nil
}

func processingPhoto(ctx context.Context, args interface{}) (interface{}, error) {
	photo, ok := args.(Photo)
	if !ok {
		l.ErrorLogger.Print("Unknown type parameter processingPhoto")
		return "", NewErrorWrapper("processingPhoto", errors.New("Unknown type parameter processingPhoto"), "failed to convert type")
	}

	var err error
	raw, err := get(photo.url)
	if err != nil {
		return "", NewErrorWrapper("processingPhoto", err, fmt.Sprintf("Произошла ошибка при скачивании %s", photo.url))
	}
	err = Create_folder(photo.path.Folder)
	if err != nil {
		return "", NewErrorWrapper("processingPhoto", err, fmt.Sprintf("Произошла ошибка при создании каталога  %s", photo.path.Folder))
	}
	err = ioutil.WriteFile(photo.path.Path, raw, 0644)
	if err != nil {
		return "", NewErrorWrapper("processingPhoto", err, fmt.Sprintf("Произошла ошибка сохранения %s", photo.path.Path))
	}
	return fmt.Sprintf("Done processing %s", photo.url), nil
}

func get_albums_and_photos(ctx context.Context, album_url string, photos_url string) (map[int]models.Album, models.Photos, error) {
	var err error
	var albums models.Albums
	var photos models.Photos

	errorsCh := make(chan error)
	albumsCh := make(chan []byte)
	photosCh := make(chan []byte)

	go Download(album_url, albumsCh, errorsCh)
	l.GeneralLogger.Print("Start download albums")
	go Download(photos_url, photosCh, errorsCh)
	l.GeneralLogger.Print("Start download photos")

	select {
	case r := <-albumsCh:
		albums, err = albums.Parse(r)
		if err != nil {
			return albums.AsDict(), photos, NewErrorWrapper("get_albums_and_photos", err, "failed to marshal json albums")
		}
		l.GeneralLogger.Print("Finished download albums")
	case r := <-errorsCh:
		return albums.AsDict(), photos, NewErrorWrapper("get_albums_and_photos", r, "failed to download")
	case <-ctx.Done():
		return albums.AsDict(), photos, NewErrorWrapper("get_albums_and_photos", ctx.Err(), "Context done")
	}

	select {
	case r := <-photosCh:
		photos, err = photos.Parse(r)
		if err != nil {
			return albums.AsDict(), photos, NewErrorWrapper("get_albums_and_photos", err, "failed to marshal json photos")
		}
		l.GeneralLogger.Print("Finished download photos")
	case r := <-errorsCh:
		return albums.AsDict(), photos, NewErrorWrapper("get_albums_and_photos", r, "failed to download")
	case <-ctx.Done():
		return albums.AsDict(), photos, NewErrorWrapper("get_albums_and_photos", ctx.Err(), "Context done")
	}

	return albums.AsDict(), photos, nil
}

func Core(album_url string, photos_url string, folder string) {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	// ctx, cancel = context.WithTimeout(ctx, time.Duration(150)*time.Millisecond) через 150 мс все горутины отменятся
	defer cancel()

	albums, photos, err := get_albums_and_photos(ctx, album_url, photos_url)

	if err != nil {
		err = (NewErrorWrapper("get_albums_and_photos", err, fmt.Sprintf("Произошла ошибка загрузки")))
		var ew ErrorWrapper
		if errors.As(err, &ew) {
			l.ErrorLogger.Print(ew.Message)
			l.ErrorLogger.Print(ew.Context)
			panic(err)
		}

	}

	jobs := make([]wpool.Job, 5000)
	for i, photo := range photos {
		jobs[i] = wpool.Job{
			Descriptor: wpool.JobDescriptor{
				ID:       wpool.JobID(fmt.Sprintf("%v", i)),
				JType:    "Download",
				Metadata: nil,
			},
			ExecFn: wpool.ExecutionFn(processingPhoto),
			Args: Photo{
				url:  photo.Url,
				path: buildPhotoPath(albums[photo.Albumid].Title, photo.Title, folder),
			},
		}
	}

	wp := wpool.New(workerCount)

	go wp.GenerateFrom(jobs)

	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}
			i, _ := strconv.ParseInt(string(r.Descriptor.ID), 10, 64)
			if r.Err != nil {
				var ew ErrorWrapper
				if errors.As(r.Err, &ew) {
					l.ErrorLogger.Print(ew.Message)
					l.ErrorLogger.Print(ew.Context)
					l.ErrorLogger.Print(ew.Err)
				}
			} else {
				l.GeneralLogger.Print(r.Value)
			}
			l.GeneralLogger.Printf("Task: %v finished", i)
		case <-wp.Done:
			l.GeneralLogger.Print("Done")
			return
		}
	}

}

func buildPhotoPath(albumTitle string, photoTitle string, folder string) models.FilePath {
	albumPath := filepath.Join(folder, albumTitle)
	photoPath := filepath.Join(albumPath, photoTitle+"."+models.PhotoExt) // вместо точки для винды нужно использовать / возможно можно подставить какой нибудь универсальный разделитель
	return models.FilePath{
		Path:   photoPath,
		Folder: albumPath,
	}
}

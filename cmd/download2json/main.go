package main

import (
	"download2json/internal/core"
	l "download2json/internal/logger"
	"download2json/internal/utils"
)

const FOLDER string = "albums"
const ALBUM_URL string = "https://jsonplaceholder.typicode.com/albums/"
const PHOTOS_URL string = "https://jsonplaceholder.typicode.com/photos/"

func main() {
	l.GeneralLogger.Println("Началось скачивание")
	go core.Get_albums(ALBUM_URL)
	go core.Get_photos(PHOTOS_URL)
	utils.Create_folder(FOLDER)
	core.Worker_poll(FOLDER)
	l.GeneralLogger.Println("Закончилось скачивание")
}

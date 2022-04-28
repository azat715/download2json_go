package main

import (
	"download2json/internal/core"
	l "download2json/internal/logger"
)

const FOLDER string = "albums"
const ALBUM_URL string = "https://jsonplaceholder.typicode.com/albums/"
const PHOTOS_URL string = "https://jsonplaceholder.typicode.com/photos/"

func main() {
	l.GeneralLogger.Println("Началось скачивание")
	core.Create_folder(FOLDER)
	core.Core(ALBUM_URL, PHOTOS_URL, FOLDER)
	l.GeneralLogger.Println("Закончилось скачивание")
}

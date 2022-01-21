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
	utils.Create_folder(FOLDER)
	core.DownloadAll(ALBUM_URL, PHOTOS_URL, FOLDER)
	l.GeneralLogger.Println("Закончилось скачивание")
}

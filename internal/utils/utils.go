package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

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

func Create_folder(folder string) {
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folder, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

package core

import (
	"fmt"
	"io"
	"net/http"
)

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Http Error status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}

func Download(url string, resultCh chan []byte, errorCh chan error) {
	r, err := get(url)
	if err != nil {
		errorCh <- NewErrorWrapper("download_url", err, fmt.Sprintf("failed download url: %s", url))
	}
	resultCh <- r
}

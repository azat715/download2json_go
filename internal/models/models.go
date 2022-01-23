package models

import (
	"encoding/json"
)

type Album struct {
	Userid int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
}

func (c Album) Serialize(data []byte) ([]Album, error) {
	albums := []Album{}
	err := json.Unmarshal(data, &albums)
	if err != nil {
		return nil, err
	}
	return albums, nil
}

type Photo struct {
	Albumid      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	Thumbnailurl string `json:"thumbnailUrl"`
}

func (c Photo) Serialize(data []byte) ([]Photo, error) {
	photos := []Photo{}
	err := json.Unmarshal(data, &photos)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

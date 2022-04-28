package models

import (
	"encoding/json"
)

const PhotoExt string = "png"

type Album struct {
	Userid int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
}

type Albums []Album

func (c Albums) Parse(data []byte) (Albums, error) {
	albums := Albums{}
	return albums, unmarshal(data, &albums)
}

func (c Albums) AsDict() map[int]Album {
	albumsDict := make(map[int]Album)
	for _, i := range c {
		albumsDict[i.Id] = i
	}
	return albumsDict
}

type Photo struct {
	Albumid      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	Thumbnailurl string `json:"thumbnailUrl"`
}

type Photos []Photo

func (c Photos) Parse(data []byte) (Photos, error) {
	photos := Photos{}
	return photos, unmarshal(data, &photos)
}

func unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	return nil
}

type FilePath struct {
	Path   string
	Folder string
}

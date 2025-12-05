package main

import (
	"time"
)
type Image struct {
    ID         int       `json:"id"`
    Filename   string    `json:"filename"`
    UploadTime time.Time `json:"upload_time"`
    Memo       string    `json:"memo"`
    UserID     int       `json:"user_id"`
    WallID     int       `json:"wall_id"`
    AlbumDate  time.Time `json:"album_date"`
    URL        string    `json:"url"`  
}


type createAlbumRequest struct{
	Name string `json:"name"`
	AlbumDate string `json:"album_date"`
	UserID int `json:"user_id"`
}

type Album struct{
	ID int `json:"id"`
	Name string `json:"name"`
	AlbumDate time.Time `json:"album_date"`
}
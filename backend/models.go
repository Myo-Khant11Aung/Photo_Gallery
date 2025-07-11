package main

import "time"

type Image struct{
	ID         int       `json:"id"`
	Filename   string    `json:"filename"`
	UploadTime time.Time    `json:"upload_time"`
	Memo       string    `json:"memo"`
	UserID     int       `json:"user_id"`
	WallID     int       `json:"wall_id"`
}
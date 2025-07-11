package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)




func imageHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*") 
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") 
	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()
	var images []Image
	rows, err := pool.Query(ctx, "SELECT id, filename, upload_time, memo, user_id, wall_id FROM images")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var img Image
		err := rows.Scan(
			&img.ID,
			&img.Filename,
			&img.UploadTime,
			&img.Memo,
			&img.UserID,
			&img.WallID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	images = append(images, img)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
	
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()

	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Cannot Parse Form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		dst, err := os.Create("uploads/" + header.Filename)
		if err != nil {
			http.Error(w, "Error Creating File", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error Copying File", http.StatusInternalServerError)
			return
		}
		memo := ""
		userID := 1
		wallID := 1

		sql := `
			INSERT INTO images (filename, upload_time, memo, user_id, wall_id)
			VALUES ($1, NOW(), $2, $3, $4)
		`

		_, err = pool.Exec(ctx, sql, header.Filename, memo, userID, wallID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Upload Successful!",
		})
	} else {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}


func main(){
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/images", imageHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("uploads"))))
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
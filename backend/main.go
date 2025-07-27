package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)




func imageHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*") 
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") 
	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()
	var images []Image
	rows, err := pool.Query(ctx, `SELECT id, filename, upload_time, memo, user_id, wall_id, album_date FROM images ORDER BY album_date DESC, upload_time ASC, id ASC`)

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
			&img.AlbumDate,
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

func updateMemoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT") 
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT method is allowd", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 5 || parts[4] != "memo" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	idStr := parts[3]
	id , err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Memo string `json:"memo"`
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid Json Body", http.StatusBadRequest)
		return
	}
	
	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()

	_, err = pool.Exec(ctx, "UPDATE images SET memo = $1 WHERE id = $2", body.Memo, id)
	if err != nil {
		http.Error(w, "Failed to update memo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})

}

func cleanupMissingFiles() {
	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()

	// Step 1: Get all files from uploads/ folder
	filesOnDisk, err := os.ReadDir("uploads")
	if err != nil {
		log.Fatalf("Failed to read uploads folder: %v", err)
	}

	// Build a map of existing filenames
	fileMap := make(map[string]bool)
	for _, file := range filesOnDisk {
		fileMap[file.Name()] = true
	}

	// Step 2: Get all DB image records
	rows, err := pool.Query(ctx, "SELECT id, filename FROM images")
	if err != nil {
		log.Fatalf("Failed to query images: %v", err)
	}
	defer rows.Close()

	var deleted int

	for rows.Next() {
		var id int
		var filename string
		rows.Scan(&id, &filename)

		if !fileMap[filename] {
			// File is missing – delete the DB record
			_, err := pool.Exec(ctx, "DELETE FROM images WHERE id = $1", id)
			if err != nil {
				log.Printf("Failed to delete image ID %d: %v", id, err)
			} else {
				log.Printf("Deleted image ID %d (missing file %s)", id, filename)
				deleted++
			}
		}
	}

	log.Printf("Cleanup complete. %d entries deleted.", deleted)
}


func main(){
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/images", imageHandler)
	http.HandleFunc("/api/photo/", updateMemoHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("uploads"))))
	// cleanupMissingFiles()
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
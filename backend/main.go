package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil 
}

var jwtSecret = []byte(os.Getenv("HS256_SECRET"))


type contextKey string

const userContextKey = contextKey("userID")
const wallContextKey = contextKey("wallID")

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
			w.WriteHeader(http.StatusOK)
			return
		}

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
            return
        }
		wallIDFloat, ok := claims["wall_id"].(float64)
		if !ok {
			http.Error(w, "Invalid wall ID in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, int(userIDFloat))
		ctx = context.WithValue(ctx, wallContextKey, int(wallIDFloat))
		next.ServeHTTP(w, r.WithContext(ctx))
    }
}


func generateJWT(userID int, wallID int) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"wall_id": wallID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func registerHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT") 
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var Input struct {
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password_hash"`
	}

	err := json.NewDecoder(r.Body).Decode(&Input)
	if err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
	}

	hashedpassword, err := HashPassword(Input.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
	}

	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()
	_, err = pool.Exec(ctx, "INSERT INTO users (username, email, password_hash, wall_id) VALUES ($1, $2, $3, NULL)", Input.Username, Input.Email, hashedpassword)
	if err != nil {
        http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
        return
    }
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})

}

func loginHandler(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT") 
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post Allowed", http.StatusMethodNotAllowed)
		return
	}

	var Input struct{
		Email string `json:"email"`
		Password string `json:"password_hash"`
	}
	err := json.NewDecoder(r.Body).Decode(&Input)
	if err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()

	var storedPassword string
	var userID int
	var wallID sql.NullInt16

	err = pool.QueryRow(ctx, "SELECT id, password_hash, wall_id FROM users WHERE email = $1", Input.Email).Scan(&userID, &storedPassword, &wallID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if !CheckPasswordHash(Input.Password, storedPassword) {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}
if !wallID.Valid {
	http.Error(w, "User has no wall assigned", http.StatusForbidden)
	return
}
	token, err := generateJWT(userID, int(wallID.Int16))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}


	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "login successful",
        "user_id": userID,
		"wall_id": wallID,
        "token": token,
	})

}

func imageHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

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
	// --- CORS ---
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	pool := connectDatabase()
	defer pool.Close()

	// ~10 MB total memory buffer, rest to temp files
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Cannot Parse Form", http.StatusBadRequest)
		return
	}

	// Expect multiple files under the same field name "image"
	files := r.MultipartForm.File["image"]
	if len(files) == 0 {
		http.Error(w, "No files provided (field: image)", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(userContextKey).(int)
	wallID := r.Context().Value(wallContextKey).(int)

	type Saved struct {
		Filename string `json:"filename"`
	}
	var saved []Saved

	// Ensure uploads dir exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		http.Error(w, "Error ensuring uploads dir", http.StatusInternalServerError)
		return
	}

	for _, fh := range files {
		src, err := fh.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusBadRequest)
			return
		}
		defer src.Close()

		name := filepath.Base(fh.Filename)
		tsName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), name)

		dst, err := os.Create(filepath.Join("uploads", tsName))
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, src); err != nil {
			dst.Close()
			http.Error(w, "Error copying file", http.StatusInternalServerError)
			return
		}
		dst.Close()

		const sql = `
			INSERT INTO images (filename, upload_time, memo, user_id, wall_id)
			VALUES ($1, NOW(), $2, $3, $4)
		`
		if _, err := pool.Exec(ctx, sql, tsName, "", userID, wallID); err != nil {
			http.Error(w, "DB insert failed", http.StatusInternalServerError)
			return
		}

		saved = append(saved, Saved{Filename: tsName})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":  "Upload successful",
		"count":    len(saved),
		"files":    saved,
	})
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
	http.HandleFunc("/api/upload",jwtMiddleware(uploadHandler))
	http.HandleFunc("/api/images", jwtMiddleware(imageHandler))
	http.HandleFunc("/api/photo/", jwtMiddleware(updateMemoHandler))
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("uploads"))))
	// cleanupMissingFiles()
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
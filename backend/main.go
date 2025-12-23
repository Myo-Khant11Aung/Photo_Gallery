package main

import (
	"bytes"
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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var jwtSecret = []byte(os.Getenv("HS256_SECRET"))

type contextKey string

const userContextKey = contextKey("userID")
const wallContextKey = contextKey("wallID")

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		// if r.Method == http.MethodOptions {
		// 	w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		// 	w.WriteHeader(http.StatusOK)
		// 	return
		// }
		if r.Method == http.MethodOptions {
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
	})
}

func generateJWT(userID int, wallID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"wall_id": wallID,
		"exp":     time.Now().Add(time.Hour * 6).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var Input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post Allowed", http.StatusMethodNotAllowed)
		return
	}

	var Input struct {
		Email    string `json:"email"`
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
		"status":  "login successful",
		"user_id": userID,
		"wall_id": wallID,
		"token":   token,
	})

}

func (c *R2Client) PresignURL(ctx context.Context, key string, expires time.Duration, w http.ResponseWriter) (string, error) {
	presigner := s3.NewPresignClient(c.s3)
	out, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		log.Printf("Presign error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Presign error")
	}

	return out.URL, nil
}

func imageHandler(pool *pgxpool.Pool, r2 *R2Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("imageHandler called")
		ctx := context.Background()
		wallID := r.Context().Value(wallContextKey).(int)
		var images []Image
		rows, err := pool.Query(ctx, `SELECT id, filename, upload_time, memo, user_id, wall_id, album_date
            FROM images
            WHERE wall_id = $1
            ORDER BY album_date DESC, upload_time ASC, id ASC`, wallID)

		if err != nil {
			http.Error(w, "DB query failed: "+err.Error(), 500)
			return
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
			url, err := r2.PresignURL(ctx, img.Filename, 10*time.Minute, w)
			if err != nil {
				log.Printf("Presign error: %v", err)
				writeJSONError(w, http.StatusInternalServerError, "Presign error")
			}
			img.URL = url
			images = append(images, img)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)

	}
}

func (c *R2Client) UploadToR2(ctx context.Context, key string, contentType string, data []byte) error {
	_, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	return err
}

func uploadHandler(pool *pgxpool.Pool, r2 *R2Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := context.Background()

		// ~10 MB total memory buffer, rest to temp files
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Cannot Parse Form", http.StatusBadRequest)
			return
		}

		albumIDStr := r.FormValue("album_id")
		if albumIDStr == "" {
			http.Error(w, "album_id is required", http.StatusBadRequest)
			return
		}

		albumID, err := strconv.Atoi(albumIDStr)
		if err != nil {
			http.Error(w, "Invalid album_id", http.StatusBadRequest)
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

		var exists bool
		err = pool.QueryRow(
			ctx,
			`SELECT EXISTS (
				SELECT 1 FROM albums WHERE id = $1 AND wall_id = $2
			)`,
			albumID,
			wallID,
		).Scan(&exists)

		if err != nil || !exists {
			http.Error(w, "Album not found or unauthorized", http.StatusUnauthorized)
			return
		}

		type Saved struct {
			Filename string `json:"filename"`
		}
		var saved []Saved

		for _, fh := range files {
			src, err := fh.Open()
			if err != nil {
				http.Error(w, "Error opening file", http.StatusBadRequest)
				return
			}
			data, err := io.ReadAll(src)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusBadRequest)
				return
			}
			src.Close()

			ct := http.DetectContentType(data)
			ts := time.Now().UnixNano()
			base := filepath.Base(fh.Filename)
			key := fmt.Sprintf("walls/%d/%d_%s", wallID, ts, base)

			var payload []byte
			var contentType string
			payload = data
			contentType = ct
			// }
			if err := r2.UploadToR2(ctx, key, contentType, payload); err != nil {
				http.Error(w, "R2 upload failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			// Save only the key in DB (not a URL)
			const sql = `
			INSERT INTO images (filename, upload_time, memo, user_id, wall_id, album_id)
VALUES ($1, NOW(), $2, $3, $4, $5)

		`
			if _, err := pool.Exec(ctx, sql, key, "", userID, wallID, albumID); err != nil {
				http.Error(w, "DB insert failed", http.StatusInternalServerError)
				return
			}

			saved = append(saved, Saved{Filename: key})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Upload successful",
			"count":   len(saved),
			"files":   saved,
		})
	}
}

func updateMemoHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 5 || parts[4] != "memo" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	idStr := parts[3]
	id, err := strconv.Atoi(idStr)
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

func createAlbumHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		// userID := r.Context().Value(userContextKey).(int)
		wallID := r.Context().Value(wallContextKey).(int)

		var req createAlbumRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			http.Error(w, "Album name required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		var albumID int
		var albumDate time.Time

		err := pool.QueryRow(ctx, `INSERT INTO albums (name, wall_id) VALUES ($1, $2) RETURNING id, album_date`, req.Name, wallID).Scan(&albumID, &albumDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":         albumID,
			"name":       req.Name,
			"album_date": albumDate,
		})
	}
}

func getAlbumsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}
		wallID := r.Context().Value(wallContextKey).(int)

		ctx := r.Context()
		var albums []Album
		rows, err := pool.Query(ctx, `SELECT id, name FROM albums WHERE wall_id = $1
            ORDER BY album_date DESC`, wallID)

		if err != nil {
			http.Error(w, "DB query failed: "+err.Error(), 500)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var album Album
			if err := rows.Scan(&album.ID, &album.Name); err != nil {
				http.Error(w, "Failed to scan album: "+err.Error(), http.StatusInternalServerError)
				return
			}
			albums = append(albums, album)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Row iteration error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(albums)
	}
}
func meHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userContextKey).(int)
	wallID := r.Context().Value(wallContextKey).(int)

	json.NewEncoder(w).Encode(map[string]any{
		"user_id": userID,
		"wall_id": wallID,
		"status":  "ok",
	})
}

func deletePhotoHandler(db *pgxpool.Pool, r2 *R2Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract ID from URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		idStr := parts[4]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		userWallID := ctx.Value(wallContextKey).(int)

		// STEP 1 — Verify image exists & belongs to that wall
		var filename string
		var wallID int

		err = db.QueryRow(ctx,
			`SELECT filename, wall_id FROM images WHERE id = $1`,
			id,
		).Scan(&filename, &wallID)

		if err == pgx.ErrNoRows {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "DB query failed", http.StatusInternalServerError)
			return
		}

		// Security check — only delete images in your own wall
		if wallID != userWallID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// STEP 2 — Delete from R2
		err = r2.DeleteObject(ctx, filename)
		if err != nil {
			http.Error(w, "Failed to delete from storage: "+err.Error(), 500)
			return
		}

		// STEP 3 — Delete from DB
		_, err = db.Exec(ctx,
			`DELETE FROM images WHERE id = $1`,
			id,
		)
		if err != nil {
			http.Error(w, "Failed to delete from database", http.StatusInternalServerError)
			return
		}

		// STEP 4 — Return success JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success":    true,
			"deleted_id": id,
		})
	}
}

func getAlbumImagesHandler(pool *pgxpool.Pool, r2 *R2Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}

		// URL: /api/albums/{id}/images
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		albumID, err := strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "Invalid album id", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		wallID := ctx.Value(wallContextKey).(int)

		// Security check: album must belong to wall
		var exists bool
		err = pool.QueryRow(ctx,
			`SELECT EXISTS (
				SELECT 1 FROM albums WHERE id = $1 AND wall_id = $2
			)`,
			albumID, wallID,
		).Scan(&exists)

		if err != nil || !exists {
			http.Error(w, "Album not found", http.StatusNotFound)
			return
		}

		rows, err := pool.Query(ctx, `
			SELECT id, filename, upload_time, memo, user_id, wall_id
			FROM images
			WHERE album_id = $1
			ORDER BY upload_time ASC, id ASC
		`, albumID)

		if err != nil {
			http.Error(w, "DB query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var images []Image

		for rows.Next() {
			var img Image
			if err := rows.Scan(
				&img.ID,
				&img.Filename,
				&img.UploadTime,
				&img.Memo,
				&img.UserID,
				&img.WallID,
			); err != nil {
				http.Error(w, "Scan failed", http.StatusInternalServerError)
				return
			}

			url, err := r2.PresignURL(ctx, img.Filename, 10*time.Minute, w)
			if err != nil {
				http.Error(w, "Presign failed", http.StatusInternalServerError)
				return
			}

			img.URL = url
			images = append(images, img)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(images)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from the system")
	}

	r2Client, err := NewR2Client()
	if err != nil {
		log.Fatalf("Failed to create R2 client: %v", err)
		return
	}

	db := connectDatabase()
	defer db.Close()

	mux := http.NewServeMux()

	// Protected routes (JWT + CORS)
	mux.HandleFunc("/api/upload", jwtMiddleware(uploadHandler(db, r2Client)))

	mux.HandleFunc("/api/images", jwtMiddleware(imageHandler(db, r2Client)))

	mux.HandleFunc("/api/photo/", jwtMiddleware(http.HandlerFunc(updateMemoHandler)))

	mux.HandleFunc("/api/verifyToken", jwtMiddleware(http.HandlerFunc(meHandler)))

	mux.HandleFunc("/api/create_album", jwtMiddleware(http.HandlerFunc(createAlbumHandler(db))))

	// Public routes (CORS, no JWT)
	mux.HandleFunc("/api/register", registerHandler)

	mux.HandleFunc("/api/login", loginHandler)

	mux.Handle("/api/photo/delete/", jwtMiddleware(deletePhotoHandler(db, r2Client)))

	mux.HandleFunc("/api/albums", jwtMiddleware(getAlbumsHandler(db)))

	mux.Handle("/api/albums/", jwtMiddleware(getAlbumImagesHandler(db, r2Client)))
	c := cors.Options{
		AllowedOrigins: []string{"localhost:3000", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}

	handler := cors.New(c).Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s\n", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

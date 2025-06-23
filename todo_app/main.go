package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	refreshMutex sync.Mutex
	isRefreshing bool
)

const (
	picsumURL     = "https://picsum.photos/1200"
	cacheDuration = 10 * time.Minute
	storageDir    = "storage"
	photoPath     = "storage/latest_image.jpg"
)

func fetchAndSaveNewImage() error {
	log.Println("Fetching a new image from Lorem Picsum...")

	resp, err := http.Get(picsumURL)
	if err != nil {
		return fmt.Errorf("failed to make request to picsum: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("picsum returned a non-200 status code: %s", resp.Status)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read image data from response: %w", err)
	}

	if err := os.WriteFile(photoPath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to save image to disk: %w", err)
	}

	log.Printf("Successfully fetched ans saved an image")
	return nil
}

func refreshImage() {
	refreshMutex.Lock()
	if isRefreshing {
		refreshMutex.Unlock()
		return
	}
	isRefreshing = true
	refreshMutex.Unlock()

	defer func() {
		refreshMutex.Lock()
		isRefreshing = false
		refreshMutex.Unlock()
	}()

	if err := fetchAndSaveNewImage(); err != nil {
		log.Printf("Error refreshing image file in background: %v", err)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	fileInfo, err := os.Stat(photoPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Photo not found. Fetching initial image.")
			if fetchErr := fetchAndSaveNewImage(); fetchErr != nil {
				http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
				log.Printf("Error during initial fetch: %v", fetchErr)
				return
			}
			fileInfo, err = os.Stat(photoPath)
			if err != nil {
				http.Error(w, "Failed to read newly created file", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Could not access photo storage", http.StatusInternalServerError)
			log.Printf("Error stating file: %v", err)
			return
		}
	}

	if time.Since(fileInfo.ModTime()) > cacheDuration {
		log.Println("File is stale. Serving old image and refreshing in background.")
		go refreshImage()
	} else {
		log.Println("Serving fresh image from file.")
	}

	imageData, err := os.ReadFile(photoPath)
	if err != nil {
		http.Error(w, "Could not read photo from storage", http.StatusInternalServerError)
		log.Printf("Error reading file: %v", err)
		return
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	dataURI := "data:image/jpeg;base64," + encodedImage

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var htmlBuffer bytes.Buffer
	htmlBuffer.WriteString(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Filesystem Image Viewer</title>
    <style>
        body { font-family: sans-serif; display: flex; flex-direction: column; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; margin: 0; text-align: center; }
        img { max-width: 90%; max-height: 80vh; border: 5px solid white; box-shadow: 0 4px 15px rgba(0,0,0,0.2); }
		p { color: #555; margin-top: 20px; }
    </style>
</head>
<body>
	<h1>The project App</h1>
    <img src="`)
	htmlBuffer.WriteString(dataURI)
	htmlBuffer.WriteString(`" alt="A random image from Picsum, served from the filesystem">
	<p>DevOps with Kubernetes 2025</p>
</body>
</html>
`)

	w.Write(htmlBuffer.Bytes())
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.NewServeMux()
	server.HandleFunc("/", imageHandler)

	log.Printf("Server started in port %s...", port)

	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}

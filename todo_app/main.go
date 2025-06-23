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
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; display: flex; flex-direction: column; align-items: center; background-color: #f0f2f5; margin: 0; padding: 20px; }
        .container { max-width: 800px; width: 100%; }
        img { width: 100%; height: auto; border-radius: 8px; border: 5px solid white; box-shadow: 0 4px 15px rgba(0,0,0,0.1); }
        .todo-form { display: flex; margin-top: 20px; }
        .todo-form input { flex-grow: 1; border: 1px solid #ccc; border-radius: 4px; padding: 10px; font-size: 16px; }
        .todo-form button { background-color: #007bff; color: white; border: none; padding: 10px 20px; margin-left: 10px; border-radius: 4px; cursor: pointer; font-size: 16px; }
        .todo-form button:hover { background-color: #0056b3; }
        .todo-list { background-color: white; list-style-type: none; padding: 10px 20px; margin-top: 20px; border-radius: 8px; box-shadow: 0 4px 15px rgba(0,0,0,0.05); }
        .todo-list h2 { text-align: center; color: #333; }
        .todo-list li { padding: 12px 0; border-bottom: 1px solid #eee; color: #555; }
        .todo-list li:last-child { border-bottom: none; }
    </style>
</head>
<body>
	<div class="container">
		<h1>The project App</h1>
		<img src="`)
	htmlBuffer.WriteString(dataURI)
	htmlBuffer.WriteString(`" alt="A random image from Picsum, served from the filesystem">
		<form class="todo-form" onsubmit="event.preventDefault(); alert('Submitting todos is not implemented yet!');">
			<input type="text" name="todo" placeholder="What needs to be done?" maxlength="140" required>
			<button type="submit">Add Todo</button>
		</form>

		<div class="todo-list">
			<h2>My Todos</h2>
			<ul>
				<li>Learn Go concurrency</li>
				<li>Set up a Kubernetes cluster</li>
				<li>Master HTML forms and CSS</li>
				<li>Read a book on system design</li>
			</ul>
		</div>
		<p>DevOps with Kubernetes 2025</p>
	</div>
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

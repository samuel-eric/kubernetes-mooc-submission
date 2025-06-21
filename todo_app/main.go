package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		htmlContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Web Server</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; margin: 0; }
        h1 { color: #333; }
    </style>
</head>
<body>
    <h1>Hello from my Go Server!</h1>
</body>
</html>
`
		fmt.Fprint(w, htmlContent)
	})

	log.Printf("Server started in port %s...", port)

	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatal(err)
	}
}

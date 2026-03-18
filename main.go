package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	latestPostbackData struct {
		Path        string              `json:"path"`
		QueryParams map[string][]string `json:"query_params"`
	}
	dataMutex sync.RWMutex
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		next(w, r)
	}
}

func postbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received postback request: %s %s", r.Method, r.URL.String())

	path := r.URL.Path
	queryParams := r.URL.Query()

	dataMutex.Lock()
	latestPostbackData.Path = path
	latestPostbackData.QueryParams = queryParams
	dataMutex.Unlock()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Postback received successfully!")
}

func apiViewHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(latestPostbackData); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()
	data := latestPostbackData

	tmpl, err := template.New("postback").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AppMetrica Postback Viewer</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <div class="container">
        <h1>AppMetrica Postback Data</h1>
        
        <h2>Request Path</h2>
        <p><strong>{{.Path}}</strong></p>
        
        <h2>Query Parameters</h2>
        <ul>
            {{if .QueryParams}}
                {{range $key, $values := .QueryParams}}
                    <li>
                        <strong>{{$key}}:</strong> 
                        <span class="param-value">
                            {{range $values}}
                                {{.}} 
                            {{end}}
                        </span>
                    </li>
                {{end}}
            {{else}}
                <li>No query parameters received.</li>
            {{end}}
        </ul>
    </div>
</body>
</html>
`)
	if err != nil {
		log.Printf("Error parsing HTML template: %v", err)
		http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing HTML template: %v", err)
		http.Error(w, "Error executing HTML template", http.StatusInternalServerError)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT environment variable not set, defaulting to 8080")
	}

	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/view", corsMiddleware(viewHandler))
	http.HandleFunc("/api/view", corsMiddleware(apiViewHandler))
	http.HandleFunc("/postback", corsMiddleware(postbackHandler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			http.ServeFile(w, r, filepath.Join("frontend", "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})

	log.Printf("Server starting on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

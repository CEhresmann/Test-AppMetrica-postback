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
	"time"
)

type PostbackData struct {
	Path        string              `json:"path"`
	QueryParams map[string][]string `json:"query_params"`
	ReceivedAt  string              `json:"received_at"`
}

var (
	postbackHistory []PostbackData
	dataMutex       sync.RWMutex
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func postbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received postback request: %s %s", r.Method, r.URL.String())

	path := r.URL.Path
	queryParams := r.URL.Query()
	receivedAt := time.Now().Format("2006-01-02 15:04:05")

	newData := PostbackData{
		Path:        path,
		QueryParams: queryParams,
		ReceivedAt:  receivedAt,
	}

	dataMutex.Lock()
	postbackHistory = append([]PostbackData{newData}, postbackHistory...)
	dataMutex.Unlock()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Postback received successfully!")
}

func apiViewHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(postbackHistory); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	data := postbackHistory
	dataMutex.RUnlock()

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
        <h1>AppMetrica Postback History</h1>
        
        {{if .}}
            {{range .}}
                <div class="postback-item" style="border-bottom: 1px solid #ccc; margin-bottom: 20px; padding-bottom: 10px;">
                    <h2>Request Path</h2>
                    <p><strong>{{.Path}}</strong></p>

                    <h2>Received At</h2>
                    <p><strong>{{.ReceivedAt}}</strong></p>
                    
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
            {{end}}
        {{else}}
            <p>No postbacks received yet.</p>
        {{end}}
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

	fs := http.FileServer(http.Dir("docs"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/view", corsMiddleware(viewHandler))
	http.HandleFunc("/api/view", corsMiddleware(apiViewHandler))
	http.HandleFunc("/postback", corsMiddleware(postbackHandler))

	http.HandleFunc("/", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if len(r.URL.Query()) > 0 {
				postbackHandler(w, r)
				return
			}
			http.ServeFile(w, r, filepath.Join("docs", "index.html"))
			return
		}

		if stat, err := os.Stat(filepath.Join("docs", r.URL.Path)); err == nil && !stat.IsDir() {
			http.ServeFile(w, r, filepath.Join("docs", r.URL.Path))
			return
		}

		postbackHandler(w, r)
	}))

	log.Printf("Server starting on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

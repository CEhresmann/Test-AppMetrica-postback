package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var latestPostbackData struct {
	Path        string
	QueryParams map[string][]string
}

func postbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received postback request: %s %s", r.Method, r.URL.String())

	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		log.Printf("Error parsing URL: %v", err)
		http.Error(w, "Error parsing URL", http.StatusInternalServerError)
		return
	}

	path := parsedURL.Path
	queryParams := parsedURL.Query()

	latestPostbackData = struct {
		Path        string
		QueryParams map[string][]string
	}{
		Path:        path,
		QueryParams: queryParams,
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Postback received successfully!")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Path        string
		QueryParams map[string][]string
	}{
		Path:        latestPostbackData.Path,
		QueryParams: latestPostbackData.QueryParams,
	}

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

	http.HandleFunc("/view", viewHandler)

	http.HandleFunc("/postback", postbackHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("frontend", "index.html"))
			return
		}
		http.NotFound(w, r)
	})

	log.Printf("Server starting on port %s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

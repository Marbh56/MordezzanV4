package middleware

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// CustomFileServer wraps http.FileServer to ensure correct MIME types are set
func CustomFileServer(root http.FileSystem) http.Handler {
	fileServer := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get file extension from path
		ext := strings.ToLower(filepath.Ext(r.URL.Path))

		// Set appropriate Content-Type header based on file extension
		switch ext {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		default:
			// For other files, let the system determine the type
			contentType := mime.TypeByExtension(ext)
			if contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
		}

		// Serve the file
		fileServer.ServeHTTP(w, r)
	})
}

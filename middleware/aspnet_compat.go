package middleware

import (
	"net/http"
	"path/filepath"
)

// AspNetCompat implementa um middleware para tratar extensões típicas ASP.NET
// (como .axd, .aspx) roteando para index.html ou tratando static files
func AspNetCompat(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".axd", ".ashx":
			// Simplesmente trate como arquivo estático
			http.ServeFile(w, r, "wwwroot"+r.URL.Path)
			return
		case ".aspx":
			// Roteie tudo para index.html (Single Page App)
			http.ServeFile(w, r, "wwwroot/index.html")
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}
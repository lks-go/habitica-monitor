package http

import "net/http"

// CORS — простое middleware, разрешающее запросы от Web UI.
// Origin берётся из allowedOrigin ("*" по умолчанию). В dev запросы и так
// проксируются Vite, но для prod-сборки фронта на другом origin это нужно.
func CORS(allowedOrigin string, next http.Handler) http.Handler {
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

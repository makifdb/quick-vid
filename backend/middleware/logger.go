package middleware

import (
	"net/http"
	"time"

	"github.com/makifdb/quick-vid/utils"
)

// LoggerMiddleware logs incoming requests and response times
func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Set CORS headers
		utils.SetCORSHeaders(w)

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		defer func() {
			// Log the request
			utils.LogRequest(r, time.Since(start))
		}()

		// Call the next handler in the chain
		next(w, r)
	}
}

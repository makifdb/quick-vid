package utils

import (
	"log"
	"net/http"
	"time"
)

// LogInfo logs informational messages
func LogInfo(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// LogError logs error messages
func LogError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// LogResponseTime logs the response time for an HTTP request
func LogRequest(r *http.Request, duration time.Duration) {
	log.Printf(
		"%s %s %s took %v",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		duration,
	)
}

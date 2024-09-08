package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type HandlerFuncWithError = func(http.ResponseWriter, *http.Request) error

func loggingMiddleware(next HandlerFuncWithError) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logText := fmt.Sprintf("%s %s", req.Method, req.URL)
		slog.Info("Request received " + logText)

		start := time.Now()
		err := next(res, req)
		elapsedTime := time.Since(start)

		if err != nil {
			slog.Error("Request failed "+logText, "error", err, "elapsedTime", elapsedTime)
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			slog.Info("Request finished "+logText, "elapsedTime", elapsedTime)
		}
	}
}

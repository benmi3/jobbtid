/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"fmt"
	"jobbtid/pkg/config"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The current time is %s", time.Now())
}

func Serve(cfg *config.Config) {
	fmt.Println("Starting server")
	host, err := cfg.Host("server")
	if err != nil {
		panic(fmt.Errorf("fatal error getting hostname from config: %w", err))
	}
	// Server start
	http.HandleFunc("/", greet)
	http.Handle("/start", authMiddleware(http.HandlerFunc(start)))
	http.Handle("/stop", authMiddleware(http.HandlerFunc(stop)))
	http.Handle("/check", authMiddleware(http.HandlerFunc(check)))
	http.ListenAndServe(host, nil)
}

/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"fmt"
	"net/http"
	"time"
)

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

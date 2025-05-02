/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package db

import (
	"testing"
)

func TestGetDefaultConfig(t *testing.T) {
	t.Run("Can Create Default Config", func(t *testing.T) {
		_, err := setupDbCon()
		if err != nil {
			t.Error("Wanted default is true, got false")
		}
	})
}

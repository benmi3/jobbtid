/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package config

import (
	"testing"
)

func TestGetDefaultConfig(t *testing.T) {
	t.Run("Can Create Default Config", func(t *testing.T) {
		res := createDefaultConfig()
		if res.ServerConf.Default != true {
			t.Error("Wanted default is true, got false")
		}
	})
}

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

func TestGetDefaultConfigHost(t *testing.T) {
	t.Run("Can get host when configName is set to server", func(t *testing.T) {
		res := createDefaultConfig()
		if res.ServerConf.Default != true {
			t.Error("Wanted default is true, got false")
		}
		host_test, err := res.Host("server")
		if err != nil {
			t.Error("host_test returned an error")
		}
		default_host := ":8080"
		if host_test != default_host {
			t.Errorf("host_test name returned %q, wanted %q", host_test, default_host)
		}
	})
}

func TestGetDefaultConfigHostError(t *testing.T) {
	t.Run("Host will error if unsupported configName is set", func(t *testing.T) {
		res := createDefaultConfig()
		if res.ServerConf.Default != true {
			t.Error("Wanted default is true, got false")
		}
		host_test, err := res.Host("")
		if err == nil {
			t.Error("host_test returned nil instead of error")
		}
		default_host := ""
		if host_test != default_host {
			t.Errorf("host_test name returned %q, wanted %q", host_test, default_host)
		}
	})
}

func TestReadConfigFile(t *testing.T) {
	t.Run("Can Read from config", func(t *testing.T) {
		res, err := ReadConfigFile("./config_test.toml")
		if res.ServerConf.Default != false {
			t.Error("Wanted default is false, got true")
		}
		if err != nil {
			t.Error("Error happened when reading empty config file. should not happen")
		}
		host_test, err := res.Host("server")
		if err != nil {
			t.Error("host_test returned error instead of nil")
		}
		want_host := "localhost:8081"
		if host_test != want_host {
			t.Errorf("host_test name returned %q, wanted %q", host_test, want_host)
		}
	})
}

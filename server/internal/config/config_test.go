package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/jsanca/abacus/server/internal/config"
)

func TestLoad_defaults(t *testing.T) {
	serverConfig, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error loading defaults, got: %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"ListenAddress", serverConfig.ListenAddress, ":8080"},
		{"ShutdownTimeout", serverConfig.ShutdownTimeout, 10 * time.Second},
		{"ReadTimeout", serverConfig.ReadTimeout, 5 * time.Second},
		{"WriteTimeout", serverConfig.WriteTimeout, 10 * time.Second},
		{"IdleTimeout", serverConfig.IdleTimeout, 120 * time.Second},
		{"RequestTimeout", serverConfig.RequestTimeout, 30 * time.Second},
		{"AllowedOrigin", serverConfig.AllowedOrigin, "http://localhost:5173"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.got != testCase.want {
				t.Errorf("expected %v, got %v", testCase.want, testCase.got)
			}
		})
	}
}

func TestLoad_envOverrides(t *testing.T) {
	t.Setenv("LISTEN_ADDRESS", ":9090")
	t.Setenv("SHUTDOWN_TIMEOUT", "5s")
	t.Setenv("ALLOWED_ORIGIN", "http://localhost:3000")

	serverConfig, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if serverConfig.ListenAddress != ":9090" {
		t.Errorf("expected ListenAddress %q, got %q", ":9090", serverConfig.ListenAddress)
	}
	if serverConfig.ShutdownTimeout != 5*time.Second {
		t.Errorf("expected ShutdownTimeout %v, got %v", 5*time.Second, serverConfig.ShutdownTimeout)
	}
	if serverConfig.AllowedOrigin != "http://localhost:3000" {
		t.Errorf("expected AllowedOrigin %q, got %q", "http://localhost:3000", serverConfig.AllowedOrigin)
	}
}

func TestLoad_invalidDuration(t *testing.T) {
	tests := []struct {
		name        string
		envKey      string
		envValue    string
		wantErrFrag string
	}{
		{
			name:        "invalid SHUTDOWN_TIMEOUT",
			envKey:      "SHUTDOWN_TIMEOUT",
			envValue:    "notaduration",
			wantErrFrag: "invalid SHUTDOWN_TIMEOUT",
		},
		{
			name:        "invalid READ_TIMEOUT",
			envKey:      "READ_TIMEOUT",
			envValue:    "notaduration",
			wantErrFrag: "invalid READ_TIMEOUT",
		},
		{
			name:        "invalid WRITE_TIMEOUT",
			envKey:      "WRITE_TIMEOUT",
			envValue:    "notaduration",
			wantErrFrag: "invalid WRITE_TIMEOUT",
		},
		{
			name:        "invalid IDLE_TIMEOUT",
			envKey:      "IDLE_TIMEOUT",
			envValue:    "notaduration",
			wantErrFrag: "invalid IDLE_TIMEOUT",
		},
		{
			name:        "invalid REQUEST_TIMEOUT",
			envKey:      "REQUEST_TIMEOUT",
			envValue:    "notaduration",
			wantErrFrag: "invalid REQUEST_TIMEOUT",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv(testCase.envKey, testCase.envValue)
			_, err := config.Load()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), testCase.wantErrFrag) {
				t.Errorf("expected error containing %q, got %q", testCase.wantErrFrag, err.Error())
			}
		})
	}
}

// Package config provides environment-driven configuration for the Abacus backend.
package config

import (
	"fmt"
	"os"
	"time"
)

// ServerConfig holds all configuration values for the Abacus HTTP server.
type ServerConfig struct {
	// ListenAddress is the TCP address the HTTP server will listen on, e.g. ":8080".
	ListenAddress string
	// ShutdownTimeout is the maximum duration allowed for graceful shutdown to complete.
	ShutdownTimeout time.Duration
	// ReadTimeout is the maximum duration for reading the entire incoming request.
	ReadTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration
	// IdleTimeout is the maximum duration to wait for the next request on a keep-alive connection.
	IdleTimeout time.Duration
	// RequestTimeout is the maximum duration for a single HTTP request to complete.
	RequestTimeout time.Duration
	// AllowedOrigin is the CORS origin permitted for local frontend development.
	AllowedOrigin string
}

// Load reads configuration from environment variables and returns a validated ServerConfig.
// It fails fast with contextual errors if any value is invalid or unparseable.
func Load() (ServerConfig, error) {
	listenAddress := envStringOrDefault("LISTEN_ADDRESS", ":8080")
	allowedOrigin := envStringOrDefault("ALLOWED_ORIGIN", "http://localhost:5173")

	shutdownTimeout, err := envDurationOrDefault("SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}

	readTimeout, err := envDurationOrDefault("READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := envDurationOrDefault("WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid WRITE_TIMEOUT: %w", err)
	}

	idleTimeout, err := envDurationOrDefault("IDLE_TIMEOUT", 120*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid IDLE_TIMEOUT: %w", err)
	}

	requestTimeout, err := envDurationOrDefault("REQUEST_TIMEOUT", 30*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid REQUEST_TIMEOUT: %w", err)
	}

	serverConfig := ServerConfig{
		ListenAddress:   listenAddress,
		ShutdownTimeout: shutdownTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		IdleTimeout:     idleTimeout,
		RequestTimeout:  requestTimeout,
		AllowedOrigin:   allowedOrigin,
	}

	if err := serverConfig.validate(); err != nil {
		return ServerConfig{}, err
	}

	return serverConfig, nil
}

func (serverConfig ServerConfig) validate() error {
	if serverConfig.ListenAddress == "" {
		return fmt.Errorf("LISTEN_ADDRESS must not be empty")
	}
	if serverConfig.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be positive, got %v", serverConfig.ShutdownTimeout)
	}
	if serverConfig.ReadTimeout <= 0 {
		return fmt.Errorf("READ_TIMEOUT must be positive, got %v", serverConfig.ReadTimeout)
	}
	if serverConfig.WriteTimeout <= 0 {
		return fmt.Errorf("WRITE_TIMEOUT must be positive, got %v", serverConfig.WriteTimeout)
	}
	if serverConfig.IdleTimeout <= 0 {
		return fmt.Errorf("IDLE_TIMEOUT must be positive, got %v", serverConfig.IdleTimeout)
	}
	if serverConfig.RequestTimeout <= 0 {
		return fmt.Errorf("REQUEST_TIMEOUT must be positive, got %v", serverConfig.RequestTimeout)
	}
	return nil
}

func envStringOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func envDurationOrDefault(key string, defaultValue time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue, nil
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("could not parse %q as a duration: %w", raw, err)
	}
	return parsed, nil
}

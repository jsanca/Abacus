// Package observability defines the telemetry extension point for the Abacus backend.
//
// Production telemetry is deferred: Abacus is a self-contained exercise and
// introducing a real metrics or tracing backend would add external dependencies
// and operational complexity disproportionate to the project scope. The boundary
// is retained so a production implementation can be wired in without restructuring
// the codebase. See ADR 0001 for the accepted decision.
package observability

import "time"

// Observer is the telemetry extension point for the Abacus HTTP server.
// The zero value is a ready-to-use no-op implementation.
type Observer struct{}

// OnServerStart is called when the HTTP server begins accepting connections.
func (observer *Observer) OnServerStart(address string) {}

// OnServerStop is called when graceful shutdown completes.
func (observer *Observer) OnServerStop() {}

// OnRequestHandled is called after each HTTP request completes.
func (observer *Observer) OnRequestHandled(method, path string, statusCode int, duration time.Duration) {
}

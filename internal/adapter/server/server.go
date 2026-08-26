// Package server isolates HTTP server start and shutdown operations from callers.
package server

import (
	"context"
	"net/http"
)

// Operations performs the server effects used by a caller.
type Operations interface {
	ListenAndServe(*http.Server) error
	Shutdown(*http.Server, context.Context) error
}

// Dependencies contains the server effects used by a caller. Its zero value is
// a safe no-op; production callers must select System explicitly.
type Dependencies struct {
	Operations Operations
}

// Selector returns the exact server selected by a caller.
type Selector func() *http.Server

// ListenAndServe starts the exact selected server through the configured dependency.
func ListenAndServe(dependencies Dependencies, selectServer Selector) error {
	if dependencies.Operations == nil || selectServer == nil {
		return nil
	}
	return dependencies.Operations.ListenAndServe(selectServer())
}

// Shutdown stops the exact selected server with the exact supplied context.
func Shutdown(dependencies Dependencies, selectServer Selector, ctx context.Context) error {
	if dependencies.Operations == nil || selectServer == nil {
		return nil
	}
	return dependencies.Operations.Shutdown(selectServer(), ctx)
}

// System returns the production dependency that operates on supplied HTTP servers.
func System() Dependencies {
	return Dependencies{Operations: systemServer{}}
}

type systemServer struct{}

func (systemServer) ListenAndServe(supplied *http.Server) error {
	return supplied.ListenAndServe()
}

func (systemServer) Shutdown(supplied *http.Server, ctx context.Context) error {
	return supplied.Shutdown(ctx)
}

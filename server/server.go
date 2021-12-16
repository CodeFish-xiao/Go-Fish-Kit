package server

import "context"

// Server is transport server.
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

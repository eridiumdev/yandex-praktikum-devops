package http

import "context"

// These are the interfaces required for handling monitoring requests

// Pingable can be pinged to check live status, e.g. for database connection
type Pingable interface {
	Ping(ctx context.Context) bool
}

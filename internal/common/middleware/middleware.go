package middleware

import (
	"net/http"
)

// BasicSet logs requests and responses main data (method, status, etc.)
// as well as adds requestId to context
var BasicSet = []func(next http.Handler) http.Handler{
	BaseLoggingMiddleware,
	AddRequestID,
	LogRequests,
	LogResponses,
}

// ExtendedSet does the same as BasicSet
// but also logs request bodies and response bodies
// as well as adds support for compressing responses and decompressing requests
var ExtendedSet = []func(next http.Handler) http.Handler{
	BaseLoggingMiddleware,
	AddRequestID,
	// LogRequestsWithBody,
	CompressResponses, // should be before LogResponsesWithBody (for human-readable logs)
	LogResponsesWithBody,
	DecompressRequests, // should be after LogResponsesWithBody
}

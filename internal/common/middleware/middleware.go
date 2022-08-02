package middleware

import (
	"net/http"
)

// BasicSet logs requests and responses main data (method, status, etc.)
// as well as adds requestId to context
// and compresses responses / decompresses requests if needed
var BasicSet = []func(next http.Handler) http.Handler{
	BaseLoggingMiddleware,
	AddRequestID,
	LogRequests,
	LogResponses,
	CompressResponses,
	DecompressRequests,
}

// ExtendedSet does the same as BasicSet
// but also logs request bodies and response bodies
var ExtendedSet = []func(next http.Handler) http.Handler{
	BaseLoggingMiddleware,
	AddRequestID,
	LogRequestsWithBody,
	CompressResponses,
	LogResponsesWithBody,
	DecompressRequests,
}

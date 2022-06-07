package types

import "net/http"

// Route datatype
type APIRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - collection of all routes
type APIRoutes []APIRoute

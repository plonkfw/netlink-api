package routing

import (
	"net/http"

	"github.com/plonkfw/netlink-api/api/v1/addr"
	"github.com/plonkfw/netlink-api/api/v1/link"
)

// APIRoute - Route datatype
type APIRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// APIRoutes - Routes - collection of all routes
type APIRoutes []APIRoute

var routes = APIRoutes{
	APIRoute{
		Name:        "api/v1/addr/Add.go",
		Method:      "POST",
		Pattern:     "/v1/addr/add",
		HandlerFunc: addr.Add,
	},
	APIRoute{
		Name:        "api/v1/addr/List.go",
		Method:      "GET",
		Pattern:     "/v1/addr/list",
		HandlerFunc: addr.List,
	},
	APIRoute{
		Name:        "api/v1/link/Add.go",
		Method:      "POST",
		Pattern:     "/v1/link/add",
		HandlerFunc: link.Add,
	},
	APIRoute{
		Name:        "api/v1/link/ByName.go",
		Method:      "GET",
		Pattern:     "/v1/link/by/name",
		HandlerFunc: link.ByName,
	},
	APIRoute{
		Name:        "api/v1/link/ByIndex.go",
		Method:      "GET",
		Pattern:     "/v1/link/by/index",
		HandlerFunc: link.ByIndex,
	},
	APIRoute{
		Name:        "api/v1/link/List.go",
		Method:      "GET",
		Pattern:     "/v1/link/list",
		HandlerFunc: link.List,
	},
	APIRoute{
		Name:        "api/v1/link/SetMaster.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/master",
		HandlerFunc: link.SetMaster,
	},
	APIRoute{
		Name:        "api/v1/link/SetNoMaster.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/nomaster",
		HandlerFunc: link.SetNoMaster,
	},
}

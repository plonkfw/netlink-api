package routingv1

import (
	"net/http"

	addrv1 "github.com/plonkfw/netlink-api/api/v1/addr"
	linkv1 "github.com/plonkfw/netlink-api/api/v1/link"
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
		HandlerFunc: addrv1.Add,
	},
	APIRoute{
		Name:        "api/v1/addr/List.go",
		Method:      "GET",
		Pattern:     "/v1/addr/list",
		HandlerFunc: addrv1.List,
	},
	APIRoute{
		Name:        "api/v1/link/Add.go",
		Method:      "POST",
		Pattern:     "/v1/link/add",
		HandlerFunc: linkv1.Add,
	},
	APIRoute{
		Name:        "api/v1/link/ByName.go",
		Method:      "GET",
		Pattern:     "/v1/link/by/name",
		HandlerFunc: linkv1.ByName,
	},
	APIRoute{
		Name:        "api/v1/link/ByIndex.go",
		Method:      "GET",
		Pattern:     "/v1/link/by/index",
		HandlerFunc: linkv1.ByIndex,
	},
	APIRoute{
		Name:        "api/v1/link/Del.go",
		Method:      "POST",
		Pattern:     "/v1/link/del",
		HandlerFunc: linkv1.Del,
	},
	APIRoute{
		Name:        "api/v1/link/List.go",
		Method:      "GET",
		Pattern:     "/v1/link/list",
		HandlerFunc: linkv1.List,
	},
	APIRoute{
		Name:        "api/v1/link/SetDown.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/down",
		HandlerFunc: linkv1.SetDown,
	},
	APIRoute{
		Name:        "api/v1/link/SetUp.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/up",
		HandlerFunc: linkv1.SetUp,
	},
	APIRoute{
		Name:        "api/v1/link/SetMaster.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/master",
		HandlerFunc: linkv1.SetMaster,
	},
	APIRoute{
		Name:        "api/v1/link/SetMasterByIndex.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/masterbyindex",
		HandlerFunc: linkv1.SetMasterByIndex,
	},
	APIRoute{
		Name:        "api/v1/link/SetName.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/name",
		HandlerFunc: linkv1.SetName,
	},
	APIRoute{
		Name:        "api/v1/link/SetNoMaster.go",
		Method:      "POST",
		Pattern:     "/v1/link/set/nomaster",
		HandlerFunc: linkv1.SetNoMaster,
	},
}

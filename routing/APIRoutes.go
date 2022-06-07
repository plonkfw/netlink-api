package routing

import (
	"github.com/plonkfw/netlink-api/api/v1/addr"
	link "github.com/plonkfw/netlink-api/api/v1/link"
	"github.com/plonkfw/netlink-api/types"
)

var routes = types.APIRoutes{
	types.APIRoute{
		Name:        "api/v1/addr/Add.go",
		Method:      "POST",
		Pattern:     "/v1/addr/add",
		HandlerFunc: addr.Add,
	},
	types.APIRoute{
		Name:        "api/v1/link/Add.go",
		Method:      "POST",
		Pattern:     "/v1/link/add",
		HandlerFunc: link.Add,
	},
	types.APIRoute{
		Name:        "api/v1/link/ByName.go",
		Method:      "GET",
		Pattern:     "/v1/link/by-name",
		HandlerFunc: link.ByName,
	},
	types.APIRoute{
		Name:        "api/v1/link/List.go",
		Method:      "GET",
		Pattern:     "/v1/link/list",
		HandlerFunc: link.List,
	},
}

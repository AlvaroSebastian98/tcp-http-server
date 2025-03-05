package utils

import (
	"strings"
)

type Router struct {
	Get func(string, func(RouteContext))
	Post func(string, func(RouteContext))
	Put func(string, func(RouteContext))
	Delete func(string, func(RouteContext))
	// Execute for all HTTP methods
	Use func(string, func(RouteContext))
}

type RouteContext struct {
	Params map[string]string
	Status func(int)
	Send func(string)
}

// Match route defined using the Router and the request path.
// Set data in routeCtx param.
// Returns true if both routes match
func matchRoute(httpMethod string, routerPath string, req HTTPRequest, routeCtx *RouteContext) bool {
	if req.Method != httpMethod && httpMethod != "ALL" {
		return false
	}

	routerSegments := strings.Split(routerPath, "/")
	requestSegments := strings.Split(req.Path, "/")

	if len(routerSegments) != len(requestSegments) {
		return false
	}

	for i, str := range routerSegments {
		isParam := strings.HasPrefix(str, ":")

		if requestSegments[i] == str {
			continue
		} else if isParam {
			param := strings.TrimPrefix(str, ":")
			routeCtx.Params[param] = requestSegments[i]
			continue
		} else {
			return false
		}
	}

	return true
}


// Prepare methods that handle routes and params
func BuildRouter(req HTTPRequest) Router {

	var router Router
	var matchedRoute bool

	router.Get = func(route string, callback func(ctx RouteContext)) {
		if matchedRoute {
			return
		}

		routeCtx := RouteContext{
			Params: make(map[string]string),
		}

		if matchedRoute = matchRoute("GET", route, req, &routeCtx); matchedRoute {
			callback(routeCtx)
		}
	}

	router.Post = func(route string, callback func(ctx RouteContext)) {
		if matchedRoute {
			return
		}

		routeCtx := RouteContext{
			Params: make(map[string]string),
		}

		if matchedRoute = matchRoute("POST", route, req, &routeCtx); matchedRoute {
			callback(routeCtx)
		}
	}

	// Execute for all HTTP methods
	router.Use = func(route string, callback func(ctx RouteContext)) {
		if matchedRoute {
			return
		}

		routeCtx := RouteContext{
			Params: make(map[string]string),
		}

		if route == "*" {
			callback(routeCtx)
		} else if matchedRoute = matchRoute("ALL", route, req, &routeCtx); matchedRoute {
			callback(routeCtx)
		}
	}

	return router

}
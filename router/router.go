package router

import (
	"main/http"
	"net"
	"strings"
)

type Router struct {
	Get func(string, func(Request, Response))
	Post func(string, func(Request, Response))
	Put func(string, func(Request, Response))
	Delete func(string, func(Request, Response))
	// Execute for all HTTP methods
	Use func(string, func(Request, Response))
}

// type RouteContext struct {
// 	Params map[string]string
// 	Status func(int)
// 	Send func(string)
// }

type Request struct {
	http.HTTPRequest
	Params map[string]string
}

type Response struct {
	http.HTTPResponse
}

// Match route defined using the Router and the request path.
// Set data in req param.
// Returns true if both routes match
func matchRoute(httpMethod string, routerPath string, req *Request) bool {
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
			req.Params[param] = requestSegments[i]
			continue
		} else {
			return false
		}
	}

	return true
}


// Prepare methods that handle routes and params
func BuildRouter(conn net.Conn) Router {

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)

	req := Request{
		HTTPRequest: http.ParseRequest(buff[:n]),
	}

	res := Response{
		HTTPResponse: http.BuildResponse(conn),
	}

	var router Router
	var matchedRoute bool

	router.Get = func(route string, callback func(req Request, res Response)) {
		if matchedRoute {
			return
		}

		req.Params = make(map[string]string)

		if matchedRoute = matchRoute("GET", route, &req); matchedRoute {
			callback(req, res)
		}
	}

	router.Post = func(route string, callback func(req Request, res Response)) {
		if matchedRoute {
			return
		}

		req.Params = make(map[string]string)

		if matchedRoute = matchRoute("POST", route, &req); matchedRoute {
			callback(req, res)
		}
	}

	// Execute for all HTTP methods
	router.Use = func(route string, callback func(req Request, res Response)) {
		if matchedRoute {
			return
		}

		req.Params = make(map[string]string)

		if route == "*" {
			callback(req, res)
		} else if matchedRoute = matchRoute("ALL", route, &req); matchedRoute {
			callback(req, res)
		}
	}

	return router

}
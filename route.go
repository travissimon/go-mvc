package mvc

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Parses a path string where a named parameter is surrounded in brackets. Returns
// a regexp that matches the path, and an array of named parameters
// For example, parsePathString("/item/{Id}") would return:
// "/item/(.*)", []string {"Id"}, true
func parsePathString(path string) (string, []string) {
	var buffer bytes.Buffer
	parameters := make([]string, 0, 4)
	parsingParameters := false
	start := 0
	for idx, char := range path {
		if char == '{' {
			parsingParameters = true
			start = idx + 1
			continue
		} else if char == '}' {
			parsingParameters = false
			buffer.WriteString("(.*)")
			parameters = append(parameters, path[start:idx])
			continue
		} else {
			if !parsingParameters {
				buffer.WriteRune(char)
			}
		}
	}

	return buffer.String(), parameters
}

// Creates a Route from a spec (formatted query string, e.g. /Products/{Id})
func NewRoute(name string, spec string, method HttpMethod, controllerFunc ControllerFunc) *Route {
	route := new(Route)
	route.Name = name
	route.Spec = spec
	regStr, params := parsePathString(spec)
	route.Regexp = regexp.MustCompile(regStr)
	route.Params = params
	route.Method = method
	route.Controller = controllerFunc
	return route
}

type Route struct {
	Name       string
	Spec       string
	Regexp     *regexp.Regexp
	Params     []string
	Method     HttpMethod
	Controller ControllerFunc
}

type MismatchedParameterCountError int

func (e MismatchedParameterCountError) Error() string {
	return "Wrong number of parameters parsed: " + string(e)
}

// Gets the values implicitly passed in a path
// E.g., for a route of /p/{Id}, the path:
// /p/23 would return {Id:23}
func (r *Route) GetParameterValues(path string) (url.Values, error) {
	vals := r.Regexp.FindStringSubmatch(path)

	if len(vals) < 2 {
		return url.Values{}, nil
	}

	// The first value is the entire match, which we'll skip
	vals = vals[1:]

	if len(r.Params) != len(vals) {
		return nil, MismatchedParameterCountError(len(vals))
	}

	params := url.Values{}
	for idx, val := range vals {
		params.Add(r.Params[idx], val)
	}
	return params, nil
}

type RouteHandler struct {
	getRoutes    []*Route
	postRoutes   []*Route
	deleteRoutes []*Route
	headRoutes   []*Route
	putRoutes    []*Route
}

func NewRouteHandler() *RouteHandler {
	rh := new(RouteHandler)
	rh.getRoutes = make([]*Route, 0, 10)
	rh.postRoutes = make([]*Route, 0, 10)
	rh.deleteRoutes = make([]*Route, 0, 10)
	rh.headRoutes = make([]*Route, 0, 1)
	rh.putRoutes = make([]*Route, 0, 1)
	return rh
}

func (rh *RouteHandler) AddNewRoute(name string, path string, method HttpMethod, controllerFunc ControllerFunc) {
	rh.AddRoute(NewRoute(name, path, method, controllerFunc))
}

func (rh *RouteHandler) AddRoute(route *Route) {
	switch route.Method {
	case GET:
		rh.getRoutes = append(rh.getRoutes, route)
		break
	case POST:
		rh.postRoutes = append(rh.postRoutes, route)
		break
	case DELETE:
		rh.deleteRoutes = append(rh.deleteRoutes, route)
		break
	case HEAD:
		rh.headRoutes = append(rh.headRoutes, route)
		break
	case PUT:
		rh.putRoutes = append(rh.putRoutes, route)
		break
	}
}

func (rh *RouteHandler) GetRouteFromRequest(r *http.Request) (*Route, bool) {
	route, found := rh.GetRoute(r.URL.Path, r.Method)
	return route, found
}

func (rh *RouteHandler) GetRoute(path string, method string) (*Route, bool) {
	var routes []*Route
	switch strings.ToUpper(method) {
	case "GET":
		routes = rh.getRoutes
		break
	case "POST":
		routes = rh.postRoutes
		break
	case "DELETE":
		routes = rh.deleteRoutes
		break
	case "HEAD":
		routes = rh.headRoutes
		break
	case "PUT":
		routes = rh.putRoutes
		break
	}

	for _, route := range routes {
		if route == nil {
			fmt.Printf("Found nil route :-/")
			continue
		}

		indicies := route.Regexp.FindStringIndex(path)
		if len(indicies) == 2 && indicies[0] == 0 && indicies[1] == len(path) {
			return route, true
		}
	}

	return nil, false
}

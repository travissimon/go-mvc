package mvc

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// parsePathString parses a path string where a named parameter is surrounded
// in brackets. Returns a regexp that matches the path, and an array of named
// parameters. For example, parsePathString("/item/{Id}") would return:
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

// NewRoute creates a Route from a spec (formatted query string, e.g. /Products/{Id})
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

// Route encapsulates infomration needed to route http requests, including the regex
// it uses to match requests, the parameter names to use in value mapping and the
// controller to which the route should be mapped
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

// GetParameterValues extracts the values implicitly passed in a URL.
// E.g., for a route of /p/{Id}, the path:
// /p/23 would return {Id:23}
func (r *Route) GetParameterValues(path string) (url.Values, error) {
	vals := r.Regexp.FindStringSubmatch(path)
	debug(fmt.Sprintf("Parameter values for %s: %v\n", path, vals))

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

// Route Handler multiplexes URL requests to controller functions.
type RouteHandler struct {
	getRoutes    []*Route
	postRoutes   []*Route
	deleteRoutes []*Route
	headRoutes   []*Route
	putRoutes    []*Route
}

// NewRouteHandler initialises and returns a new route handler.
func NewRouteHandler() *RouteHandler {
	rh := new(RouteHandler)
	rh.getRoutes = make([]*Route, 0, 10)
	rh.postRoutes = make([]*Route, 0, 10)
	rh.deleteRoutes = make([]*Route, 0, 10)
	rh.headRoutes = make([]*Route, 0, 1)
	rh.putRoutes = make([]*Route, 0, 1)
	return rh
}

// AddNewRoute associates a route to a controller and adds it to the RouteHandler
func (rh *RouteHandler) AddNewRoute(name string, path string, method HttpMethod, controllerFunc ControllerFunc) {
	rh.AddRoute(NewRoute(name, path, method, controllerFunc))
}

// AddRoute adds an existing route to the RouteHandler
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

// GetRouteFromRequests returns a route from an http.Request instance.
func (rh *RouteHandler) GetRouteFromRequest(r *http.Request) (*Route, bool) {
	route, found := rh.GetRoute(r.URL.Path, r.Method)
	return route, found
}

// GetRoute retrieves a route given a URL and request method
func (rh *RouteHandler) GetRoute(path string, method string) (*Route, bool) {
	debug("GetRoute, path is: " + path + ", method is: " + method)
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
		debug(fmt.Sprintf("Indicies for path '%s': %v", path, indicies))
		if len(indicies) == 2 && indicies[0] == 0 && indicies[1] == len(path) {
			return route, true
		}
	}

	return nil, false
}

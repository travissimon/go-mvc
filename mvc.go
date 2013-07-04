package mvc

import (
	"fmt"
	"net/http"
	"net/url"
)

// This needs to be moved
type WebContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
}

func NewWebContext(w http.ResponseWriter, r *http.Request, s *Session) *WebContext {
	return &WebContext{
		ResponseWriter: w,
		Request:        r,
		Session:        s,
	}
}

type HttpMethod int

// HttpMethods that we will handle
const (
	GET  HttpMethod = iota
	HEAD            // do we care about this?
	POST
	PUT // and this?
	DELETE
)

func NewMvcHandler() *MvcHandler {
	return &MvcHandler{
		Routes:          NewRouteHandler(),
		Sessions:        NewSessionManager(),
		SessionsEnabled: true,
		NotFoundHandler: NotFoundFunc,
	}
}

// TODO: Fill this out
func NotFoundFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404 Not Found")
}

// Http handler function to enable MVC functionality
type MvcHandler struct {
	Routes          *RouteHandler
	Sessions        *SessionManager
	SessionsEnabled bool
	NotFoundHandler func(http.ResponseWriter, *http.Request)
}

// Adds a new route to the MVC handler
func (mvc *MvcHandler) AddRoute(name string, path string, method HttpMethod, controllerFunc ControllerFunc) {
	mvc.Routes.AddNewRoute(name, path, method, controllerFunc)
}

// Main handler function, responsible for multiplexing routes and
// adding session data
func (mvc *MvcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, found := mvc.Routes.GetRouteFromRequest(r)
	if !found {
		mvc.NotFoundHandler(w, r)
		return
	}

	// get any parameter values from the path (from named parameters)
	params, _ := route.GetParameterValues(r.URL.Path)

	// add parameters from form posts
	// Do we need to call parse form?
	MergeValues(params, r.Form)

	var session *Session
	if mvc.SessionsEnabled {
		session = mvc.Sessions.GetSession(w, r)
	}

	ctx := NewWebContext(w, r, session)

	result := route.Controller(ctx, params)

	result.Execute()
}

func MergeValues(vals, valsToMerge url.Values) {
	for key, valSlice := range valsToMerge {
		for _, item := range valSlice {
			vals.Add(key, item)
		}
	}
}

// Return value from a controller
type ControllerResult interface {
	Execute()
}

// Method signature expected for a controller function
type ControllerFunc func(ctx *WebContext, params url.Values) ControllerResult

// Haml definition of a controller result
type HamlTemplate interface {
	Execute(http.ResponseWriter, *http.Request)
}

// Haml Result, containing the template, data to display and
// the web context within which we are working
type HamlResult struct {
	Template HamlTemplate
	Data     interface{}
	Context  *WebContext
}

// Execute method for a Haml template
func (h *HamlResult) Execute() {
	tmpl := h.Template
	ctx := *h.Context
	tmpl.Execute(ctx.ResponseWriter, h.Context.Request)
}

// Utility method to build a Controller Result for a Haml template
func Haml(templ HamlTemplate, data interface{}, ctx *WebContext) ControllerResult {
	return &HamlResult{
		Template: templ,
		Data:     data,
		Context:  ctx,
	}
}

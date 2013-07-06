package mvc

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// This needs to be moved
type WebContext struct {
	mvcHandler     *MvcHandler
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
}

func GetTestControllerParameters() (ctx *WebContext, params url.Values) {
	ctx = NewWebContext(nil, nil, nil, NewSession("Test Session"))
	params = url.Values{}
	return
}

func NewWebContext(m *MvcHandler, w http.ResponseWriter, r *http.Request, s *Session) *WebContext {
	return &WebContext{
		mvcHandler:     m,
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
		Templates:       nil,
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
	Templates       *template.Template // Go Html Templates
	NotFoundHandler func(http.ResponseWriter, *http.Request)
}

// Adds a new route to the MVC handler
func (mvc *MvcHandler) AddRoute(name string, path string, method HttpMethod, controllerFunc ControllerFunc) {
	mvc.Routes.AddNewRoute(name, path, method, controllerFunc)
}

// Adds all Templates to the MVC Hanlder.
// Template value should be the result of calling 'template.ParseFiles(...)'
func (mvc *MvcHandler) SetTemplates(template *template.Template) {
	mvc.Templates = template
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

	ctx := NewWebContext(mvc, w, r, session)

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

// ControllerResult Execute() method for a Haml template
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

// Template result for combining go Templates with a data itme
type TemplateResult struct {
	TemplateName string
	Data         interface{}
	Context      *WebContext
}

// ControllerResult Execute() method for a Haml template
func (t *TemplateResult) Execute() {
	templateName := t.TemplateName
	ctx := *t.Context
	templates := *ctx.mvcHandler.Templates

	err := templates.ExecuteTemplate(ctx.ResponseWriter, templateName, t.Data)
	if err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func Template(templateName string, data interface{}, ctx *WebContext) ControllerResult {
	return &TemplateResult{
		TemplateName: templateName,
		Data:         data,
		Context:      ctx,
	}
}

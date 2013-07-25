package mvc

import (
	"encoding/json"
	"fmt"
	"github.com/travissimon/cache"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

// WebContext provides access to request and session information
type WebContext struct {
	mvcHandler     *MvcHandler
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
	User           *User
}

func (ctx *WebContext) Login(username, password string) (*User, error) {
	sessionId, ipAddress := ctx.getSessionIdAndIPAddress()
	return ctx.mvcHandler.Authenticator.Login(username, password, ipAddress, sessionId)
}

func (ctx *WebContext) CreateUser(username, password, emailAddress string) (user *User, err error) {
	sessionId, ipAddress := ctx.getSessionIdAndIPAddress()
	return ctx.mvcHandler.Authenticator.CreateUser(username, password, emailAddress, ipAddress, sessionId)
}

func (ctx *WebContext) IsUserLoggedIn() bool {
	return ctx.User != nil
}

func (ctx *WebContext) getSessionIdAndIPAddress() (sessionId, ipAddress string) {
	ipAddress = ctx.Request.RemoteAddr
	if idx := strings.Index(ipAddress, ":"); idx > 0 {
		ipAddress = ipAddress[:idx]
	}
	sessionId = ctx.Session.Id
	return
}

// Returns empty WebContext and Values objects for testing
func GetTestControllerParameters() (ctx *WebContext, params url.Values) {
	ctx = NewWebContext(nil, nil, nil, NewSession("Test Session"))
	params = url.Values{}
	return
}

// Creates a new Web Context
func NewWebContext(m *MvcHandler, w http.ResponseWriter, r *http.Request, s *Session) *WebContext {
	ctx := &WebContext{
		mvcHandler:     m,
		ResponseWriter: w,
		Request:        r,
		Session:        s,
	}

	var user *User = nil
	sessionId, ipAddress := ctx.getSessionIdAndIPAddress()
	cacheKey := sessionId + ipAddress
	val, found := m.userCache.Get(cacheKey)
	if found {
		user = val.(*User)
	} else {
		_, user, _ = m.Authenticator.GetAuthentication(sessionId, ipAddress)
	}

	found = (user != nil)
	if found {
		m.userCache.Add(cacheKey, user)
		ctx.User = user
	}

	return ctx
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

// TODO: Fill out this method
func NotFoundFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "404 Not Found")
}

// MvcHandler provides routing, sessions and an mvc patter for the http.Handle() function
type MvcHandler struct {
	Routes          *RouteHandler
	Sessions        *SessionManager
	SessionsEnabled bool
	Templates       *template.Template // Go Html Templates
	NotFoundHandler func(http.ResponseWriter, *http.Request)
	Authenticator   *Authenticator
	userCache       *cache.LRUCache
}

// NewMvcHandler creates an http handler for the MVC package. You can use this handler
// to route requests from Go's http server like this: http.Handle("/", handler)
func NewMvcHandler() *MvcHandler {
	return &MvcHandler{
		Routes:          NewRouteHandler(),
		Sessions:        NewSessionManager(),
		SessionsEnabled: true,
		Templates:       nil,
		NotFoundHandler: NotFoundFunc,
		Authenticator:   NewAuthenticator(),
		userCache:       cache.NewLRUCache(10),
	}
}

// Adds a new route to the MVC handler
func (mvc *MvcHandler) AddRoute(name string, path string, method HttpMethod, controllerFunc ControllerFunc) {
	mvc.Routes.AddNewRoute(name, path, method, controllerFunc)
}

// Adds all (parsed) go Templates to the MVC Hanlder.
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
	r.ParseForm()
	mergeValues(params, r.Form)

	var session *Session
	if mvc.SessionsEnabled {
		session = mvc.Sessions.GetSession(w, r)
	}

	ctx := NewWebContext(mvc, w, r, session)

	result := route.Controller(ctx, params)

	result.Execute()
}

// mergeValues combines url.Values into the the first argument
func mergeValues(vals, valsToMerge url.Values) {
	for key, valSlice := range valsToMerge {
		for _, item := range valSlice {
			vals.Add(key, item)
		}
	}
}

// ControllerResult is the return interface value from a controller.
type ControllerResult interface {
	Execute()
}

// ControllerFunc is the signature expected for a controller function
type ControllerFunc func(ctx *WebContext, params url.Values) ControllerResult

// HamlTemplate is the interface definition for executing a generated Haml template
type HamlTemplate interface {
	SetData(data interface{})
	Execute(http.ResponseWriter, *http.Request)
}

// HamlResult contains the template, data to display and
// the web context within which we are working
type HamlResult struct {
	Template HamlTemplate
	Data     interface{}
	Context  *WebContext
}

// Execute() executes the Haml template and writes the response to the ResponseWriter
func (h *HamlResult) Execute() {
	tmpl := h.Template
	tmpl.SetData(h.Data)
	ctx := *h.Context
	tmpl.Execute(ctx.ResponseWriter, h.Context.Request)
}

// Haml is a utility method to create a controller result for executing Haml templates
func Haml(templ HamlTemplate, data interface{}, ctx *WebContext) ControllerResult {
	return &HamlResult{
		Template: templ,
		Data:     data,
		Context:  ctx,
	}
}

// TemplateResult combines a Go Template and the data for its execution context
type TemplateResult struct {
	TemplateName string
	Data         interface{}
	Context      *WebContext
}

// Execute executes the template and writes the result to the Http response
func (t *TemplateResult) Execute() {
	templateName := t.TemplateName
	ctx := *t.Context
	templates := *ctx.mvcHandler.Templates

	err := templates.ExecuteTemplate(ctx.ResponseWriter, templateName, t.Data)
	if err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

// Template is a utility method to create a controller result for executing go templates
func Template(templateName string, data interface{}, ctx *WebContext) ControllerResult {
	return &TemplateResult{
		TemplateName: templateName,
		Data:         data,
		Context:      ctx,
	}
}

// JsonResult is a ControllerResult for returning Json to the client
type JsonResult struct {
	Data    interface{}
	Context *WebContext
}

// Execute marshalls the Json object and returns the result to the client
func (j *JsonResult) Execute() {
	respWriter := j.Context.ResponseWriter
	respWriter.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(j.Data)
	if err != nil {
		fmt.Fprintf(respWriter, fmt.Sprintf("{Error: '%s'}", err.Error()))
		return
	}
	jsonStr := string(json)
	fmt.Fprintf(j.Context.ResponseWriter, jsonStr)
}

// Json is a utility function for creating ControllerResults to return Json to the client
func Json(data interface{}, ctx *WebContext) ControllerResult {
	return &JsonResult{
		Data:    data,
		Context: ctx,
	}
}

package main

import (
	"fmt"
	"github.com/travissimon/go-mvc"
	"html/template"
	"net/http"
	"net/url"
)

func SessionController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	var val interface{}
	exists := false
	if val, exists = ctx.Session.Get("count"); !exists {
		val = -1
	}

	count := val.(int)
	count++
	ctx.Session.Put("count", count)

	wr := NewSessionWriter()
	return mvc.Haml(wr, count, ctx)
}

func GreetingController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	name := params.Get("name")
	if name == "" {
		name = "there"
	}

	wr := NewHeyWriter()
	return mvc.Haml(wr, name, ctx)
}

// Go templates
type Article struct {
	Author string
	Title  string
	Body   string
}

func GetTestArticle() *Article {
	return &Article{
		Author: "Travis Simon",
		Title:  "Test Article",
		Body:   "This is the body of the article. It's just a test. Enjoy!",
	}
}

var templates = template.Must(template.ParseFiles("article.html"))

func ArticleController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	article := GetTestArticle()
	return mvc.Template("article.html", article, ctx)
}

func JsonController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	article := GetTestArticle()
	return mvc.Json(article, ctx)
}

func LoginController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	sessionId := ctx.Session.Id
	ipAddress := ctx.Request.RemoteAddr
	authenticator := mvc.NewAuthenticator()
	_, user, err := authenticator.GetAuthentication(sessionId, ipAddress)

	if err != nil {
		fmt.Println(err)
	}

	login := new(LoginResult)
	login.LoginSource = Session
	login.User = user
	login.IsLoggedIn = (user != nil)

	wr := NewLoginWriter()
	return mvc.Haml(wr, login, ctx)
}

type LoginSource int

const (
	Session LoginSource = iota
	Form
)

type LoginResult struct {
	IsLoggedIn  bool
	LoginSource LoginSource
	User        *mvc.User
}

func LoginPostController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	submit := params.Get("submit")
	if submit == "Submit" {
		return handleExistingLogin(ctx, params)
	} else {
		return handleNewLogin(ctx, params)
	}
}

func handleExistingLogin(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	ipAddress := ctx.Request.RemoteAddr
	sessionId := ctx.Session.Id
	username := params.Get("txtUsername")
	password := params.Get("txtPassword")
	authenticator := mvc.NewAuthenticator()
	_, user := authenticator.Login(username, password, ipAddress, sessionId)

	login := new(LoginResult)
	login.LoginSource = Form
	login.User = user
	login.IsLoggedIn = (user != nil)

	wr := NewLoginWriter()
	return mvc.Haml(wr, login, ctx)
}

func handleNewLogin(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	ipAddress := ctx.Request.RemoteAddr
	sessionId := ctx.Session.Id
	username := params.Get("txtNewUsername")
	password := params.Get("txtNewPassword")
	email := params.Get("txtNewEmailAddress")

	authenticator := mvc.NewAuthenticator()
	user, err := authenticator.CreateUser(sessionId, ipAddress, username, password, email)

	login := new(LoginResult)
	login.LoginSource = Form
	login.User = user
	login.IsLoggedIn = err == nil

	wr := NewLoginWriter()
	return mvc.Haml(wr, login, ctx)
}

func main() {

	// insert a user into database
	//auth_db := mvc.NewAuthenticationDatabase()
	//userId, err := auth_db.CreateUser("1", "1.1", "tsimon", "tsimon@gmail.com", "Secret!")

	fmt.Println("Listening on: http://localhost:8080/")

	handler := mvc.NewMvcHandler()

	// Set go html-templates
	handler.SetTemplates(templates)
	handler.AddRoute("Article", "/Article", mvc.GET, ArticleController)

	// Add routes
	handler.AddRoute("Homepage", "/", mvc.GET, SessionController)
	handler.AddRoute("Hey", "/Hey/{name}", mvc.GET, GreetingController)
	handler.AddRoute("Article", "/Article", mvc.GET, ArticleController)
	handler.AddRoute("Json", "/json", mvc.GET, JsonController)
	handler.AddRoute("Login", "/login", mvc.GET, LoginController)
	handler.AddRoute("Login Post", "/login", mvc.POST, LoginPostController)

	http.Handle("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}

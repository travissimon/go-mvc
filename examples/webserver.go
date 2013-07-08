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

	wr := NewSessionWriter(count)
	return mvc.Haml(wr, count, ctx)
}

func GreetingController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	name := params.Get("name")
	if name == "" {
		name = "there"
	}

	wr := NewHeyWriter(name)
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

func main() {
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

	http.Handle("/", handler)
	http.ListenAndServe("localhost:8080", nil)
}

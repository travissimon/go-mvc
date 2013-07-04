package main

import (
	"fmt"
	"github.com/travissimon/go-mvc"
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

func main() {
	fmt.Println("Listening on: http://localhost:4040/")

	handl := mvc.NewMvcHandler()
	handl.AddRoute("Homepage", "/", mvc.GET, SessionController)
	handl.AddRoute("Hey", "/Hey/{name}", mvc.GET, GreetingController)
	http.Handle("/", handl)
	http.ListenAndServe("localhost:4040", nil)
}

package mvc

import (
	"log"
	"net/url"
	"testing"
)

func Test_RouteHookup(t *testing.T) {
	log.Printf("Hookup Succeeded - testing Routes")
}

func Test_ParsePathString(t *testing.T) {
	reg, params := parsePathString("/Calendar")
	if reg != "/Calendar" || len(params) != 0 {
		t.Error("parsing not as expected")
	}
	reg, params = parsePathString("/Calendar/{year}")
	if reg != "/Calendar/(.*)" || len(params) != 1 || params[0] != "year" {
		t.Error("parsing not as expected")
	}
	reg, params = parsePathString("/Calendar/{year}/{month}")
	if reg != "/Calendar/(.*)/(.*)" || len(params) != 2 || params[0] != "year" || params[1] != "month" {
		t.Error("parsing not as expected")
	}
	reg, params = parsePathString("/Calendar/{year}/{month}/{day}")
	if reg != "/Calendar/(.*)/(.*)/(.*)" || len(params) != 3 || params[0] != "year" || params[1] != "month" || params[2] != "day" {
		t.Error("parsing not as expected")
	}
}

func Test_MatchingRoutes(t *testing.T) {
	indexRoute := func(ctx *WebContext, vals url.Values) ControllerResult { return nil }

	rh := NewRouteHandler()
	rh.AddNewRoute("Calendar day", "/Calendar/{year}/{month/{day}", GET, indexRoute)
	rh.AddNewRoute("Calendar month", "/Calendar/{year}/{month}", GET, indexRoute)
	rh.AddNewRoute("Calendar year", "/Calendar/{year}", GET, indexRoute)
	rh.AddNewRoute("Calendar index", "/Calendar/", GET, indexRoute)

	r1, found := rh.GetRoute("/Calendar/", "GET")
	if !found || r1 == nil {
		t.Error("Could not find route")
	}
	r2, found := rh.GetRoute("/Calendar/2013/", "GET")
	if !found || r2 == nil {
		t.Error("Could not find route")
	}
	r3, found := rh.GetRoute("/Calendar/2013/", "GET")
	if !found || r3 == nil {
		t.Error("Could not find route")
	}
	r4, found := rh.GetRoute("/Calendar/2013/01/", "GET")
	if !found || r4 == nil {
		t.Error("Could not find route")
	}
	r5, found := rh.GetRoute("/Calendar/2013/01/02/", "GET")
	if !found || r5 == nil {
		t.Error("Could not find route")
	}
	r6, found := rh.GetRoute("/Calendar/", "get")
	if !found || r6 == nil {
		t.Error("Could not find route")
	}
	r7, found := rh.GetRoute("/Calendar/", "post")
	if found || r7 != nil {
		t.Error("Found unexpected route for POST method")
	}
}

func Test_IncompleteMatches(t *testing.T) {
	testRoute := func(ctx *WebContext, vals url.Values) ControllerResult { return nil }

	rh := NewRouteHandler()
	rh.AddNewRoute("Base Route", "/", GET, testRoute)
	_, found := rh.GetRoute("/somestuff", "get")
	if found {
		t.Error("partial matching should not occur")
	}
}

func Test_GetParamValues(t *testing.T) {
	testRoute := func(ctx *WebContext, vals url.Values) ControllerResult { return nil }

	rh := NewRouteHandler()
	rh.AddNewRoute("Product", "/p/{id}", GET, testRoute)
	route, _ := rh.GetRoute("/p/123", "get")

	values, err := route.GetParameterValues("/p/123")
	if err != nil || len(values) != 1 {
		t.Error("Value count not as expected")
	}

	val := values.Get("id")
	if val != "123" {
		t.Error("id value not as expected")
	}

	rh.AddNewRoute("Another test", "/a/{b}/{c}", GET, testRoute)
	route, _ = rh.GetRoute("/a/12/34", "get")

	values, err = route.GetParameterValues("/a/12/34")
	if err != nil || len(values) != 2 {
		t.Error("Value count not as expected")
	}

	if values.Get("b") != "12" {
		t.Error("value 'b' not as expected")
	}

	if values.Get("c") != "34" {
		t.Error("value 'c' not as expected")
	}
}

package main

import (
	"github.com/travissimon/go-mvc"
	"log"
	"testing"
)

func Test_WebserverHookup(t *testing.T) {
	log.Printf("Hookup Succeeded - testing Controllers")
}

func SessionController_Test(t *testing.T) {
	ctx, _ := mvc.GetTestControllerParameters()
	res := SessionController(ctx, nil)

	// check that 'count' has been set in the session
	val, exists := ctx.Session.Get("count")
	if !exists || val.(int) != 0 {
		t.Error("Session variable not set properly")
	}

	// check that our output data is correct
	hRes := res.(*mvc.HamlResult)
	if hRes.Data != 0 {
		t.Error("Data not as expected: %s")
	}
}

func GreetingController_Test(t *testing.T) {
	ctx, params := mvc.GetTestControllerParameters()
	res := GreetingController(ctx, params)

	hamlRes := res.(*mvc.HamlResult)
	if hamlRes.Data != "there" {
		t.Error("Data is not 'there'")
	}

	params.Add("name", "test")
	res = GreetingController(ctx, params)
	hamlRes = res.(*mvc.HamlResult)

	if hamlRes.Data != "test" {
		t.Error("Data is not 'test'")
	}
}

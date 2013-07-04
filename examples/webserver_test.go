package main

import (
	"github.com/travissimon/mvc"
	"log"
	"net/url"
	"testing"
)

func Test_WebserverHookup(t *testing.T) {
	log.Printf("Hookup Succeeded - testing Controllers")
}

func Test_Index(t *testing.T) {
	ctx := NewWebContext(nil, nil, NewSession("test"))
	res := MvcTest(ctx, nil)

	// check that 'count' has been set in the session
	val, exists := ctx.Session.Get("count")
	if !exists || val.(int) != 0 {
		t.Error("Session variable not set properly")
	}

	// check that our output data is correct
	hRes := res.(*HamlResult)
	if hRes.Data != 0 {
		t.Error("Data not as expected: %s")
	}
}

func Test_Hey(t *testing.T) {
	ctx := NewWebContext(nil, nil, nil)
	params := url.Values{}

	res := Hey(ctx, params)

	hamlRes := res.(*HamlResult)
	if hamlRes.Data != "there" {
		t.Error("Data is not 'there'")
	}

	params.Add("name", "test")
	res = Hey(ctx, params)
	hamlRes = res.(*HamlResult)

	if hamlRes.Data != "test" {
		t.Error("Data is not 'test'")
	}
}

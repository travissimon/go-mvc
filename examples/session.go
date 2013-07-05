package main

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"github.com/travissimon/formatting"
	"net/http"
)

func NewSessionWriter(data int) (*SessionWriter) {
	wr := &SessionWriter {
		data: data,
	}
	
	return wr
}

type SessionWriter struct {
	data int
}

func (wr SessionWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *SessionWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data int) {
	html := formatting.NewIndentingWriter(w)

	html.Print(
`<html>
	<head>
		<title>MVC Example</title>
	</head>
	<body>
		<style>
			 body { font-family: Helvetica, Arial, sans-serif; background: #ddd } #content { width: 80%; background:
			#fff; border-color: #333; margin: 20px; padding: 10px; -webkit-border-radius: 10px; -moz-border-radius:
			10px; border-radius: 10px; } pre, code { font-family: Menlo, monospace; font-size: 14px; } pre { line-height:
			18px; }
		</style>
		<div id="content">
			<h1>
				`)

	html.Print("You have visited this page ", data, " times")

	html.Print(
`
			</h1>
			<div> <a href="/">Click here</a> to reload this page to see the magic of sessions.</div>
			<div> Go MVC also provides easy-to-use <a href="/Hey/Mvc User">parameterised routes</a>.</div>
			<p>Dynamic processing:</p>
			<ol>
				`)

	for i := 0; i < 10; i++ {
	html.Print(
`<li>
					`)

	html.Print("Item: ", (i + 1))

	html.Print(
`
				</li>
				`)

	}
	html.Print(
`</ol>
		</div>
	</body>
</html>
`)
}

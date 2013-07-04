package main

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"github.com/travissimon/formatting"
	"net/http"
)

func NewHeyWriter(data string) (*HeyWriter) {
	wr := &HeyWriter {
		data: data,
	}
	
	return wr
}

type HeyWriter struct {
	data string
}

func (wr HeyWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *HeyWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data string) {
	html := formatting.NewIndentingWriter(w)

	html.Print(
`<html>
	<head>
		<title>
			`)

	html.Print("Hey, ", data)

	html.Print(
`
		</title>
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

	html.Print("Hey, ", data)

	html.Print(
`
			</h1>
			<p>
				 This page demostrates the use of parameterised routes. The controller for this page defines this route
				as "/Hey/{name}". 
			</p>
			<p>
				 In the controller, the developer accesses this paramer by name, like this:
				<code> name := params.Get("name")</code>
				<p>
					 Here are some other exciting URLs to try:
					<ul>
						<li><a href="/Hey/World">/Hey/World</a></li>
						<li><a href="/Hey/">/Hey</a> (default value set in controller)</li>
					</ul>
				</p>
				<p> Return <a href="/">home</a></p>
			</p>
		</div>
	</body>
</html>
`)
}

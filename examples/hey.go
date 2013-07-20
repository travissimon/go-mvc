package main

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"fmt"
	"html/template"
	"net/http"
)

func NewHeyWriter() (*HeyWriter) {
	wr := &HeyWriter{}
	
	for idx, pattern := range HeyTemplatePatterns {
		tmpl, err := template.New("HeyTemplates" + string(idx)).Parse(pattern)
		if err != nil {
			fmt.Errorf("Could not parse template: %d", idx)
			panic(err)
		}
		HeyTemplates = append(HeyTemplates, tmpl)
	}
	return wr
}

type HeyWriter struct {
	data string
}

func (wr *HeyWriter) SetData(data interface{}) {
	wr.data = data.(string)
}

var HeyHtml = [...]string{
`<html>
	<head>
		<title>
			`,
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
				`,
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
`,
}

var HeyTemplatePatterns = []string{
	`Hey, {{.}}`,
	`Hey, {{.}}`,
}

var HeyTemplates = make([]*template.Template, 0, len(HeyTemplatePatterns))

func (wr HeyWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *HeyWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data string) {
	var err error = nil
	fmt.Fprint(w, HeyHtml[0])
	err = HeyTemplates[0].Execute(w, data)
	handleHeyError(err)
	fmt.Fprint(w, HeyHtml[1])
	err = HeyTemplates[1].Execute(w, data)
	handleHeyError(err)
	fmt.Fprint(w, HeyHtml[2])
}

func handleHeyError(err error) {
	if err != nil {fmt.Println(err)}}
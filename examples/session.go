package main

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"fmt"
	"html/template"
	"net/http"
)

func NewSessionWriter() (*SessionWriter) {
	wr := &SessionWriter{}
	
	for idx, pattern := range SessionTemplatePatterns {
		tmpl, err := template.New("SessionTemplates" + string(idx)).Parse(pattern)
		if err != nil {
			fmt.Errorf("Could not parse template: %d", idx)
			panic(err)
		}
		SessionTemplates = append(SessionTemplates, tmpl)
	}
	return wr
}

type SessionWriter struct {
	data int
}

func (wr *SessionWriter) SetData(data interface{}) {
	wr.data = data.(int)
}

var SessionHtml = [...]string{
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
				`,
				`
			</h1>
			<div> <a href="/">Click here</a> to reload this page to see the magic of sessions.</div>
			<div> Go MVC also provides easy-to-use <a href="/Hey/Mvc User">parameterised routes</a>.</div>
			<div> You can also use <a href="/Article">go templates</a></div>
			<p>Dynamic processing:</p>
			<ol>
				`,
				`
				<li>
					`,
					`
				</li>
			</ol>
		</div>
	</body>
</html>
`,
}

var SessionTemplatePatterns = []string{
	`You have visited this page {{.}} times`,
	`Item: {{.}}`,
}

var SessionTemplates = make([]*template.Template, 0, len(SessionTemplatePatterns))

func (wr SessionWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *SessionWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data int) {
	var err error = nil
	fmt.Fprint(w, SessionHtml[0])
	err = SessionTemplates[0].Execute(w, data)
	handleSessionError(err)
	fmt.Fprint(w, SessionHtml[1])
	for i := 0; i < 10; i++ {
		fmt.Fprint(w, SessionHtml[2])
		err = SessionTemplates[1].Execute(w, data)
		handleSessionError(err)
	}
	fmt.Fprint(w, SessionHtml[3])
}

func handleSessionError(err error) {
	if err != nil {fmt.Println(err)}}
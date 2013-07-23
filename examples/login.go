package main

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"fmt"
	"html/template"
	"net/http"
)

func NewLoginWriter() (*LoginWriter) {
	wr := &LoginWriter{}
	
	for idx, pattern := range LoginTemplatePatterns {
		tmpl, err := template.New("LoginTemplates" + string(idx)).Parse(pattern)
		if err != nil {
			fmt.Errorf("Could not parse template: %d", idx)
			panic(err)
		}
		LoginTemplates = append(LoginTemplates, tmpl)
	}
	return wr
}

type LoginWriter struct {
	data *LoginResult
}

func (wr *LoginWriter) SetData(data interface{}) {
	wr.data = data.(*LoginResult)
}

var LoginHtml = [...]string{
`<html>
	<head>
		<title>Login example</title>
	</head>
	<body>
		<style>
			 body { font-family: Helvetica, Arial, sans-serif; background: #ddd } #content { width: 80%; background:
			#fff; border-color: #333; margin: 20px; padding: 10px; -webkit-border-radius: 10px; -moz-border-radius:
			10px; border-radius: 10px; } pre, code { font-family: Menlo, monospace; font-size: 14px; } pre { line-height:
			18px; }
		</style>
		<div id="content">
			<h1>Login</h1>
			`,
			`
			<form method="post">
				<div id="login">
					`,
					`
					<div class="error">The username and password you supplied were not correct. Please try again.</div>
					`,
					`
					<p>Username</p>
					<input type="text" name="txtUsername"></input>
					<p>Password</p>
					<input type="password" name="txtPassword"></input>
					<p>
						<button name="submit" type="submit" value="Submit"> Submit</button>
					</p>
				</div>
				<div id="new_login">
					<p>Username</p>
					<input type="text" name="txtNewUsername"></input>
					<p>Password</p>
					<input type="password" name="txtNewPassword"></input>
					<p>Recovery Email Address</p>
					<input type="text" name="txtNewEmailAddress"></input>
					<p>
						<button name="submit" type="submit" value="New"> Submit</button>
					</p>
				</div>
			</form>
			`,
			`
			<div>
				`,
				`
			</div>
		</div>
	</body>
</html>
`,
}

var LoginTemplatePatterns = []string{
	`Hey, you're logged in, {{.User.Username}}`,
}

var LoginTemplates = make([]*template.Template, 0, len(LoginTemplatePatterns))

func (wr LoginWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *LoginWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data *LoginResult) {
	var err error = nil
	fmt.Fprint(w, LoginHtml[0])
	if !data.IsLoggedIn {
		fmt.Fprint(w, LoginHtml[1])
		if data.LoginSource == Form {
			fmt.Fprint(w, LoginHtml[2])
		}
		fmt.Fprint(w, LoginHtml[3])
	} else {
		fmt.Fprint(w, LoginHtml[4])
		err = LoginTemplates[0].Execute(w, data)
		handleLoginError(err)
	}
	fmt.Fprint(w, LoginHtml[5])
}

func handleLoginError(err error) {
	if err != nil {fmt.Println(err)}}
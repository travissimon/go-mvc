/*
Package MVC provides a lightweight, testable and efficient Model-View-Controller framework.
The main features are testable controllers, session support, parameterised routes
and support for a Haml-like template language.

[![Build Status](https://travis-ci.org/travissimon/go-mvc.png)](https://travis-ci.org/travissimon/go-mvc)

The main component is the mvc handler, instantiated like this:

    handler := mvc.NewMvcHandler()

Once you have have an MVC handler, you can begin adding routes, which can contained
named parameters. Parameters can be specified by surrounding the parameter name
with curly brackets, like this:

    handler.AddRoute("Hello route", "/Hello/{name}", mvc.GET, GreetingController)

The code above creates a route named 'Hello route'. The route will match requests to
'/Hello/...', and will associate the value passed in the URL with a parameter called
'name'. Finally, this route will be handled by a function named GreetingController,
which might look like this:

    func GreetingController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
        name := params.Get("name")
        if name == "" {
            name = "there"
        }

        wr := NewHelloWorldWriter(name)
        return mvc.Haml(wr, name, ctx)
    }

The Greeting Controller retrieves the parameter called 'name' and passes it to a
Haml template to render to the user. Note that the function returns a
mvc.ControllerResult object, which allows us to call this controller method
in a test scenario and test the resulting ControllerResult.


*/
package mvc

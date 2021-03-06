Go-MVC - Testable MVC Framework for Go
======================================

Go-MVC is a lightweight MVC framework for the Go language. Its goals are
to provide an efficient, testable framework for creating web applications.

[![Build Status](https://travis-ci.org/travissimon/go-mvc.png)](https://travis-ci.org/travissimon/go-mvc)

Go-MVC's main features are:
* Routing with Named Parameters
* Testable controllers
* Session support
* Support for different templating engines and Json requests

Development is still underway, so contributions are encouraged.

Routing with Named Parameters
=============================

Go-MVC lets you create routes with named parameters. For example, you
can create a route as follows:

	mvcHandle := mvc.NewMvcHandler()
	mvcHandle.AddRoute("Hello", "/Hello/{name}", mvc.GET, GreetingController)
	http.Handle("/", mvcHandle)
	http.ListenAndServe("localhost:8080", nil)

The code above creates a route named 'Hello' that intercepts GET requests for
routes that match 'localhost:8080/Hello/*' and directs them to a controller
named 'GreetingController'.

Controller methods accept a parameter called 'params' which contains query string
parameters, form post parameters and values from named routes. You can then
access the values inside of your controller. For example, if a request is made to
'http://localhost:8080/Hello/World', you could access the value of 'name' as follows:

	name := params.Get("name")

Testable Controllers
====================

Go-MVC provides a clean separation between controller logic and view
presentation, allowing for easy testing. For example, the following
controller converts a URL part (specified in a route: see below) to uppercase:

    func ToUpperController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResults {
	    in := params.Get("input")
		upper := strings.ToUpper(in)
		hamlWriter := NewViewWriter(upper)
		return mvc.Haml(hamlWriter, upper, ctx)
    }

The code above could be tested like this:

    func ToUpperController_Test(t *testing.T) {
	    ctx, params = GetTestControllerParameters()
		params.Add("input", "tesTIng")

		result := ToUpperController(ctx, params)
		hamlRes := res.(*HamlResult)
		if hamlRes.Data != "TESTING" {
		    t.Error("Data not as expected: " + hamlRes.Data)
		}
    }

Sessions
========

Go-MVC provides server-side sessions transparently. As a developer you 
just add and remove items from the session:

	var val interface{}
	if val, exists = ctx.Session.Get("count"); !exists {
		val = -1
	}

	count := val.(int)
	count++
	ctx.Session.Put("count", count)

Copyright / License
===================
Copyright 2013 Travis Simon

Licensed under the GNU General Public License Version 2.0 (or later); you may not use this work except in compliance with the License. You may obtain a copy of the License in the LICENSE file, or at:

http://www.gnu.org/licenses/old-licenses/gpl-2.0.txt

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
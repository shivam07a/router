/*Package router contains structs and functions for a generic http muxer.

  Usage

  The main usage of this package is through the struct "Router". Like in a simple server,
  instead of "http.HandleFunc", "router.HandleFunc" is used. The function pattern, however,
  varies.
      func (w http.ResponseWriter, r *http.Request) => func (w http.ResponseWriter, r *http.Request, params router.Params)

  A sample usage is written below :

      r := router.NewRouter()
      r.HandleFunc("/:username", UserPage)
      r.HandleFunc("/:username//////:repo", Home)
      r.HandleFunc("/", Home)
      http.ListenAndServe(":8080", r)

  In the handler function, the params is a simple map of string to string. Consider the above
  example.

      func UserPage(w http.ResponseWriter, r *http.Request, params router.Params) {
              username := params["username"]
              ...
      }

  All the URL Path variables can be accessed from the params parameter.

  Normalised URLs

  In the second route in the above example, the route "/:username//////:repo" will automatically
  be normalised to "/:username/:repo".
*/
package router

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// A Route contains string pattern corresponding to the URLs it has to be matched and
// the respective handler function for that URL
type Route struct {
	pattern string
	handler func(http.ResponseWriter, *http.Request, Params)
}

// Match method is used to check if the given path matches the pattern.
// It returns a boolean.
func (route *Route) Match(path string) bool {
	psp := strings.Split(route.pattern, "/")
	sp := strings.Split(path, "/")
	if len(sp) != len(psp) {
		fmt.Println("Lengths Don't match")
		return false
	}
	for index, val := range psp {
		if len(val) > 0 {
			if val[0] == ':' {
				continue
			}
		}
		if val != sp[index] {
			return false
		}
	}
	return true
}

// GetParams is used to get the URL parameters for a given pattern.
// It returns Params and an error.
func (route *Route) GetParams(path string) (Params, error) {
	params := make(Params, 0)
	if !route.Match(path) {
		return nil, errors.New("The given path doesn't match the pattern")
	}
	psp := strings.Split(route.pattern, "/")
	sp := strings.Split(path, "/")
	for index, val := range psp {
		if len(val) > 0 {
			if val[0] == ':' {
				params[val[1:]] = sp[index]
			}
		}
	}
	return params, nil
}

// Router is a struct that satisfies the handler interface.
// It can be used as a http.Handler in the http.ListenAndServe.
type Router struct {
	routes []Route
}

// Params is used to store the URL parameters
// This is a simple key value map(string).
type Params map[string]string

// NewRouter is used for making a new Router variables
func NewRouter() *Router {
	r := &Router{}
	r.routes = make([]Route, 0)
	return r
}

// ServeHTTP method makes the router a qualified handler.
// This is where all the pattern based routing takes place.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = normalizePath(r.URL.Path)
	for _, val := range rt.routes {
		if !val.Match(r.URL.Path) {
			continue
		}
		params, err := val.GetParams(r.URL.Path)
		if err != nil {
			http.Error(w, "Internal Server Error", 503)
			return
		}
		fmt.Println(params)
		val.handler(w, r, params)
		return
	}
	http.NotFound(w, r)
}

// HandleFunc method is used to register Routes for the corresponding handler function.
// This can be used like http.HandleFunc, but the blueprint of the handler function is different.
//     func(http.ResponseWriter, *http.Request, Params)
func (rt *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request, Params)) {
	pattern = normalizePath(pattern)
	for _, val := range rt.routes {
		if val.pattern == pattern {
			panic("Same Route Exists")
		}
	}
	nr := Route{pattern, handler}
	rt.routes = append(rt.routes, nr)
	fmt.Println(rt.routes)
}

// normalizePath is used to sanitize the pattern passed to the HandleFunc.
func normalizePath(path string) string {
	ret := "/"
	sp := strings.Split(path, "/")
	for _, val := range sp {
		if val == "" || val == "." || val == ".." {
			continue
		}
		ret += val + "/"
	}
	return ret[:len(ret)-1]
}

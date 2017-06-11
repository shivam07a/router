package router

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	pattern string
	handler func(http.ResponseWriter, *http.Request, Params)
}

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

type Router struct {
	routes []Route
}

type Params map[string]string

func NewRouter() *Router {
	r := &Router{}
	r.routes = make([]Route, 0)
	return r
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (rt *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request, Params)) {
	pattern = normalizePath(pattern)
	for _, val := range rt.routes {
		if val.pattern == pattern {
			panic("Same Route Exists")
			return
		}
	}
	nr := Route{pattern, handler}
	rt.routes = append(rt.routes, nr)
	fmt.Println(rt.routes)
}

func normalizePath(path string) string {
	ret := "/"
	sp := strings.Split(path, "/")
	for _, val := range sp {
		if val == "" {
			continue
		}
		ret += val + "/"
	}
	return ret[:len(ret)-1]
}

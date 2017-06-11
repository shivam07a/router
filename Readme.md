# router
A simple bare minimum pattern based URL routing package for golang.

## Usage
This package can be used as a drop in replacement for the default http router provided by the `net/http` package
with some changes. Refer [documentation](https://godoc.org/github.com/shivam07a/router) for more.

#### Sample Usage

```go
package main

import (
	"net/http"
	
	"github.com/shivam07a/router"
)

func handleName (w http.ResponseWriter , r *http.Request, p utils.Params) {
	name, ok := p["name"]
	if !ok {
		http.Error(w, "Internal Server Error", 503)
		return
	}
	w.Write([]byte("Hello, " + name + "!"))
}

func home (w http.ResponseWriter , r *http.Request, p utils.Params) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	r := router.NewRouter()
	r.HandleFunc("/:name", handleName)
	r.HandleFunc("/", home)
	http.ListenAndServe(":8080", r)
}
```

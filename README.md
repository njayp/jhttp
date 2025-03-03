
# jhttp

`jhttp` is a Go package for encoding and decoding JSON in HTTP requests and responses with validation support.

## Installation

To install the package, run:

```sh
go get github.com/yourusername/jhttp
```

## Usage

### Decoding JSON

To decode JSON from an `http.Request` and validate it, use the `Decode` function:

```go
package main

import (
	"net/http"
	"github.com/yourusername/jhttp"
)

type Request struct {
	Name string `json:"name"`
}

func (r Request) Valid() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := jhttp.Decode[Request](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process the request...
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

### Encoding JSON

To encode a struct into JSON and write it to an `http.ResponseWriter`, use the `Encode` function:

```go
package main

import (
	"net/http"
	"github.com/yourusername/jhttp"
)

type Response struct {
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp := Response{Message: "Hello, World!"}
	if err := jhttp.Encode(w, http.StatusOK, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```


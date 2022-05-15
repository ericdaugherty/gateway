Package gateway provides a drop-in replacement for net/http's `ListenAndServe` for use in [AWS Lambda](https://aws.amazon.com/lambda/) & [Lambda Function URLs](https://docs.aws.amazon.com/lambda/latest/dg/lambda-urls.html), simply swap it out for `gateway.ListenAndServe`. 

This project is a fork of [apex/Gateway](https://github.com/apex/gateway) which provides the same functionality for [API Gateway](https://aws.amazon.com/api-gateway/) and [HTTP APIs](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api.html).

# Installation

```
go get github.com/ericdaugherty/gateway
```

# Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ericdaugherty/gateway"
)

func main() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	// example retrieving values from the api gateway proxy request context.
	requestContext, ok := gateway.RequestContext(r.Context())
	if !ok || requestContext.Authorizer["sub"] == nil {
		fmt.Fprint(w, "Hello World from Go")
		return
	}

	userID := requestContext.Authorizer["sub"].(string)
	fmt.Fprintf(w, "Hello %s from Go", userID)
}
```

---

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/ericdaugherty/gateway)
[![License](https://img.shields.io/github/license/ericdaugherty/gateway)](https://github.com/ericdaugherty/gateway/blob/master/LICENSE)

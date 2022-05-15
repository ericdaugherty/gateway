package gateway_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ericdaugherty/gateway"
	"github.com/tj/assert"
)

func Example() {
	http.HandleFunc("/", hello)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World from Go")
}

func TestGateway_Invoke(t *testing.T) {

	e := []byte(`{"version": "2.0", "rawPath": "/pets/luna", "requestContext": {"http": {"method": "POST"}}}`)

	gw := gateway.NewGateway(http.HandlerFunc(hello))

	payload, err := gw.Invoke(context.Background(), e)
	assert.NoError(t, err)
	res := events.LambdaFunctionURLResponse{}
	err = json.Unmarshal(payload, &res)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World from Go\n", res.Body)
	assert.Equal(t, "text/plain; charset=utf8", res.Headers["Content-Type"])
	assert.Equal(t, 200, res.StatusCode)
}

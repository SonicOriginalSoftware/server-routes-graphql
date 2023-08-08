package graphql_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"git.sonicoriginal.software/routes/graphql.git"
	"git.sonicoriginal.software/server.git/v2"
)

const portEnvKey = "TEST_PORT"

var (
	certs []tls.Certificate
	mux   *http.ServeMux = nil
)

func TestHandler(t *testing.T) {
	route := graphql.New(mux)

	t.Logf("Handler registered for route [%v]\n", route)

	ctx, cancelFunction := context.WithCancel(context.Background())
	address, serverErrorChannel := server.Run(ctx, &certs, mux, portEnvKey)

	t.Logf("Serving on [%v]\n", address)

	// TODO modify the request to send a proper graphql request
	url := fmt.Sprintf("http://%v%v", address, route)

	t.Logf("Requesting [%v]\n", url)

	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelFunction()

	serverError := <-serverErrorChannel
	if serverError.Close != nil {
		t.Fatalf("Error closing server: %v", serverError.Close.Error())
	}
	contextError := serverError.Context.Error()

	t.Logf("%v\n", contextError)

	if contextError != server.ErrContextCancelled.Error() {
		t.Fatalf("Server failed unexpectedly: %v", contextError)
	}

	t.Log("Response:")
	t.Logf("  Status code: %v", response.StatusCode)
	t.Logf("  Status text: %v", response.Status)

	if response.Status != http.StatusText(http.StatusBadRequest) && response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Server returned: %v", response.Status)
	}
}

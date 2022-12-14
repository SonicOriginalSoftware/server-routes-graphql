package graphql_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"git.sonicoriginal.software/routes/graphql"
	lib "git.sonicoriginal.software/server"
)

var certs []tls.Certificate

func TestHandler(t *testing.T) {
	route := fmt.Sprintf("localhost/%v/", graphql.Name)
	t.Setenv(fmt.Sprintf("%v_SERVE_ADDRESS", strings.ToUpper(graphql.Name)), route)

	_ = graphql.New()

	ctx, cancelFunction := context.WithCancel(context.Background())

	exitCode, address := lib.Run(ctx, certs)
	defer close(exitCode)

	// TODO modify the request to send a proper graphql request
	url := fmt.Sprintf("http://%v/%v/", address, graphql.Name)
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	cancelFunction()

	if returnCode := <-exitCode; returnCode != 0 {
		t.Fatalf("Server errored: %v", returnCode)
	}

	if response.Status != http.StatusText(http.StatusBadRequest) && response.StatusCode != http.StatusBadRequest {
		t.Fatalf("Server returned: %v", response.Status)
	}
}

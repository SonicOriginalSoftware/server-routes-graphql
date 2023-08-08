//revive:disable:package-comments

package graphql

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"

	"git.sonicoriginal.software/logger.git"
	"git.sonicoriginal.software/server.git/v2"
)

// name is the name used to identify the service
const name = "gql"

type postData struct {
	Variables map[string]interface{} `json:"variables"`
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
}

// Handler handles GraphQL API requests
type handler struct {
	logger logger.Log
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var jsonData postData
	if err := json.NewDecoder(request.Body).Decode(&jsonData); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	result := graphql.Do(graphql.Params{
		Context:        request.Context(),
		Schema:         Schema,
		RequestString:  jsonData.Query,
		VariableValues: jsonData.Variables,
		OperationName:  jsonData.Operation,
	})

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		h.logger.Error("Could not write result to response: %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

// New returns a new Handler
func New(mux *http.ServeMux) (route string) {
	logger := logger.New(
		name,
		logger.DefaultSeverity,
		os.Stdout,
		os.Stderr,
	)

	return server.RegisterHandler(name, &handler{logger}, mux)
}

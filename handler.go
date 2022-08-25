//revive:disable:package-comments

package graphql

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"

	"git.nathanblair.rocks/server/logging"
)

const prefix = "graphql"

type postData struct {
	Variables map[string]interface{} `json:"variables"`
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
}

// Handler handles GraphQL API requests
type Handler struct {
	logger *logging.Logger
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var jsonData postData
	if err := json.NewDecoder(request.Body).Decode(&jsonData); err != nil {
		writer.WriteHeader(400)
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
		handler.logger.Error("Could not write result to response: %s", err)
	}
}

// Prefix is the subdomain prefix
func (handler *Handler) Prefix() string {
	return prefix
}

// New returns a new Handler
func New() *Handler {
	return &Handler{
		logger: logging.New(prefix),
	}
}

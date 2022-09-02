//revive:disable:package-comments

package graphql

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"

	"git.sonicoriginal.software/server/handlers"
	"git.sonicoriginal.software/server/logging"
)

// Name is the name used to identify the service
const Name = "gql"

type postData struct {
	Variables map[string]interface{} `json:"variables"`
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
}

// Handler handles GraphQL API requests
type Handler struct {
	logger logging.Log
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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
		handler.logger.Error("Could not write result to response: %s", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
}

// New returns a new Handler
func New() (handler *Handler) {
	logger := logging.New(Name)
	handler = &Handler{logger}
	handlers.Register(Name, "", Name, handler, logger)

	return
}

package main

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daniilsolovey/graphql-service/graph"
	"github.com/daniilsolovey/graphql-service/graph/generated"
	"github.com/daniilsolovey/graphql-service/internal/config"
	"github.com/daniilsolovey/graphql-service/internal/database"
	"github.com/daniilsolovey/graphql-service/internal/operator"
	"github.com/docopt/docopt-go"
	"github.com/go-chi/chi"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

var version = "[manual build]"

var usage = `graphql-service

Simple service for study

Usage:
  graphql-service [options]

Options:
  -c --config <path>                Read specified config file. [default: config.yaml]
  --debug                           Enable debug messages.
  -v --version                      Print version.
  -h --help                         Show this help.
`

const defaultPort = "8080"

func main() {
	args, err := docopt.ParseArgs(
		usage,
		nil,
		"graphql-service "+version,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("version", version),
		"graphql-service started",
	)

	if args["--debug"].(bool) {
		log.SetLevel(log.LevelDebug)
	}

	log.Infof(nil, "loading configuration file: %q", args["--config"].(string))

	config, err := config.Load(args["--config"].(string))
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("database", config.Database.Name),
		"connecting to the database",
	)

	database := database.NewDatabase(
		config.Database.Name, config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password,
	)

	newOperator := operator.NewOperator(
		config, database,
	)

	port := config.Server.Port
	router := chi.NewRouter()
	router.Use(CreateMiddleware())
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Operator: newOperator,
			},
		},
	))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Infof(nil, "connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}

func CreateMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")
			ctx := context.WithValue(request.Context(), "Authorization", token)
			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request)
		})
	}
}

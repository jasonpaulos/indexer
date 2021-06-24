package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/algorand/indexer/api/graph/generated"
	"github.com/labstack/echo/v4"
)

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router *echo.Echo, m ...echo.MiddlewareFunc) {

	graphqlHandler := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &Resolver{}},
		),
	)

	router.POST("/query", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	}, m...)

	playgroundHandler := playground.Handler("GraphQL", "/query")

	router.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	}, m...)
}

package api

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/algorand/indexer/api/graph/generated"
	"github.com/labstack/echo/v4"
)

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterGraphQLHandlers(router *echo.Echo, si *ServerImplementation, m ...echo.MiddlewareFunc) {

	graphqlHandler := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &Resolver{si: si}},
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

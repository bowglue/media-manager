package server

// import (
// 	"context"
// 	"mms/common/gql"
// 	"net/http"
// 	"time"

// 	"github.com/99designs/gqlgen/graphql"
// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/handler/extension"
// 	"github.com/99designs/gqlgen/graphql/handler/lru"
// 	"github.com/99designs/gqlgen/graphql/handler/transport"
// 	"github.com/99designs/gqlgen/graphql/playground"
// 	"github.com/vektah/gqlparser/v2/gqlerror"
// )

// func (s *Server) graphqlHandler() http.Handler {
// 	c := gql.Config{Resolvers: s.resolver}
// 	srv := handler.New(gql.NewExecutableSchema(c))

// 	srv.AddTransport(transport.Options{})
// 	srv.AddTransport(transport.GET{})
// 	srv.AddTransport(transport.POST{})
// 	srv.AddTransport(transport.MultipartForm{})
// 	srv.AddTransport(transport.Websocket{KeepAlivePingInterval: 10 * time.Second})

// 	srv.Use(extension.Introspection{})
// 	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

// 	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
// 		return graphql.DefaultErrorPresenter(ctx, err)
// 	})

// 	return srv
// }

// func (s *Server) playgroundHandler() http.Handler {
// 	return playground.Handler("GraphQL Playground", "/query")
// }

// func (s *Server) setupHandlers() {

// 	s.router.Handle("/query", s.graphqlHandler())
// 	s.router.Handle("/", s.playgroundHandler())
// }

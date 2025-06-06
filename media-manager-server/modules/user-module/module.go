package usermodule

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"mms/common/database/repository"
	"mms/modules/user-module/gql"
	"mms/modules/user-module/gql/resolvers"
	"mms/modules/user-module/services"
)

// Only public thing: the module interface
type UserModule struct {
	handler http.Handler // only exposed element
}

// Only exported function
func NewUserModule(queries *repository.Queries) *UserModule {
	service := services.NewUserService(queries)
	resolver := &resolvers.Resolver{UserService: service}

	// 4. Create gqlgen server config
	cfg := gql.Config{Resolvers: resolver}
	srv := handler.New(gql.NewExecutableSchema(cfg))

	// 5. Add extensions & transports (optional but common)
	// srv.AddTransport(transport.Options{})
	// srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	// srv.AddTransport(transport.MultipartForm{})
	// srv.AddTransport(transport.Websocket{KeepAlivePingInterval: 10 * time.Second})

	srv.Use(extension.Introspection{})
	// srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	// srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
	// 	return graphql.DefaultErrorPresenter(ctx, err)
	// })

	return &UserModule{
		handler: srv,
	}
}

// Public accessor
func (m *UserModule) Handler() http.Handler {
	return m.handler
}

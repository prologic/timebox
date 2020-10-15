package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"banking/model"
	"banking/server"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/kode4food/timebox/store/local"
)

// Environment variables
const (
	PORT     = "PORT"
	DATABASE = "DATABASE"

	defaultPort     = "8080"
	defaultDatabase = "events.db"
)

func main() {
	port := envLookup(PORT, defaultPort)
	database := envLookup(DATABASE, defaultDatabase)

	base, _ := os.Getwd()

	// Create a Local Event Store (using BitCask) and register a set of
	// Decoders for our application's event types (ex: AccountOpened)
	db, err := local.Open(
		local.Path(path.Join(base, database)),
		local.Decoder(model.TypedInstantiator.Decoder()),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	router := mux.NewRouter()
	router.Use(setAllowOriginMiddleware)

	// srv gets a ResolverRoot that utilizes our db for event sourcing
	srv := makeServer(server.NewResolver(db))

	// Create routes for both the GraphQL playground and our resolvers
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	// Start listening for requests
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func setAllowOriginMiddleware(next http.Handler) http.Handler {
	// This middleware is not secure, best to not use it in production
	// unless combining it with strong authentication and authorization
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}

func makeServer(root server.ResolverRoot) *handler.Server {
	// Create a GraphQL executable schema instance and register our
	// resolver root with it
	srv := handler.New(
		server.NewExecutableSchema(server.Config{
			Resolvers: root,
		}),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	return srv
}

func envLookup(name string, defaultValue string) string {
	if res, ok := os.LookupEnv(name); ok {
		return res
	}
	return defaultValue
}

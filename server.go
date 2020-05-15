//go:generate go generate ./resources/

package main

import (
	"fmt"
	"log"
	"net/http"

	"messenger/appconfig"
	"messenger/handlers"
	"messenger/models"
	"messenger/resolvers"
	"messenger/resources"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func main() {
	cfg := appconfig.Config()
	session := connectDB(cfg.DB.URI, cfg.DB.Name)

	schemaSource, err := resources.GetSchema()
	if err != nil {
		log.Fatalln(err)
	}

	schema := graphql.MustParseSchema(
		schemaSource,
		resolvers.NewRootResolver(session),
		graphql.UseStringDescriptions())
	graphQLHandler := graphqlws.NewHandlerFunc(schema, &handlers.GraphQL{
		Schema:  schema,
		Session: session,
	})

	http.HandleFunc("/graphql", graphQLHandler)
	http.Handle("/", &handlers.GraphiQL{})
	http.Handle("/play", &handlers.GraphQLPlayground{})
	http.Handle("/static/", http.FileServer(&assetfs.AssetFS{
		Asset:     resources.Asset,
		AssetDir:  resources.AssetDir,
		AssetInfo: resources.AssetInfo,
		Prefix:    "",
	}))

	addr := fmt.Sprintf("0.0.0.0:%v", cfg.Port)
	log.Println(addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func connectDB(uri, dbname string) *r.Session {
	session, err := r.Connect(r.ConnectOpts{
		Address:  uri,
		Database: dbname,
	})
	if err != nil {
		panic(err)
	}
	createDBTables(session)
	return session
}

func createDBTables(session *r.Session) {
	for _, tbl := range models.Tables() {
		res, err := r.TableList().Contains(tbl).Run(session)
		if err != nil {
			panic(err)
		}
		exists := false
		res.One(&exists)
		if !exists {
			r.TableCreate(tbl).RunWrite(session)
		}
	}
}

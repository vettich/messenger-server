package handlers

import (
	"net/http"

	"messenger/resources"
)

type GraphiQL struct{}

func (g GraphiQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !g.isSupport(r.Method) {
		w.Write([]byte("only GET requests are supported"))
		return
	}

	graphiql, err := resources.Asset("graphiql.html")
	if err != nil {
		w.Write([]byte(""))
		return
	}
	w.Write(graphiql)
}

func (g GraphiQL) isSupport(method string) bool {
	return (method == "GET")
}

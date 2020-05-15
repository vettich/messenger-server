package handlers

import (
	"net/http"

	"messenger/resources"
)

type GraphQLPlayground struct{}

func (g GraphQLPlayground) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !g.isSupport(r.Method) {
		w.Write([]byte("only GET requests are supported"))
		return
	}

	playground, err := resources.Asset("playground.html")
	if err != nil {
		w.Write([]byte(""))
		return
	}
	w.Write(playground)
}

func (g GraphQLPlayground) isSupport(method string) bool {
	return (method == "GET")
}

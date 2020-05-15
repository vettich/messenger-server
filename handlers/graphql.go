package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"messenger/models"

	"github.com/graph-gophers/graphql-go"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type GraphQL struct {
	Schema  *graphql.Schema
	Session *r.Session
}

type requestParams struct {
	isMultiple bool
	params     []requestParam
}

func newRequestParams() *requestParams {
	return &requestParams{
		isMultiple: false,
		params:     make([]requestParam, 0),
	}
}

type requestParam struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

func (g GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !g.preAction(w, r) {
		return
	}

	tokenValue := g.parseToken(r)
	token, _ := models.GetTokenByValue(g.Session, tokenValue)
	ctx := r.Context()
	ctx = context.WithValue(ctx, "tokenValue", tokenValue)
	ctx = context.WithValue(ctx, "token", token)

	params, err := g.parseParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var responses []*graphql.Response
	for _, param := range params.params {
		resExec := g.Schema.Exec(ctx, param.Query, param.OperationName, param.Variables)
		responses = append(responses, resExec)
	}
	var responseJSON []byte
	if params.isMultiple {
		responseJSON, err = json.Marshal(responses)
	} else if len(responses) > 0 {
		responseJSON, err = json.Marshal(responses[0])
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(responseJSON)
}

func (g GraphQL) preAction(w http.ResponseWriter, r *http.Request) bool {
	method := r.Method
	if method != "GET" && method != "POST" && method != "OPTIONS" {
		w.Write([]byte("only GET, POST and OPTIONS requests are supported"))
		return false
	}

	if method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		return false
	}

	return true
}

func (g GraphQL) parseToken(r *http.Request) string {
	token := r.URL.Query().Get("token")
	if token == "" {
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Bearer" {
			return ""
		}
		token = auth[1]
	}
	return token
}

// TODO: implement parsing variables in GET and POST parameters
func (g GraphQL) parseParams(r *http.Request) (*requestParams, error) {
	contentType := r.Header.Get("Content-Type")
	if strings.Index(contentType, "application/json") != -1 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		var params requestParams
		var param requestParam
		err = json.Unmarshal(body, &param)
		if err != nil {
			err = json.Unmarshal(body, &params.params)
			if err != nil {
				return nil, err
			}
			params.isMultiple = true
		} else {
			params.params = append(params.params, param)
		}
		return &params, err
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return nil, err
		}
		param := requestParam{
			Query:         r.Form.Get("query"),
			OperationName: r.Form.Get("operationName"),
		}
		json.Unmarshal([]byte(r.Form.Get("variables")), &param.Variables)
		return &requestParams{params: []requestParam{param}}, nil
	}

	if r.Method != "GET" {
		return nil, errors.New("Request method must be GET or POST")
	}
	query := r.URL.Query()
	param := requestParam{
		Query:         query.Get("query"),
		OperationName: query.Get("operationName"),
	}
	json.Unmarshal([]byte(query.Get("variables")), &param.Variables)
	return &requestParams{params: []requestParam{param}}, nil
}

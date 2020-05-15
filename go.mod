module messenger

go 1.14

replace github.com/eviot/log => ../../eviot/log

require (
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/eviot/log v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/graph-gophers/graphql-go v0.0.0-20200309224638-dae41bde9ef9
	github.com/graph-gophers/graphql-transport-ws v0.0.0-20190611222414-40c048432299
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	gopkg.in/rethinkdb/rethinkdb-go.v6 v6.2.1
)

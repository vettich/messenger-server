package models

import (
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func fetch(cur *r.Cursor, to interface{}) error {
	err := cur.One(to)
	cur.Close()
	return err
}

func fetchAll(cur *r.Cursor, to interface{}) error {
	err := cur.All(to)
	cur.Close()
	return err
}

package resolvers

import (
	"errors"

	"github.com/eviot/log"
	"github.com/graph-gophers/graphql-go"
)

const EmptyID = graphql.ID("")

// Logging inside error and return new outside error
// The function is made to reduce the number of repeated lines
func writeError(inside error, outside string) error {
	log.Info(inside)
	return errors.New(outside)
}

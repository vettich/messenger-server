package resolvers

import (
	"context"

	"messenger/models"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type RootResolver struct {
	session *r.Session
}

func NewRootResolver(s *r.Session) *RootResolver {
	return &RootResolver{s}
}

func (self *RootResolver) tokenFromCtx(ctx context.Context) *models.Token {
	if self == nil {
		return nil
	}

	tokenValue, ok := ctx.Value("tokenValue").(string)
	if !ok {
		return nil
	}

	token, err := models.GetTokenByValue(self.session, tokenValue)
	if err != nil {
		return nil
	}
	return token
}

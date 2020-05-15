package resolvers

import (
	"context"
	"messenger/models"

	"github.com/graph-gophers/graphql-go"
)

type TokenResolver struct {
	*RootResolver
	token *models.Token
}

func (self *TokenResolver) ID() graphql.ID {
	if self == nil || self.token == nil {
		return EmptyID
	}
	return graphql.ID(self.token.ID)
}

func (self *TokenResolver) Value() string {
	if self == nil || self.token == nil {
		return ""
	}
	return self.token.Value
}

func (self *TokenResolver) User(ctx context.Context) *UserResolver {
	if self == nil || self.token == nil {
		return nil
	}

	user := self.token.User()
	if user == nil {
		return nil
	}

	return &UserResolver{self.RootResolver, user}
}

func (self *TokenResolver) UserID() graphql.ID {
	if self == nil || self.token == nil {
		return EmptyID
	}
	return graphql.ID(self.token.UserID)
}

func (self *TokenResolver) Active() bool {
	if self == nil || self.token == nil {
		return false
	}
	return self.token.Active
}

func (self *TokenResolver) CreatedAt() graphql.Time {
	if self == nil || self.token == nil {
		return graphql.Time{}
	}
	return graphql.Time{Time: self.token.CreatedAt}
}

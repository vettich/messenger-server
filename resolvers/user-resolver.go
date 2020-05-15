package resolvers

import (
	"messenger/models"

	"github.com/graph-gophers/graphql-go"
)

type UserResolver struct {
	*RootResolver
	user *models.User
}

func (self *UserResolver) ID() graphql.ID {
	if self == nil || self.user == nil {
		return EmptyID
	}
	return graphql.ID(self.user.ID)
}

func (self *UserResolver) Username() string {
	if self == nil || self.user == nil {
		return ""
	}
	return self.user.Username
}

func (self *UserResolver) Blocked() bool {
	if self == nil || self.user == nil {
		return false
	}
	return self.user.Blocked
}

func (self *UserResolver) CreatedAt() graphql.Time {
	if self == nil || self.user == nil {
		return graphql.Time{}
	}
	return graphql.Time{self.user.CreatedAt}
}

func (self *UserResolver) Chats() *[]*ChatResolver {
	if self == nil || self.user == nil {
		return nil
	}
	chats, err := models.ListChats(self.session, self.user.ID)
	if err != nil {
		return nil
	}
	resolvers := make([]*ChatResolver, len(chats))
	for i, chat := range chats {
		resolvers[i] = &ChatResolver{self.RootResolver, chat}
	}
	return &resolvers
}

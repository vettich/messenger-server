package resolvers

import (
	"messenger/models"

	"github.com/graph-gophers/graphql-go"
)

type MessageResolver struct {
	*RootResolver
	message *models.Message
}

func (self *MessageResolver) ID() graphql.ID {
	if self == nil || self.message == nil {
		return EmptyID
	}
	return graphql.ID(self.message.ID)
}

func (self *MessageResolver) Text() string {
	if self == nil || self.message == nil {
		return ""
	}
	return self.message.Text
}

func (self *MessageResolver) ChatID() string {
	if self == nil || self.message == nil {
		return ""
	}
	return self.message.ChatID
}

func (self *MessageResolver) UserID() string {
	if self == nil || self.message == nil {
		return ""
	}
	return self.message.UserID
}

func (self *MessageResolver) User() *UserResolver {
	if self == nil || self.message == nil {
		return nil
	}
	user, err := models.GetUser(self.session, self.message.UserID)
	if err != nil {
		return nil
	}
	return &UserResolver{self.RootResolver, user}
}

func (self *MessageResolver) CreatedAt() graphql.Time {
	if self == nil || self.message == nil {
		return graphql.Time{}
	}
	return graphql.Time{self.message.CreatedAt}
}

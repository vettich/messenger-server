package resolvers

import (
	"context"
	"messenger/models"

	"github.com/graph-gophers/graphql-go"
)

type ChatResolver struct {
	*RootResolver
	chat  *models.Chat
	token *models.Token
}

func (self *ChatResolver) ID() graphql.ID {
	if self == nil || self.chat == nil {
		return EmptyID
	}
	return graphql.ID(self.chat.ID)
}

func (self *ChatResolver) Name(ctx context.Context) string {
	if self == nil || self.chat == nil {
		return ""
	}

	if self.chat.IsPersonal {
		token := self.token
		if self.token == nil {
			var ok bool
			if token, ok = ctx.Value("token").(*models.Token); !ok {
				return ""
			}
		}

		self.chat.SetSession(self.session)
		users, err := self.chat.Users(token.UserID)
		if err != nil || len(users) == 0 {
			return ""
		}
		return users[0].Username
	}

	return self.chat.Name
}

func (self *ChatResolver) UserIDs() []string {
	if self == nil || self.chat == nil {
		return nil
	}
	return self.chat.UserIDs
}

func (self *ChatResolver) IsPersonal() bool {
	if self == nil || self.chat == nil {
		return false
	}
	return self.chat.IsPersonal
}

func (self *ChatResolver) CreatedAt() graphql.Time {
	if self == nil || self.chat == nil {
		return graphql.Time{}
	}
	return graphql.Time{self.chat.CreatedAt}
}

func (self *ChatResolver) Preview() *MessageResolver {
	if self == nil || self.chat == nil {
		return nil
	}

	msg, err := models.LastMessage(self.session, self.chat.ID)
	if err != nil {
		return nil
	}

	return &MessageResolver{self.RootResolver, msg}
}

func (self *ChatResolver) Messages() *[]*MessageResolver {
	if self == nil || self.chat == nil {
		return nil
	}

	messages, err := models.ListMessages(self.session, self.chat.ID)
	if err != nil {
		return nil
	}

	resolvers := make([]*MessageResolver, len(messages))
	for i, message := range messages {
		resolvers[i] = &MessageResolver{self.RootResolver, message}
	}

	return &resolvers
}

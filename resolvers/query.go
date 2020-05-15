package resolvers

import (
	"context"
	"errors"

	"messenger/models"

	"github.com/eviot/log"
)

func (r *RootResolver) User(ctx context.Context) (*UserResolver, error) {
	token, ok := ctx.Value("token").(*models.Token)
	if !ok || token == nil {
		return nil, errors.New("token not found")
	}
	user := token.User()
	if user == nil {
		return nil, errors.New("user not found")
	}
	return &UserResolver{r, user}, nil
}

func (r *RootResolver) UserByUsername(
	ctx context.Context,
	args struct {
		Username string
	},
) (*UserResolver, error) {
	_, ok := ctx.Value("token").(*models.Token)
	if !ok {
		return nil, errors.New("token not found")
	}
	user, err := models.GetUserByUsername(r.session, args.Username)
	if log.Debug(err) {
		return nil, errors.New("user not found")
	}
	return &UserResolver{r, user}, nil
}

func (r *RootResolver) Chats(ctx context.Context) (*[]*ChatResolver, error) {
	token, ok := ctx.Value("token").(*models.Token)
	if !ok {
		return nil, errors.New("token not found")
	}
	chats, err := models.ListChats(r.session, token.UserID)
	if err != nil {
		return nil, errors.New("chats not found")
	}
	resolvers := make([]*ChatResolver, len(chats))
	for i, chat := range chats {
		resolvers[i] = &ChatResolver{r, chat, nil}
	}
	return &resolvers, nil
}

func (r *RootResolver) Messages(
	ctx context.Context,
	args struct {
		ChatID string
	},
) (*[]*MessageResolver, error) {
	_, ok := ctx.Value("token").(*models.Token)
	if !ok {
		return nil, errors.New("token not found")
	}
	messages, err := models.ListMessages(r.session, args.ChatID)
	if err != nil {
		return nil, errors.New("chats not found")
	}
	resolvers := make([]*MessageResolver, len(messages))
	for i, message := range messages {
		resolvers[i] = &MessageResolver{r, message}
	}
	return &resolvers, nil
}

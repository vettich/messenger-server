package resolvers

import (
	"context"
	"errors"
	"messenger/models"
)

func (r *RootResolver) Login(
	ctx context.Context,
	args struct {
		Username string
		Password string
	},
) (*TokenResolver, error) {
	user, err := models.GetUserByUsername(r.session, args.Username)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if !user.CheckPassword(args.Password) {
		return nil, errors.New("wrong password")
	}

	token, err := user.CreateToken()
	if err != nil {
		return nil, errors.New("fail create token")
	}

	return &TokenResolver{r, token}, nil
}

func (r *RootResolver) Signup(
	ctx context.Context,
	args struct {
		Username string
		Password string
	},
) (*TokenResolver, error) {
	if models.UserByUsernameExists(r.session, args.Username) {
		return nil, errors.New("user exists")
	}

	user, err := models.CreateUser(r.session, args.Username, args.Password)
	if err != nil {
		return nil, errors.New("fail create user")
	}

	token, err := user.CreateToken()
	if err != nil {
		return nil, errors.New("fail create token")
	}
	return &TokenResolver{r, token}, nil
}

func (r *RootResolver) ChangeUserPassword(
	ctx context.Context,
	args struct {
		Old string
		New string
	},
) (bool, error) {
	token, ok := ctx.Value("token").(*models.Token)
	if !ok {
		return false, errors.New("token not found")
	}

	user := token.User()
	if user == nil {
		return false, errors.New("user not found")
	}

	if !user.CheckPassword(args.Old) {
		return false, errors.New("wrong old password")
	}

	user.ChangePassword(args.New)

	return true, nil
}

func (r *RootResolver) CreatePersonalChat(
	ctx context.Context,
	args struct {
		Username string
	},
) (*ChatResolver, error) {
	token, ok := ctx.Value("token").(*models.Token)
	if !ok {
		return nil, errors.New("token not found")
	}

	user, err := models.GetUserByUsername(r.session, args.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userIDs := []string{token.UserID, user.ID}

	if models.IssetPersonalChat(r.session, userIDs) {
		return nil, errors.New("chat already exists")
	}

	chat, err := models.CreatePersonalChat(r.session, userIDs)
	if err != nil {
		return nil, errors.New("fail create chat")
	}

	return &ChatResolver{r, chat}, nil
}

func (r *RootResolver) CreateMessage(
	ctx context.Context,
	args struct {
		Input struct {
			Text   string
			ChatID string
		}
	},
) (*MessageResolver, error) {
	token, ok := ctx.Value("token").(*models.Token)
	if !ok {
		return nil, errors.New("token not found")
	}

	message, err := models.CreateMessage(r.session, args.Input.Text, args.Input.ChatID, token.UserID)
	if err != nil {
		return nil, errors.New("fail create message")
	}

	return &MessageResolver{r, message}, nil
}

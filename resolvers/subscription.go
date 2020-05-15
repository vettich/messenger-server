package resolvers

import (
	"context"
	"errors"
	"messenger/models"

	"github.com/eviot/log"
)

func (r *RootResolver) WatchChatList(
	ctx context.Context,
	args struct {
		Token string
	},
) (chan *ChatResolver, error) {
	token, err := models.GetTokenByValue(r.session, args.Token)
	if log.Debug(err) {
		return nil, errors.New("token not found")
	}

	chatCancel := make(chan bool)
	chatChan, err := models.WatchChatList(r.session, token.UserID, chatCancel)
	if err != nil {
		close(chatCancel)
		return nil, err
	}

	chatList, err := models.ListChats(r.session, token.UserID)
	if err != nil {
		chatList = []*models.Chat{}
	}

	chatIDs := []string{}
	for _, c := range chatList {
		chatIDs = append(chatIDs, c.ID)
	}

	msgCancel := make(chan bool)
	msgChan, _ := models.WatchNewMessages(r.session, chatIDs, msgCancel)

	c := make(chan *ChatResolver)
	go func() {
		done := false
		for !done {
			select {
			case chat := <-chatChan:
				c <- &ChatResolver{r, chat}
				chatIDs = append(chatIDs, chat.ID)
				msgCancel <- true
				msgChan, _ = models.WatchNewMessages(r.session, chatIDs, msgCancel)

			case msg := <-msgChan:
				chat, err := models.GetChat(r.session, msg.ChatID)
				if err == nil {
					c <- &ChatResolver{r, chat}
				}

			case <-ctx.Done():
				chatCancel <- true
				msgCancel <- true
				done = true
			}
		}
		close(chatCancel)
		close(msgCancel)
		close(c)
	}()

	return c, nil
}

func (r *RootResolver) WatchNewMessagesInChat(
	ctx context.Context,
	args struct {
		Token  string
		ChatID string
	},
) (chan *MessageResolver, error) {
	_, err := models.GetTokenByValue(r.session, args.Token)
	if log.Debug(err) {
		return nil, errors.New("token not found")
	}

	cancel := make(chan bool)
	msgChan, err := models.WatchNewMessagesInChat(r.session, args.ChatID, cancel)
	if err != nil {
		close(cancel)
		return nil, err
	}

	c := make(chan *MessageResolver)
	go func() {
		done := false
		for !done {
			select {
			case msg := <-msgChan:
				c <- &MessageResolver{r, msg}

			case <-ctx.Done():
				cancel <- true
				done = true
			}
		}
		close(cancel)
		close(c)
	}()

	return c, nil
}

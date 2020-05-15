package models

import (
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type MessageChangesDoc struct {
	New *Message `rethinkdb:"new_val"`
}

type Message struct {
	ID        string    `rethinkdb:"id,omitempty" json:"id"`
	Text      string    `rethinkdb:"text"         json:"text"`
	ChatID    string    `rethinkdb:"chat_id"      json:"chat_id"`
	UserID    string    `rethinkdb:"user_id"      json:"user_id"`
	CreatedAt time.Time `rethinkdb:"created_at"   json:"created_at"`

	session *r.Session
}

func (self *Message) SetSession(s *r.Session) {
	if self == nil {
		return
	}
	self.session = s
}

func CreateMessage(s *r.Session, text, chatID, userID string) (*Message, error) {
	msg := &Message{"", text, chatID, userID, time.Now(), nil}
	res, err := r.Table(Messages).Insert(msg).RunWrite(s)
	if err != nil {
		return nil, err
	}
	if res.Inserted == 0 {
		return nil, NotInsertedError
	}
	msg.ID = res.GeneratedKeys[0]
	return msg, nil
}

func ListMessages(s *r.Session, chatID string) ([]*Message, error) {
	cur, err := r.Table(Messages).Filter(r.Row.Field("chat_id").Eq(chatID)).OrderBy("created_at").Run(s)
	if err != nil {
		return nil, err
	}
	messages := []*Message{}
	return messages, fetchAll(cur, &messages)
}

func LastMessage(s *r.Session, chatID string) (*Message, error) {
	cur, err := r.Table(Messages).Filter(r.Row.Field("chat_id").Eq(chatID)).OrderBy(r.Desc("created_at")).Limit(1).Run(s)
	if err != nil {
		return nil, err
	}
	message := &Message{}
	return message, fetch(cur, message)
}

func WatchNewMessagesInChat(s *r.Session, chatID string, cancel chan bool) (chan *Message, error) {
	cur, err := r.Table(Messages).Filter(r.Row.Field("chat_id").Eq(chatID)).Changes().Run(s)
	if err != nil {
		return nil, err
	}
	return baseWatchMessages(cur, cancel)
}

func WatchNewMessages(s *r.Session, chatIDs []string, cancel chan bool) (chan *Message, error) {
	filter := func(doc r.Term) r.Term {
		return r.Expr(chatIDs).Contains(doc.Field("chat_id"))
	}
	cur, err := r.Table(Messages).Filter(filter).Changes().Run(s)
	if err != nil {
		return nil, err
	}
	return baseWatchMessages(cur, cancel)
}

func baseWatchMessages(cur *r.Cursor, cancel chan bool) (chan *Message, error) {
	messages := make(chan *Message)
	c := make(chan MessageChangesDoc)
	cur.Listen(c)

	go func() {
		for {
			select {
			case msg := <-c:
				messages <- msg.New
			case <-cancel:
				cur.Close()
				return
			}
		}
	}()

	return messages, nil
}

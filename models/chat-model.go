package models

import (
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type ChatChangesDoc struct {
	New *Chat `rethinkdb:"new_val"`
}

type Chat struct {
	ID         string    `rethinkdb:"id,omitempty" json:"id"`
	Name       string    `rethinkdb:"name,omitempty" json:"name"`
	UserIDs    []string  `rethinkdb:"user_ids"     json:"user_ids"`
	IsPersonal bool      `rethinkdb:"is_personal"  json:"is_personal"`
	CreatedAt  time.Time `rethinkdb:"created_at"   json:"created_at"`

	session *r.Session
}

func (self *Chat) SetSession(s *r.Session) {
	if self == nil {
		return
	}
	self.session = s
}

func (self *Chat) Messages() ([]*Message, error) {
	if self == nil {
		return nil, nil
	}
	return ListMessages(self.session, self.ID)
}

func (self *Chat) Users(excludeUserIDs ...string) ([]*User, error) {
	if self == nil {
		return nil, nil
	}
	users := []*User{}
	for _, userID := range self.UserIDs {
		isExclude := false
		for _, exclude := range excludeUserIDs {
			if userID == exclude {
				isExclude = true
				break
			}
		}
		if isExclude {
			continue
		}
		user, err := GetUser(self.session, userID)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func CreatePersonalChat(s *r.Session, userIDs []string) (*Chat, error) {
	doc := &Chat{"", "", userIDs, true, time.Now(), s}
	res, err := r.Table(Chats).Insert(doc).RunWrite(s)
	if err != nil {
		return nil, err
	}
	if res.Inserted == 0 {
		return nil, NotInsertedError
	}
	doc.ID = res.GeneratedKeys[0]
	return doc, nil
}

func IssetPersonalChat(s *r.Session, userIDs []string) bool {
	cur, err := r.Table(Chats).Filter(func(doc r.Term) r.Term {
		return doc.Field("is_personal").Eq(true).
			And(doc.Field("user_ids").Contains(userIDs[0])).
			And(doc.Field("user_ids").Contains(userIDs[1]))
	}).Count().Run(s)
	if err != nil {
		return false
	}

	var cnt int
	cur.One(&cnt)
	cur.Close()
	return cnt != 0
}

func GetChat(s *r.Session, chatID string) (*Chat, error) {
	cur, err := r.Table(Chats).Get(chatID).Run(s)
	if err != nil {
		return nil, err
	}
	doc := &Chat{session: s}
	return doc, fetch(cur, doc)
}

func ListChats(s *r.Session, userID string) ([]*Chat, error) {
	cur, err := r.Table(Chats).Filter(r.Row.Field("user_ids").Contains(userID)).Run(s)
	if err != nil {
		return nil, err
	}
	chats := []*Chat{}
	return chats, fetchAll(cur, &chats)
}

func WatchChatList(s *r.Session, userID string, cancel chan bool) (chan *Chat, error) {
	cur, err := r.Table(Chats).Filter(r.Row.Field("user_ids").Contains(userID)).Changes().Run(s)
	if err != nil {
		return nil, err
	}

	chats := make(chan *Chat)
	chatListener := make(chan ChatChangesDoc)
	cur.Listen(chatListener)

	go func() {
		for {
			select {
			case chat := <-chatListener:
				chats <- chat.New
			case <-cancel:
				cur.Close()
				return
			}
		}
	}()

	return chats, nil
}

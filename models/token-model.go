package models

import (
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type Token struct {
	ID          string    `rethinkdb:"id,omitempty"`
	Value       string    `rethinkdb:"value,omitempty"`
	UserID      string    `rethinkdb:"user_id,omitempty"`
	Active      bool      `rethinkdb:"active,omitempty"`
	DeactivedAt time.Time `rethinkdb:"deactived_at,omitempty"`
	CreatedAt   time.Time `rethinkdb:"created_at,omitempty"`
	UpdatedAt   time.Time `rethinkdb:"updated_at,omitempty"`

	session *r.Session
	user    *User
}

func (token *Token) Validate() bool {
	if token == nil {
		return false
	}
	u := token.User()
	if u == nil || u.Blocked {
		return false
	}
	return token.Active
}

func (token *Token) User() *User {
	if token == nil {
		return nil
	}
	if token.user == nil {
		token.user = GetUserWithoutErr(token.session, token.UserID)
	}
	return token.user
}

func CreateToken(s *r.Session, userID string) (*Token, error) {
	token := Token{
		Value:     TokenHash(userID),
		UserID:    userID,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		session:   s,
	}
	res, err := r.Table(Tokens).Insert(&token).RunWrite(s)
	if err != nil {
		return nil, err
	}

	if res.Inserted == 0 {
		return nil, NotInsertedError
	}
	id := res.GeneratedKeys[0]
	token.ID = id
	return &token, nil
}

func GetTokenByValue(s *r.Session, value string) (*Token, error) {
	cur, err := r.Table(Tokens).Filter(r.Row.Field("value").Eq(value)).Run(s)
	if err != nil {
		return nil, err
	}

	token := Token{}
	err = cur.One(&token)
	cur.Close()
	token.session = s
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func ValidateTokenValue(s *r.Session, value string) (bool, error) {
	cur, err := r.Table(Tokens).Filter(r.Row.Field("value").Eq(value)).Run(s)
	if err != nil {
		return false, err
	}
	token := Token{}
	cur.One(&token)
	cur.Close()
	token.session = s
	return token.Validate(), nil
}

func DeactivateToken(s *r.Session, id string) error {
	token := map[string]interface{}{
		"active":      false,
		"deactivedAt": time.Now(),
	}
	_, err := r.Table(Tokens).Get(id).Update(token).Run(s)
	return err
}

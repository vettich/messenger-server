package models

import (
	"time"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// User is document of user
type User struct {
	ID           string    `rethinkdb:"id,omitempty" json:"id,omitempty"`
	Username     string    `rethinkdb:"username" json:"username,omitempty"`
	PasswordHash string    `rethinkdb:"password_hash" json:"-"`
	Blocked      bool      `rethinkdb:"blocked" json:"blocked"`
	CreatedAt    time.Time `rethinkdb:"created_at" json:"created_at,omitempty"`
	UpdatedAt    time.Time `rethinkdb:"updated_at" json:"updated_at,omitempty"`

	session *r.Session
}

func (self *User) CheckPassword(pass string) bool {
	if self == nil {
		return false
	}
	return CheckPasswordHash(pass, self.PasswordHash)
}

func (self *User) ChangePassword(newpass string) error {
	if self == nil {
		return nil
	}
	return SetNewUserPassword(self.session, self.ID, newpass)
}

func (self *User) CreateToken() (*Token, error) {
	if self == nil {
		return nil, nil
	}
	return CreateToken(self.session, self.ID)
}

func (self *User) Chats() ([]*Chat, error) {
	if self == nil {
		return nil, nil
	}
	return ListChats(self.session, self.ID)
}

func UserByUsernameExists(s *r.Session, username string) bool {
	cur, err := r.Table(Users).Filter(r.Row.Field("username").Eq(username)).Count().Run(s)
	if err != nil {
		return false
	}

	var exists bool
	cur.One(&exists)
	cur.Close()
	return exists
}

func CreateUser(s *r.Session, uname, pass string) (*User, error) {
	passHash, err := HashPassword(pass)
	if err != nil {
		return nil, err
	}
	user := User{
		Username:     uname,
		PasswordHash: passHash,
		Blocked:      false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		session:      s,
	}
	res, err := r.Table(Users).Insert(&user).RunWrite(s)
	if err != nil {
		return nil, err
	}

	if res.Inserted == 0 {
		return nil, NotInsertedError
	}
	id := res.GeneratedKeys[0]
	user.ID = id
	return &user, err
}

func SetNewUserPassword(s *r.Session, id, pass string) error {
	passHash, err := HashPassword(pass)
	if err != nil {
		return err
	}
	upd := map[string]interface{}{
		"password_hash": passHash,
	}
	_, err = r.Table(Users).Get(id).Update(upd).RunWrite(s)
	return err
}

func GetUser(s *r.Session, id string) (*User, error) {
	cur, err := r.Table(Users).Get(id).Run(s)
	if err != nil {
		return nil, err
	}

	user := User{}
	err = cur.One(&user)
	cur.Close()
	user.session = s
	return &user, err
}

func GetUserWithoutErr(s *r.Session, id string) *User {
	u, _ := GetUser(s, id)
	return u
}

func GetUserByUsername(s *r.Session, username string) (*User, error) {
	cur, err := r.Table(Users).Filter(r.Row.Field("username").Eq(username)).Run(s)
	if err != nil {
		return nil, err
	}

	user := User{}
	err = cur.One(&user)
	cur.Close()
	user.session = s
	return &user, err
}

func DeleteUser(s *r.Session, id string) error {
	_, err := GetUser(s, id)
	if err != nil {
		return err
	}
	upd := map[string]interface{}{
		"deleted":    true,
		"deleted_at": time.Now(),
	}
	_, err = r.Table(Users).Get(id).Update(upd).RunWrite(s)
	return err
}

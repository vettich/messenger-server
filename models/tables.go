package models

var (
	Users    = "users"
	Tokens   = "tokens"
	Chats    = "chats"
	Messages = "messages"
)

func Tables() []string {
	tables := []string{
		Users,
		Tokens,
		Chats,
		Messages,
	}
	return tables
}

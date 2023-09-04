package model

type Token struct {
	ID     string
	UserID string
	AppID  string
	Token  string
}

type TokenUpdate struct {
	ID    string
	OldID string
	Token string
}

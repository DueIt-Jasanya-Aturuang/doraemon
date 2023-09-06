package model

type Token struct {
	ID           int
	UserID       string
	AppID        string
	RememberMe   bool
	AcceesToken  string
	RefreshToken string
}

package _usecase

import (
	"errors"
)

var InvalidAppID = errors.New("invalid app id")
var InvalidEmailOrUsernameOrPassword = errors.New("invalid email or username or password")
var EmailIsExist = errors.New("email sudah terdaftar")
var UsernameIsExist = errors.New("username sudah terdaftar")
var InvalidTokenOauth = errors.New("invalid token")

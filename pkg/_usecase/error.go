package _usecase

import (
	"errors"
)

var InvalidAppID = errors.New("invalid app id")
var InvalidEmailOrUsernameOrPassword = errors.New("invalid email or username or password")
var EmailIsExist = errors.New("email sudah terdaftar")
var UsernameIsExist = errors.New("username sudah terdaftar")
var InvalidTokenOauth = errors.New("invalid token")
var InvalidEmail = errors.New("invalid email")
var InvalidUserID = errors.New("invalid user id")
var EmailIsActivited = errors.New("email anda sudah teraktivasi silahkan login")
var InvalidEmailOrOTP = errors.New("invalid email atau otp")
var InvalidToken = errors.New("invalid token")
var TokenExpired = errors.New("access token expired")
var JwtUserIDAndHeaderUserIDNotMatch = errors.New("invalid token")
var JwtAppIDAndHeaderAppIDNotMatch = errors.New("invalid token")
var UserIsNotActivited = errors.New("email user belum di activasi")
var InvalidOldPassword = errors.New("password lama anda tidak benar")

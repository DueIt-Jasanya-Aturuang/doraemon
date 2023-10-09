package usecase

import (
	"errors"
)

var InvalidAppID = errors.New("invalid app id")
var InvalidEmailOrUsernameOrPassword = errors.New("invalid email or username or password")
var EmailIsExist = errors.New("email sudah terdaftar")
var UsernameIsExist = errors.New("username sudah terdaftar")
var InvalidTokenOauth = errors.New("invalid Token")
var InvalidEmail = errors.New("invalid email")
var InvalidUserID = errors.New("invalid User id")
var EmailIsActivited = errors.New("email anda sudah teraktivasi silahkan login")
var InvalidEmailOrOTP = errors.New("invalid email atau otp")
var InvalidToken = errors.New("invalid Token")
var TokenExpired = errors.New("access Token expired")
var JwtUserIDAndHeaderUserIDNotMatch = errors.New("invalid Token")
var JwtAppIDAndHeaderAppIDNotMatch = errors.New("invalid Token")
var UserIsNotActivited = errors.New("email User belum di activasi")
var InvalidOldPassword = errors.New("password lama anda tidak benar")

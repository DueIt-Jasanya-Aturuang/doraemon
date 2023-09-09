package usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/conv"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type AuthUsecaseImpl struct {
	userRepo     repository.UserSqlRepo
	accessRepo   repository.AccessSqlRepo
	accountApi   repository.AccountApiRepo
	securityRepo repository.SecuritySqlRepo
}

func NewAuthUsecaseImpl(
	userRepo repository.UserSqlRepo,
	accessRepo repository.AccessSqlRepo,
	accountApi repository.AccountApiRepo,
	securityRepo repository.SecuritySqlRepo,
) usecase.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:     userRepo,
		accessRepo:   accessRepo,
		accountApi:   accountApi,
		securityRepo: securityRepo,
	}
}

func (a *AuthUsecaseImpl) Login(
	ctx context.Context, req *dto.LoginReq,
) (userResp *dto.UserResp, profileResp *dto.ProfileResp, tokenResp *dto.JwtTokenResp, err error) {
	// OpenConn membuka koneksi dari userRepo dan defer untuk close connetion
	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.userRepo.CloseConn()

	// GetUserByEmailOrUsername untuk memvalidasi email atau username tersedia atau gak
	// jika tidak akan return badlogin dan itu isinya map error400 pesan error http
	user, err := a.userRepo.GetUserByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info().Msgf("user login tapi email or username nya tidak tersedia")
			return nil, nil, nil, _error.BadLogin()
		}
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// jika ini bukan login dari oauth maka check apakah password dan old password match
	if !req.Oauth2 {
		checkPassword := helper.BcryptPasswordCompare(req.Password, user.Password)
		if !checkPassword {
			log.Info().Msgf("password user tidak sama dengan yang ada di database | req %s | db %s", req.Password, user.Password)
			return nil, nil, nil, _error.BadLogin()
		}
	}

	// mengambil data profile dari account service
	profile, err := a.accountApi.GetProfileByUserID(user.ID)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusBadGateway)
	}

	// helper untuk generate refrest token dan access token
	rtat, err := helper.GenerateRTAT(user.ID, req.AppID, req.RememberMe)
	if err != nil {
		return nil, nil, nil, err
	}

	// StartTx untuk insert ke database dan defer jika untuk melakukan rollback atau commit
	// jika ada kesalah pada rollback dan commit maka akan error 500
	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
			userResp, profileResp, tokenResp = nil, nil, nil
		}
	}()

	// CreateToken memasukan refresh token dan access token ke database
	err = a.securityRepo.CreateToken(ctx, rtat)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// formatting email untuk kebutuhan fe
	// convert dari model user dan profile ke dalam response
	// convert dari rtat atau token model ke dalam token response
	emailFormat := util.EmailFormat(user.Email)
	userResp, profileResp = conv.LoginModelToResp(user, profile, emailFormat)
	tokenResp = &dto.JwtTokenResp{
		Token: rtat.AcceesToken,
	}

	return userResp, profileResp, tokenResp, nil
}

func (a *AuthUsecaseImpl) Register(
	ctx context.Context, req *dto.RegisterReq,
) (userResp *dto.UserResp, profileResp *dto.ProfileResp, tokenResp *dto.JwtTokenResp, err error) {
	// OpenConn membuka koneksi dari userRepo dan defer untuk close connetion
	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.userRepo.CloseConn()

	// check email apakah sudah terdaftar atau belum
	// jika sudah maka akan mengembalikan error400 dan itu isisnya map pesan error.
	exists, err := a.userRepo.CheckUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	if exists {
		log.Info().Msgf("user register tapi email sudah terdaftar")
		return nil, nil, nil, _error.BadExistField("email", "email has been registered")
	}

	// check username apakah sudah terdaftar atau belum
	// jika sudah maka akan mengembalikan error400 dan itu isisnya map pesan error.
	exists, err = a.userRepo.CheckUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	if exists {
		log.Info().Msgf("user register tapi username sudah terdaftar")
		return nil, nil, nil, _error.BadExistField("username", "username has been registered")
	}

	// hashing password menggunakan bcrypt
	// saat ini ini saya hanya tau kalau itu error ketika password lebih dari 70 an
	// dan jadikan req.password denga passwordhash
	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	req.Password = passwordHash

	// StartTx start transaction untuk insert data
	// defer untuk commit atau rollback, jika terjadi error maka akan return 500
	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
			userResp, profileResp, tokenResp = nil, nil, nil
		}
	}()

	// generate uuid string
	// convert request kedalam usermodel dan accessmodel untuk insert ke database
	id := uuid.NewV4().String()
	userConv, accessConv := conv.RegisterReqToModel(req, id)

	// CreateUser melakukan insert user ke db menggunakan hasil convert yang tadi
	err = a.userRepo.CreateUser(ctx, userConv)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// CreateAccess melakukan insert user ke db menggunakan hasil convert yang tadi
	_, err = a.accessRepo.CreateAccess(ctx, accessConv)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	log.Info().Msgf(req.AppID)
	// GenerateRTAT helper untuk generate refrest token dan access token
	rtat, err := helper.GenerateRTAT(userConv.ID, req.AppID, false)
	if err != nil {
		return nil, nil, nil, err
	}

	// CreateToken memasukan refresh token dan access token ke database
	err = a.securityRepo.CreateToken(ctx, rtat)
	if err != nil {
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// SerializeProfile convert data insert profile kedalam json
	// ini static untuk convert profile model ke dalam json
	profileJson, err := util.SerializeProfile(userConv.ID)
	if err != nil {
		log.Err(err).Msg("failed marshal profile req")
		return nil, nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// CreateProfile melakukan insert profile dan hit api account service untuk melakukan insert
	profile, err := a.accountApi.CreateProfile(profileJson)
	if err != nil {
		return nil, nil, nil, err
	}

	// formatting email untuk kebutuhan fe
	// convert dari model user dan profile ke dalam response
	// convert dari rtat atau token model ke dalam token response
	emailFormat := util.EmailFormat(userConv.Email)
	userResp, profileResp = conv.RegisterModelToResp(userConv, profile, emailFormat)
	tokenResp = &dto.JwtTokenResp{
		Token: rtat.AcceesToken,
	}

	return userResp, profileResp, tokenResp, nil
}

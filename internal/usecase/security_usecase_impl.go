package usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type SecurityUsecaseImpl struct {
	userRepo     repository.UserSqlRepo
	securityRepo repository.SecuritySqlRepo
}

func NewSecurityUsecaseImpl(
	userRepo repository.UserSqlRepo,
	securityRepo repository.SecuritySqlRepo,
) usecase.SecurityUsecase {
	return &SecurityUsecaseImpl{
		userRepo:     userRepo,
		securityRepo: securityRepo,
	}
}

func (s *SecurityUsecaseImpl) JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) (bool, error) {
	// claim jwt access token, jika expired maka akan return true menandakan harus generate ulang token
	// selain error itu akan return 401
	claims, err := helper.ClaimsJwtHS256(req.Authorization, config.AccessTokenKeyHS)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		log.Err(err).Msg("failed claim jwt token")
		return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	// claims sub, jika ga bisa di assertion maka return 401
	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("failed to assertion sub jwt ke string")
		return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	// check apakah sub tadi match dengan req userid nya
	// jika tidak maka akan return 401
	if req.UserId != userID {
		log.Warn().Msg("request user id dan header user id tidak sama")
		return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	// open connection dari security repo
	// dan defer untuk melakukan close connection
	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return false, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	// get token by access token request
	getToken, err := s.securityRepo.GetTokenByAT(ctx, req.Authorization)
	if err != nil {
		// jika error sql no rows maka akan di delete semua token berdasarkan userid untuk keamanan
		// dan akan return 401, selain error sql no rows akan return 500
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")
			// process delete all token by user id, jika terjadi error maka akan return 500
			// start trasaction for deleted data
			err = s.securityRepo.StartTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			if err != nil {
				return false, _error.ErrStringDefault(http.StatusInternalServerError)
			}

			if err := s.securityRepo.DeleteAllTokenByUserID(ctx, req.UserId); err != nil {
				return false, _error.ErrStringDefault(http.StatusInternalServerError)
			}
			errEndTx := s.securityRepo.EndTx(err)
			if errEndTx != nil {
				return false, _error.ErrStringDefault(http.StatusInternalServerError)
			}

			return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
		}
		return false, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// check apakah app id sama dengan id yang ada di database token
	// jika tidak maka akan return 403
	if req.AppId != getToken.AppID {
		log.Warn().Msg("app id in header and app id di token database tidak sama")
		return false, _error.ErrStringDefault(http.StatusForbidden)
	}

	// check apakah ini dari endpoint activasi account
	// jika tidak maka akan check apakah user aktif atau gak
	// if !strings.Contains(endpoint, "/activasi-account") {
	// 	activasi, err := s.userRepo.CheckActivasiUserByID(ctx, req.UserId)
	// 	if err != nil {
	// 		// jika data tidak tersedia, maka akan return 404 dan message bahwa akun tidak terdaftar
	// 		// selain error itu akan return 500
	// 		if errors.Is(err, sql.ErrNoRows) {
	// 			return false, _error.ErrString("akun anda tidak terdaftar, silahkan register", http.StatusNotFound)
	// 		}
	// 		return false, _error.ErrStringDefault(http.StatusInternalServerError)
	// 	}
	//
	// 	// check apakah user tidak aktif, jika tidak maka akan return 403, yang menandakan user belum aktivasi account
	// 	if !activasi {
	// 		log.Warn().Msg("user belum melakukan activasi tetapi mencoba request ke endpoint lain")
	// 		return false, _error.ErrString("akun anda tidak aktif, silahkan aktifkan akun anda", http.StatusForbidden)
	// 	}
	// }

	return false, nil
}

func (s *SecurityUsecaseImpl) JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (tokenResp *dto.JwtTokenResp, err error) {
	// open connection dari security repo
	// dan defer untuk melakukan close connection
	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	// get token by access token request
	getToken, err := s.securityRepo.GetTokenByAT(ctx, req.Authorization)
	if err != nil {
		// jika error sql no rows maka akan di delete semua token berdasarkan userid untuk keamanan
		// dan akan return 401, selain error sql no rows akan return 500
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")
			// process delete all token by user id, jika terjadi error maka akan return 500
			// start trasaction for deleted data
			err = s.securityRepo.StartTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			if err != nil {
				return nil, _error.ErrStringDefault(http.StatusInternalServerError)
			}

			if err := s.securityRepo.DeleteAllTokenByUserID(ctx, req.UserId); err != nil {
				return nil, _error.ErrStringDefault(http.StatusInternalServerError)
			}
			errEndTx := s.securityRepo.EndTx(err)
			if errEndTx != nil {
				return nil, _error.ErrStringDefault(http.StatusInternalServerError)
			}
			return nil, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// start transaction untuk melakukan insert dan delete token
	err = s.securityRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// claim refresh token tadi
	_, err = helper.ClaimsJwtHS256(getToken.RefreshToken, config.RefreshTokenKeyHS)
	if err != nil {
		// jika error baik itu expired atau yang lainnya maka akan melakukan delete token
		err = s.securityRepo.DeleteToken(ctx, getToken.ID, req.UserId)
		if err != nil {
			// jika terjadi error maka akan melakukan rollback
			// jika terjadi error pada rollback atau tidak error maka akan return 500
			errEndTx := s.userRepo.EndTx(err)
			if errEndTx != nil {
				err = _error.ErrStringDefault(http.StatusInternalServerError)
			}
			return nil, _error.ErrStringDefault(http.StatusInternalServerError)
		}

		// jika tidak error maka akan melakukan commit
		// jika terjadi error pada commit maka akan return 500
		// dan return 401 bahwa token sudah invalid
		errEndTx := s.userRepo.EndTx(nil)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}

		return nil, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	} else {
		// jika token valid maka akan init defer untuk rollback atau commit
		defer func() {
			errEndTx := s.userRepo.EndTx(err)
			if errEndTx != nil {
				err = _error.ErrStringDefault(http.StatusInternalServerError)
				tokenResp = nil
			}
		}()
	}

	rtat, err := helper.GenerateRTAT(req.UserId, getToken.AppID, getToken.RememberMe)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.securityRepo.UpdateToken(ctx, getToken.ID, rtat.RefreshToken, rtat.AcceesToken)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	tokenResp = &dto.JwtTokenResp{
		Token: rtat.AcceesToken,
	}

	return tokenResp, nil
}

func (s *SecurityUsecaseImpl) Logout(ctx context.Context, req *dto.LogoutReq) (err error) {
	// open connection from security repo
	// defer untuk close connection
	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	// get token by access token
	// jika return sqlnorows, maka akan langsung aja di logout, selain itu akan return internal server error
	getToken, err := s.securityRepo.GetTokenByAT(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk logout menggunakan token yang lama")
			return nil
		}
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// start transaaction dari securiy repo
	// defer untuk melakukan commit atau rollback jika terjadi error akan return 500
	err = s.securityRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := s.securityRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}()

	// delete token
	// jika terjadi error akan return 500
	err = s.securityRepo.DeleteToken(ctx, getToken.ID, req.UserID)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

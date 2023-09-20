package _usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	util2 "github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type OTPUsecaseImpl struct {
	userRepo repository.UserSqlRepo
	redis    *infra.RedisImpl
}

func NewOTPUsecaseImpl(
	userRepo repository.UserSqlRepo,
	redis *infra.RedisImpl,
) usecase.OTPUsecase {
	return &OTPUsecaseImpl{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (o *OTPUsecaseImpl) OTPGenerate(ctx context.Context, req *dto.OTPGenerateReq) (err error) {
	// OpenConn kita openconnection db dari userrepo
	// defer untuk close connection dari user repo
	err = o.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer o.userRepo.CloseConn()

	// check apakah type nya activasi atau buka
	if req.Type == "activasi-account" {
		// check apakah user sudah aktiv atau belum
		// kalau belum maka akan return err400 bahwa user sudah aktivasi
		exist, err := o.userRepo.CheckActivasiUserByID(ctx, req.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info().Msgf(req.UserID)
				return _error.ErrStringDefault(http.StatusNotFound)
			}
			return _error.ErrStringDefault(http.StatusInternalServerError)
		}

		if exist {
			log.Warn().Msg("user sudah melakukan activasi tetapi malah ngirim minta code activasi lgi")
			return _error.Err400(map[string][]string{
				"email": {
					"permintaan anda tidak dapat di proses, email anda sudah di aktivasi silahkan login",
				},
			})
		}
	}

	// check apakah ada otp user di redis
	checkOtp, err := o.redis.Client.Exists(ctx, req.Type+":"+req.Email).Result()
	if err != nil {
		log.Err(err).Msg(_msg.LogErrExistsRedisClient)
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	// condition dimana kalo ada akan set expired nya saja
	// kalau tidak ada akan set ke dalam redis
	expOtp := 5 * time.Minute
	var otp string
	if checkOtp == 1 {
		// set expire redis
		err = o.redis.Client.Expire(ctx, req.Type+":"+req.Email, expOtp).Err()
		if err != nil {
			log.Err(err).Msg(_msg.LogErrExpireRedisClient)
			return _error.ErrStringDefault(http.StatusInternalServerError)
		}

		// get code otp in redis untuk ngirim ulang ke kafka
		result, err := o.redis.Client.Get(ctx, req.Type+":"+req.Email).Result()
		if err != nil {
			log.Err(err).Msg(_msg.LogErrGetRedisClient)
			return _error.ErrStringDefault(http.StatusInternalServerError)
		}
		otp = result

	} else {
		// generate code otp dengan panjang 6
		// dan set data in redis
		otp, err = util2.RandomChar(6)
		if err != nil {
			return _error.ErrStringDefault(http.StatusInternalServerError)
		}

		err = o.redis.Client.Set(ctx, req.Type+":"+req.Email, otp, expOtp).Err()
		if err != nil {
			log.Err(err).Msg(_msg.LogErrSetRedisClient)
			return _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}

	// kita set configurasi writer kafka
	// mengggunakan function KafkaWriter didalam helper
	// defer untuk close kafka writer
	w := helper.KafkaWriter()
	defer func() {
		errCloseKafka := w.Close()
		if errCloseKafka != nil {
			log.Err(errCloseKafka).Msg("failed close kafka connection")
		}
	}()

	// message untuk ngirim data kedalam kafka
	// kita menggunakan util untuk serialize ke dalam bentu byte
	// dan melakukan publish ke kafka WriteMessages
	msg, err := util2.SerializeMsgKafka(otp, req.Email, req.Type)
	log.Info().Msgf("otp : %s", otp)
	if err != nil {
		log.Err(err).Msg("failed marshal msg activasi account")
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	err = w.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		log.Err(err).Msg("failed write message to kafka")
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

func (o *OTPUsecaseImpl) OTPValidation(ctx context.Context, req *dto.OTPValidationReq) (err error) {
	// get data in redis apakah tersedia atau gak
	getOtp, err := o.redis.Client.Get(ctx, req.Type+":"+req.Email).Result()
	if err != nil {
		log.Err(err).Msg(_msg.LogErrGetRedisClient)
		return _error.Err400(map[string][]string{
			"otp": {
				"invalid your otp or email",
			},
			"email": {
				"invalid your otp or email",
			},
		})
	}

	// validasi apakah request otp dan otp didalam redis match atau gak
	if getOtp != req.OTP {
		log.Warn().Msg("otp in redis and request not the same")
		return _error.Err400(map[string][]string{
			"otp": {
				"kode otp anda tidak valid",
			},
		})
	}

	// jika match maka akan melakukan delete otp di redis
	err = o.redis.Client.Del(ctx, req.Type+":"+req.Email).Err()
	if err != nil {
		log.Err(err).Msg(_msg.LogErrDelRedisClient)
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

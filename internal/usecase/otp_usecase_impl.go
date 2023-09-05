package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type OTPUsecaseImpl struct {
	userRepo repository.UserSqlRepo
	redis    *config.RedisImpl
}

func NewOTPUsecaseImpl(
	userRepo repository.UserSqlRepo,
	redis *config.RedisImpl,
) usecase.OTPUsecase {
	return &OTPUsecaseImpl{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (o *OTPUsecaseImpl) OTPGenerate(ctx context.Context, req *dto.OTPGenerateReq) (err error) {
	err = o.userRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer o.userRepo.CloseConn()

	user, err := o.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return _error.ErrStringDefault(http.StatusNotFound)
		}
	}

	if user.EmailVerifiedAt {
		return _error.Err400(map[string][]string{
			"email": {
				"permintaan anda tidak dapat di proses, email anda sudah di aktivasi silahkan login",
			},
		})
	}

	otp, err := util.RandomChar(6)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	expOtp := 5 * time.Minute

	checkOtp, err := o.redis.Client.Exists(ctx, req.Type+":"+req.Email).Result()
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if checkOtp == 1 {
		err = o.redis.Client.Expire(ctx, req.Type+":"+req.Email, expOtp).Err()
		if err != nil {
			return _error.ErrStringDefault(http.StatusInternalServerError)
		}
	} else {
		err = o.redis.Client.Set(ctx, req.Type+":"+req.Email, otp, expOtp).Err()
		if err != nil {
			return _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}

	msg := map[string]string{
		"value": otp,
		"to":    req.Email,
		"type":  req.Type,
	}

	kafkaMsg, err := json.Marshal(msg)
	if err != nil {
		log.Err(err).Msg("failed marshal msg activasi account")
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	w := &kafka.Writer{
		Addr:  kafka.TCP(config.KafkaBroker),
		Topic: config.KafkaTopic,
	}
	defer func() {
		errCloseKafka := w.Close()
		if errCloseKafka != nil {
			log.Err(errCloseKafka).Msg("failed close kafka connection")
		}
	}()

	err = w.WriteMessages(ctx, kafka.Message{
		Value: kafkaMsg,
	})
	if err != nil {
		log.Err(err).Msg("failed write message to kafka")
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}

func (o *OTPUsecaseImpl) OTPValidation(ctx context.Context, req *dto.OTPValidationReq) (err error) {
	getOtp, err := o.redis.Client.Get(ctx, req.Type+":"+req.Email).Result()
	if err != nil {
		return _error.ErrStringDefault(http.StatusNotFound)
	}

	if getOtp != req.OTP {
		return _error.Err400(map[string][]string{
			"otp": {
				"kode otp anda tidak valid",
			},
		})
	}

	return nil
}

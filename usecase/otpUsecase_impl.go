package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type OTPUsecaseImpl struct {
	userRepo domain.UserRepository
	redis    *infra.RedisImpl
}

func NewOTPUsecaseImpl(
	userRepo domain.UserRepository,
	redis *infra.RedisImpl,
) domain.OTPUsecase {
	return &OTPUsecaseImpl{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (o *OTPUsecaseImpl) Generate(ctx context.Context, req *domain.RequestGenerateOTP) error {
	if err := o.userRepo.OpenConn(ctx); err != nil {
		return err
	}
	defer o.userRepo.CloseConn()

	domain.GetUserByEmail = req.Email
	user, err := o.userRepo.Get(ctx, domain.GetUserByEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info().Msgf("invalid email, user generate otp tapi invalid user email nya")
			return InvalidEmail
		}
		return err
	}

	if req.Type == util.ActivasiAccount {
		if req.UserID != user.ID {
			return InvalidUserID
		}
		if user.EmailVerifiedAt {
			return EmailIsActivited
		}
	}

	checkOtp, err := o.redis.Client.Exists(ctx, req.Type+":"+req.Email).Result()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientExists, err)
		return err
	}

	expOtp := 5 * time.Minute
	var otp string

	if checkOtp == 1 {
		err = o.redis.Client.Expire(ctx, req.Type+":"+req.Email, expOtp).Err()
		if err != nil {
			log.Warn().Msgf(util.LogErrRedisClientExpire, err)
			return err
		}

		result, err := o.redis.Client.Get(ctx, req.Type+":"+req.Email).Result()
		if err != nil {
			log.Warn().Msgf(util.LogErrRedisClientGet, err)
			return err
		}
		otp = result

	} else {
		otp, err = util.RandomChar(6)
		if err != nil {
			return err
		}

		err = o.redis.Client.Set(ctx, req.Type+":"+req.Email, otp, expOtp).Err()
		if err != nil {
			log.Warn().Msgf(util.LogErrRedisClientSet, err)
			return err
		}
	}

	w := infra.KafkaWriter()
	defer func() {
		if errClose := w.Close(); errClose != nil {
			log.Warn().Msgf(util.LogErrKafkaWriterClose, errClose)
		}
	}()

	messageByte, err := helper.SerializeMsgKafka(otp, req.Email, req.Type)
	if err != nil {
		return err
	}

	if err = w.WriteMessages(ctx, kafka.Message{
		Value: messageByte,
	}); err != nil {
		log.Warn().Msgf(util.LogErrKafkaWriteMessage, err)
		return err
	}

	return nil
}

func (o *OTPUsecaseImpl) Validation(ctx context.Context, req *domain.RequestValidationOTP) error {
	getOtp, err := o.redis.Client.Get(ctx, req.Type+":"+req.Email).Result()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientGet, err)
		return InvalidEmailOrOTP
	}

	if getOtp != req.OTP {
		log.Info().Msg("otp in redis and request not the same")
		return InvalidEmailOrOTP
	}

	err = o.redis.Client.Del(ctx, req.Type+":"+req.Email).Err()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientDel, err)
		return err
	}

	return nil
}

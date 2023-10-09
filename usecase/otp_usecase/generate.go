package otp_usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (o *OTPUsecaseImpl) Generate(ctx context.Context, req *usecase.RequestGenerateOTP) error {
	repository.GetUserByEmail = req.Email
	user, err := o.userRepo.Get(ctx, repository.GetUserByEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info().Msgf("invalid email, user_repository generate otp_usecase tapi invalid user_repository email nya")
			return usecase.InvalidEmail
		}
		return err
	}

	if req.Type == util.ActivasiAccount {
		if req.UserID != user.ID {
			return usecase.InvalidUserID
		}
		if user.EmailVerifiedAt {
			return usecase.EmailIsActivited
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

	messageByte, err := usecase.SerializeMsgKafka(otp, req.Email, req.Type)
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

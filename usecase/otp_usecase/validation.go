package otp_usecase

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (o *OTPUsecaseImpl) Validation(ctx context.Context, req *usecase.RequestValidationOTP) error {
	getOtp, err := o.redis.Client.Get(ctx, req.Type+":"+req.Email).Result()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientGet, err)
		return usecase.InvalidEmailOrOTP
	}

	if getOtp != req.OTP {
		log.Info().Msg("otp in redis and request not the same")
		return usecase.InvalidEmailOrOTP
	}

	err = o.redis.Client.Del(ctx, req.Type+":"+req.Email).Err()
	if err != nil {
		log.Warn().Msgf(util.LogErrRedisClientDel, err)
		return err
	}

	return nil
}

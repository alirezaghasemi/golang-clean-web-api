package services

import (
	"fmt"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/constants"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/cache"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/service_errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type OtpService struct {
	logger      logging.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

type OtpDto struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *config.Config) *OtpService {
	logger := logging.NewLogger(cfg)
	redisClient := cache.GetRedis()

	return &OtpService{logger: logger, cfg: cfg, redisClient: redisClient}
}

func (s *OtpService) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)

	val := &OtpDto{
		Value: otp,
		Used:  false,
	}

	res, err := cache.Get[OtpDto](s.redisClient, key)

	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}

	//err = cache.Set[*OtpDto](s.redisClient, key, val, s.cfg.Otp.ExpireTime*time.Second)
	err = cache.Set(s.redisClient, key, val, s.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (s *OtpService) ValidateOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	res, err := cache.Get[OtpDto](s.redisClient, key)
	if err != nil {
		return err
	} else if res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpInvalid}
	} else if !res.Used && res.Value == otp {
		res.Used = true
		_ = cache.Set(s.redisClient, key, res, s.cfg.Otp.ExpireTime*time.Second)
	}
	return nil
}

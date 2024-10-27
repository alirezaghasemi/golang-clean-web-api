package services

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/dto"
	"github.com/alirezaghasemi/golang-clean-web-api/src/common"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/db"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
	"gorm.io/gorm"
)

type UserService struct {
	logger     logging.Logger
	cfg        *config.Config
	otpService *OtpService
	database   *gorm.DB
}

func NewUserService(cfg *config.Config) *UserService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &UserService{
		cfg:        cfg,
		database:   database,
		logger:     logger,
		otpService: NewOtpService(cfg),
	}
}

func (s *UserService) SendOtp(req *dto.GetOtpRequest) error {
	otp := common.GenerateOtp()
	err := s.otpService.SetOtp(req.MobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

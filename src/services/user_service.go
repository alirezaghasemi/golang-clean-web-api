package services

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/dto"
	"github.com/alirezaghasemi/golang-clean-web-api/src/common"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/constants"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/db"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/models"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/service_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	logger       logging.Logger
	cfg          *config.Config
	otpService   *OtpService
	tokenService *TokenService
	database     *gorm.DB
}

func NewUserService(cfg *config.Config) *UserService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &UserService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		otpService:   NewOtpService(cfg),
		tokenService: NewTokenService(cfg),
	}
}

// LoginByUsername Login by username
func (s *UserService) LoginByUsername(request *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	var user models.User
	err := s.database.Model(&models.User{}).Where("user_name = ?", request.Username).Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
		return tx.Preload("Role")
	}).First(&user).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, err
	}

	tdto := tokenDto{
		UserId:       user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		MobileNumber: user.MobileNumber,
	}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tdto.Roles = append(tdto.Roles, ur.Role.Name)
		}
	}

	token, err := s.tokenService.GenerateToken(&tdto)
	if err != nil {
		return nil, err
	}

	return token, nil

}

// RegisterByUsername Register by username
func (s *UserService) RegisterByUsername(request *dto.RegisterUserByUsernameRequest) error {
	user := models.User{
		UserName:  request.Username,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Enabled:   false,
		UserRoles: nil,
	}

	exists, err := s.existsByEmail(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
	}

	exists, err = s.existsByUsername(user.UserName)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}
	}

	password := []byte(request.Password)
	hashPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	user.Password = string(hashPassword)

	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return err
	}

	tx := s.database.Begin()

	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}

	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: user.Id}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}

	tx.Commit()

	return nil
}

// RegisterLoginByMobileNumber Register/login by mobile number
func (s *UserService) RegisterLoginByMobileNumber(request *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	err := s.otpService.ValidateOtp(request.MobileNumber, request.Otp)
	if err != nil {
		return nil, err
	}

	exists, err := s.existsByMobileNumber(request.MobileNumber)
	if err != nil {
		return nil, err
	}

	u := models.User{MobileNumber: request.MobileNumber, UserName: request.MobileNumber}

	if exists {
		var user models.User
		err = s.database.Model(&models.User{}).Where("user_name = ?", u.UserName).Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).First(&user).Error
		if err != nil {
			return nil, err
		}

		tdto := tokenDto{
			UserId:       user.Id,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			MobileNumber: user.MobileNumber,
		}

		if len(*user.UserRoles) > 0 {
			for _, ur := range *user.UserRoles {
				tdto.Roles = append(tdto.Roles, ur.Role.Name)
			}
		}

		token, err := s.tokenService.GenerateToken(&tdto)
		if err != nil {
			return nil, err
		}

		return token, nil

	} else {
		password := []byte(common.GeneratePassword())
		hashPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
			return nil, err
		}
		u.Password = string(hashPassword)

		roleId, err := s.getDefaultRole()
		if err != nil {
			s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
			return nil, err
		}

		tx := s.database.Begin()
		err = tx.Create(&u).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
			return nil, err
		}
		err = tx.Create(&models.UserRole{RoleId: roleId, UserId: u.Id}).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
			return nil, err
		}
		tx.Commit()

		var user models.User
		err = s.database.Model(&models.User{}).Where("user_name = ?", u.UserName).Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).First(&user).Error
		if err != nil {
			return nil, err
		}

		tdto := tokenDto{
			UserId:       user.Id,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Email:        user.Email,
			MobileNumber: user.MobileNumber,
		}

		if len(*user.UserRoles) > 0 {
			for _, ur := range *user.UserRoles {
				tdto.Roles = append(tdto.Roles, ur.Role.Name)
			}
		}

		token, err := s.tokenService.GenerateToken(&tdto)
		if err != nil {
			return nil, err
		}

		return token, nil

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

func (s *UserService) existsByUsername(username string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*) > 0").
		Where("user_name = ?", username).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) existsByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*) > 0").
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) existsByEmail(email string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) getDefaultRole() (roleId int, err error) {
	if err := s.database.Model(&models.Role{}).
		Select("id").
		Where("name = ?", constants.DefaultRoleName).
		First(&roleId).
		Error; err != nil {
		return 0, err
	}
	return roleId, nil
}

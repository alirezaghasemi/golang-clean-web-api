package services

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/dto"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/constants"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/service_errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenService struct {
	logger logging.Logger
	cfg    *config.Config
}

type tokenDto struct {
	UserId       int
	FirstName    string
	LastName     string
	Username     string
	MobileNumber string
	Email        string
	Roles        []string
}

func NewTokenService(cfg *config.Config) *TokenService {
	logger := logging.NewLogger(cfg)
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *TokenService) GenerateToken(token *tokenDto) (*dto.TokenDetail, error) {
	tokenDetail := &dto.TokenDetail{}
	tokenDetail.AccessTokenExpiresTime = time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	tokenDetail.RefreshTokenExpiresTime = time.Now().Add(s.cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix()

	accessTokenClaims := jwt.MapClaims{}

	accessTokenClaims[constants.UserIdKey] = token.UserId
	accessTokenClaims[constants.FirstNameKey] = token.FirstName
	accessTokenClaims[constants.LastNameKey] = token.LastName
	accessTokenClaims[constants.UsernameKey] = token.Username
	accessTokenClaims[constants.MobileNumberKey] = token.MobileNumber
	accessTokenClaims[constants.EmailKey] = token.Email
	accessTokenClaims[constants.RolesKey] = token.Roles
	accessTokenClaims[constants.AccessTokenExpireTimeKey] = tokenDetail.AccessTokenExpiresTime

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	var err error
	tokenDetail.AccessToken, err = accessToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["user_id"] = token.UserId
	refreshTokenClaims["exp"] = tokenDetail.RefreshTokenExpiresTime

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenDetail.RefreshToken, err = refreshToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}

func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	accessToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnexpectedError}
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}
	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for key, value := range claims {
			claimMap[key] = value
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimsNotFound}
}

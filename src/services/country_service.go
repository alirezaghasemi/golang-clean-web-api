package services

import (
	"database/sql"
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/dto"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/constants"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/db"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/models"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type CountryService struct {
	database *gorm.DB
	logger   logging.Logger
}

func NewCountryService(cfg *config.Config) *CountryService {
	return &CountryService{
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg),
	}
}

// CreateCountry Create
func (s *CountryService) CreateCountry(ctx context.Context, request *dto.CreateUpdateCountryRequest) (*dto.CountryResponse, error) {
	country := &models.Country{Name: request.Name}
	country.CreatedBy = int(ctx.Value(constants.UserIdKey).(float64))
	country.CreatedAt = time.Now().UTC()

	tx := s.database.WithContext(ctx).Begin()

	err := tx.Create(&country).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}

	tx.Commit()

	dtoResponse := &dto.CountryResponse{Name: country.Name, Id: country.Id}
	return dtoResponse, nil
}

// UpdateCountry Update
func (s *CountryService) UpdateCountry(ctx context.Context, id int, request *dto.CreateUpdateCountryRequest) (*dto.CountryResponse, error) {
	updateMap := map[string]interface{}{
		"name":        request.Name,
		"modified_by": &sql.NullInt64{Valid: true, Int64: int64(ctx.Value(constants.UserIdKey).(float64))},
		"modified_at": &sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}
	tx := s.database.WithContext(ctx).Begin()

	err := tx.Model(&models.Country{}).Where("id = ?", id).Updates(updateMap).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return nil, err
	}

	country := &models.Country{}

	err = tx.Model(&models.Country{}).Where("id = ? AND deleted_by is null", id).First(&country).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}

	tx.Commit()
	dtoResponse := &dto.CountryResponse{Name: request.Name, Id: country.Id}

	return dtoResponse, nil
}

// DeleteCountry Delete
func (s *CountryService) DeleteCountry(ctx context.Context, id int) error {
	deleteMap := map[string]interface{}{
		"id":         id,
		"deleted_by": &sql.NullInt64{Valid: true, Int64: int64(ctx.Value(constants.UserIdKey).(float64))},
		"deleted_at": &sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}
	tx := s.database.WithContext(ctx).Begin()
	if err := tx.Model(&models.Country{}).Where("id = ?", id).Updates(deleteMap).Error; err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Delete, err.Error(), nil)
		return err
	}
	tx.Commit()
	return nil
}

// GetCountry Get by id
func (s *CountryService) GetCountry(ctx context.Context, id int) (*dto.CountryResponse, error) {
	country := &models.Country{}
	if err := s.database.Model(&models.Country{}).Where("id = ?", id).First(&country).Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}

	dtoResponse := &dto.CountryResponse{Name: country.Name, Id: country.Id}
	return dtoResponse, nil
}

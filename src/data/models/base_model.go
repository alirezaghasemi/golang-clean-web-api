package models

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id int `gorm:"primary_key"`

	CreatedAt  time.Time    `gorm:"TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"TIMESTAMP with time zone;null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

func (bm *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}

	bm.CreatedAt = time.Now().UTC()
	bm.CreatedBy = userId
	return
}
func (bm *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: value.(int64)}
	}

	bm.ModifiedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	bm.ModifiedBy = userId
	return
}

func (bm *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: value.(int64)}
	}

	bm.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	bm.DeletedBy = userId
	return
}

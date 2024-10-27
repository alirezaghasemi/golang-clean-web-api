package models

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;size:15;unique;not null"`
	UserRoles *[]UserRole
}

package models

type UserRole struct {
	BaseModel
	UserId int
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	RoleId int
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

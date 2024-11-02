package models

type User struct {
	BaseModel
	UserName     string `gorm:"type:string;size:20;not null;unique"`
	FirstName    string `gorm:"type:string;size:15;null"`
	LastName     string `gorm:"type:string;size:15;null"`
	MobileNumber string `gorm:"type:string;size:11;unique;default:null"`
	Email        string `gorm:"type:string;size:64;unique;default:null"`
	Password     string `gorm:"type:string;size:64;not null"`
	Enabled      bool   `gorm:"type:bool;default:true"`
	UserRoles    *[]UserRole
}

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;size:15;unique;not null"`
	UserRoles *[]UserRole
}

type UserRole struct {
	BaseModel
	UserId int
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
	RoleId int
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION;"`
}

/*
user -> n roles
role -> n users

users n <-> n roles
user
user_role -> user_id,role_id,...
role
*/

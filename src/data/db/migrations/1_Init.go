package migrations

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/constants"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/db"
	"github.com/alirezaghasemi/golang-clean-web-api/src/data/models"
	"github.com/alirezaghasemi/golang-clean-web-api/src/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.GetConfig())

func Up_1() {
	database := db.GetDb()

	CreateTables(database)
	createDefaultInformation(database)
}

func CreateTables(database *gorm.DB) {
	var tables []interface{}

	country := models.Country{}
	city := models.City{}
	user := models.User{}
	role := models.Role{}
	userRole := models.UserRole{}

	tables = addNewTable(database, country, tables)
	tables = addNewTable(database, city, tables)
	tables = addNewTable(database, user, tables)
	tables = addNewTable(database, role, tables)
	tables = addNewTable(database, userRole, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Migration, "tables are not create", nil)
	}

	logger.Info(logging.Postgres, logging.Migration, "tables created", nil)
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createDefaultInformation(database *gorm.DB) {
	adminRole := models.Role{Name: constants.AdminRoleName}
	createRoleIfNotExists(database, &adminRole)
	createRoleIfNotExists(database, &models.Role{Name: constants.DefaultRoleName})

	user := models.User{
		UserName:     constants.DefaultUserName,
		FirstName:    "Alireza",
		LastName:     "Ghasemi",
		MobileNumber: "09178049681",
		Email:        "soyeiran@gmail.com",
	}
	pass := "123456"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	createAdminUserIfNotExists(database, &user, adminRole.Id)
}

func createRoleIfNotExists(database *gorm.DB, role *models.Role) {
	exists := 0
	database.Model(&models.Role{}).Select("1").Where("name = ?", role.Name).First(&exists)

	if exists == 0 {
		database.Create(role)
	}
}

func createAdminUserIfNotExists(database *gorm.DB, user *models.User, roleId int) {
	exists := 0
	database.Model(&models.Role{}).Select("1").Where("username = ?", user.UserName).First(&exists)

	if exists == 0 {
		database.Create(user)
		userRole := models.UserRole{UserId: user.Id, RoleId: roleId}
		database.Create(&userRole)
	}
}

func Down_1() {}

package services

import (
	"TradeIT/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

type userData struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Passwd string `json:"password"`
}

func (u *UserService) InitService(db *gorm.DB) {
	u.db = db
	u.db.AutoMigrate(&models.User{})
	u.db.AutoMigrate(&models.Credential{})
}

func (us *UserService) RegisterUserService(c *gin.Context) int {
	var temp userData
	if err := c.ShouldBindJSON(&temp); err != nil {
		fmt.Println(err)
	}
	id := createUserId()
	fmt.Println(temp)
	err := us.db.Create(&models.User{ //creates the user entry in user details
		Id:    id,
		Name:  temp.Name,
		Email: temp.Email,
		Funds: 10000.00,
	}).Error
	if err != nil {
		panic(err)
	}

	errr := us.db.Create(&models.Credential{ //creates the credential entry for authorisations
		Id:     id,
		Email:  temp.Email,
		Passwd: temp.Passwd,
		Token:  "",
	}).Error
	if errr != nil {
		panic(errr)
	}
	return id
}

package services

import (
	"TradeIT/middleware"
	"TradeIT/models"
	"fmt"
	"log"

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

func (u *UserService) SetDB(db *gorm.DB) {
	u.db = db
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

func (us *UserService) LoginUserService(c *gin.Context) string {
	var Cred middleware.LoginCred
	if err := c.ShouldBindJSON(&Cred); err != nil {
		log.Fatal("err in binding data to cred: ", err)
	}

	status, errr := middleware.LoginValidator(us.db, Cred.Email, Cred.Passwd)

	if status && errr == "" {
		return "login success"
	} else {
		return errr
	}
}

func (us *UserService) UpdateUserEmailService(c *gin.Context) string {
	var param middleware.UpdateEmailParameters
	if err := c.BindJSON(&param); err != nil {
		log.Fatal("Error in binding updateParameter json ", err)
	}

	return middleware.EmailUpdater(us.db, param)
}

func (us *UserService) UpdateUserPasswdService(c *gin.Context) string {
	var param middleware.UpdatePasswdParameters
	if err := c.BindJSON(&param); err != nil {
		log.Fatal("Error in binding updateParameter json ", err)
	}

	return middleware.PasswdUpdator(us.db, param)
}

package middleware

import (
	"TradeIT/models"
	"fmt"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type UpdateEmailParameters struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	Passwd string `json:"password"`
}

type LoginCred struct {
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
}

type UpdatePasswdParameters struct {
	Id        string `json:"id"`
	OldPasswd string `json:"password"`
	NewPasswd string `json:"newPassword"`
}

func LoginValidator(db *gorm.DB, email string, passwd string) (bool, string, string) {
	var Cred models.Credential
	err := db.Where("email= ?", email).Find(&Cred).Error
	if err != nil {
		log.Fatal("Error in cred from tables ", err)
	}

	if Cred.Email == email && Cred.Passwd == passwd {
		return true, CreateToken(Cred.Id), ""
	} else if Cred.Email == email && Cred.Passwd != passwd {
		return false, "", "check the credentials"
	} else {
		return false, "", "user doesn't exist "
	}
}

func PasswdValidator(db *gorm.DB, parameter UpdateEmailParameters) bool {
	var userCount models.Credential
	fmt.Println(parameter)
	id, _ := strconv.Atoi(parameter.Id)
	result := db.Where("id = ? AND passwd =?", id, parameter.Passwd).Find(&userCount)
	if result.Error != nil {
		log.Fatal("error in getting the creds from the db: ", result.Error)
	}
	if result.RowsAffected == 1 {
		fmt.Println(userCount)
		userCount.Email = parameter.Email
		fmt.Println(userCount)
		db.Save(&userCount)
		fmt.Println(userCount)
		return true
	} else {
		fmt.Printf("incorrect Passwd")
		return false
	}
}

func EmailUpdater(db *gorm.DB, parameter UpdateEmailParameters) string {
	if PasswdValidator(db, parameter) {
		var user models.User
		id, _ := strconv.Atoi(parameter.Id)
		if err := db.Where("Id = ?", id).Find(&user).Error; err != nil {
			log.Fatal("error in entry finding (mail updation): ", err)
		}
		user.Email = parameter.Email
		if err := db.Save(&user).Error; err != nil {
			log.Fatal("error in saving the updates to database: ", err)
		}
		return "value changed succesfully"
	} else {
		return "password incorrect"
	}
}

func PasswdUpdator(db *gorm.DB, parameters UpdatePasswdParameters) string {
	var authDetail models.Credential
	if err := db.Where("id = ? AND passwd = ?", parameters.Id, parameters.OldPasswd).Find(&authDetail).Error; err != nil {
		log.Fatal("Error in binding: ", err)
		return "Internal server error"
	}
	if authDetail.Passwd == parameters.OldPasswd {
		authDetail.Passwd = parameters.NewPasswd
		db.Save(&authDetail)
		return "successfully changed password"
	}
	return "incorrect Password"
}

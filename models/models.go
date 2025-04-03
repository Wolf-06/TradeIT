package models

type User struct {
	Id    int     `gorm:"PrimaryKey"`
	Name  string  `json:"name" validate:"required"`
	Email string  `json:"email" validate:"email,required"`
	Funds float32 `json:"funds"`
}

type Credential struct {
	Id     int    `gorm:"PrimaryKey"`
	Email  string `json:"email" validate:"email,  required"`
	Passwd string `json:"passwd" validate:"required, min=6"`
	Token  string
}

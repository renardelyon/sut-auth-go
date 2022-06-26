package model

type User struct {
	Id       string
	Username string
	Name     string
	Password string
	Role     string
	AdminId  string
	Admin    Admin `gorm:"foreignKey:AdminId"`
}

type Admin struct {
	Id       string
	Username string
	Name     string
	Password string
	Role     string
}

type Token struct {
	Id    string `gorm:"primaryKey"`
	Token string
}

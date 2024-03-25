package models

import (
	"errors"
	"github.com/doohee323/tz-s3-repo-mgmt/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
	"time"
)

type User struct {
	ID          uint   `gorm:"primary_key"`
	Username    string `gorm:"size:255;not null;unique" json:"username"`
	Password    string `gorm:"size:255;not null;" json:"password"`
	Name        string `gorm:"size:255;not null;" json:"name"`
	Email       string `gorm:"size:255;not null;" json:"email"`
	Description string `gorm:"size:255;" json:"description"`
	Role        string `gorm:"size:255;not null;" json:"role"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUsers() ([]User, error) {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		return []User{}, err
	}
	return users, nil
}

func GetUserByUsername(username string) (User, error) {
	var user User
	DB.Where("username=?", username).Find(&user)
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}
	user.PrepareGive()
	return user, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username, password string) (string, error) {
	var err error
	u := User{}
	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return "", err
	}
	err = VerifyPassword(password, u.Password)
	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := utils.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *User) SaveUser() (*User, error) {
	var err error
	_ = u.BeforeSave()
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil
}

func (u *User) UpdateUser() (*User, error) {
	if u.ID == 0 {
		return u, errors.New("User not found!")
	}
	_ = u.BeforeSave()
	err := DB.Model(u).Updates(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) DeleteUser() (*User, error) {
	if u.ID == 0 {
		return u, errors.New("User not found!")
	}
	err := DB.Model(u).Delete(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

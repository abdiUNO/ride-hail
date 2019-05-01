package models

import (
	"fmt"
	"os"
	"strings"
	"time"

	u "github.com/abdullahi/go-drive/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserType string

/*
JWT claims struct
*/
type Token struct {
	UserId   string
	IsDriver bool
	jwt.StandardClaims
}

type Driver struct {
	ID         string     `gorm:"primary_key;type:varchar(255);"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Online     bool       `json:"online";gorm:"default:false"`
	Location   Location   `json:"location"`
	LocationID string     `json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"-"`
}

func (driver *Driver) BeforeCreate(scope *gorm.Scope) error {
	u1 := uuid.Must(uuid.NewV4())
	scope.SetColumn("ID", u1.String())
	return nil
}

func (driver *Driver) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(driver.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(driver.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	temp := &Driver{}

	err := GetDB().Table("drivers").Where("email = ?", driver.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another driver."), false
	}

	error := GetDB().Table("drivers").Where("phone = ?", driver.Phone).First(temp).Error
	if error != nil && error != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	if temp.Phone != "" {
		return u.Message(false, "Phone number is already in use by another driver."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (driver *Driver) Create() map[string]interface{} {
	if resp, ok := driver.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(driver.Password), bcrypt.DefaultCost)
	driver.Password = string(hashedPassword)

	GetDB().Create(driver)

	tk := &Token{UserId: driver.ID, IsDriver: true}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	driver.Password = ""

	response := u.Message(true, "User has been created")
	response["driver"] = driver
	response["token"] = tokenString

	return response
}

func Login(email, password string) map[string]interface{} {

	driver := &Driver{}
	err := GetDB().Table("drivers").Preload("Location").Where("email = ?", email).First(driver).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(driver.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	driver.Password = ""

	tk := &Token{UserId: driver.ID, IsDriver: true}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	response := u.Message(true, "Logged In")
	response["driver"] = driver
	response["token"] = tokenString

	return response
}

func GetDriver(id string) *Driver {
	driver := &Driver{}
	GetDB().Table("drivers").Preload("Location").Where("id = ?", id).First(driver)

	if driver.Email == "" {
		return nil
	}

	driver.Password = ""
	return driver
}

func GetDrivers() []*Driver {

	drivers := make([]*Driver, 0)

	err := GetDB().Table("drivers").Preload("Location").Find(&drivers).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return drivers
}

func ChangeStatus(id string, status bool) map[string]interface{} {
	driver := &Driver{}
	GetDB().Table("drivers").Preload("Location").Where("id = ?", id).First(driver)

	if driver.Email == "" {
		return nil
	}

	driver.Online = status

	GetDB().Table("drivers").Update(driver)

	response := u.Message(true, "Updated status")
	response["driver"] = driver

	return response

}

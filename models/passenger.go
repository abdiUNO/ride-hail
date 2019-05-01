package models

import (
	"os"
	"strings"
	"time"

	u "github.com/abdullahi/go-drive/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Passenger struct {
	ID        string     `gorm:"primary_key;type:varchar(255);"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

func (passenger *Passenger) BeforeCreate(scope *gorm.Scope) error {
	u1 := uuid.Must(uuid.NewV4())
	scope.SetColumn("ID", u1.String())
	return nil
}

func (passenger *Passenger) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(passenger.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(passenger.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	temp := &Passenger{}

	err := GetDB().Table("passengers").Where("email = ?", passenger.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	error := GetDB().Table("passengers").Where("phone = ?", passenger.Phone).First(temp).Error
	if error != nil && error != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}

	if temp.Phone != "" {
		return u.Message(false, "Phone number is already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (passenger *Passenger) Create() map[string]interface{} {
	if resp, ok := passenger.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passenger.Password), bcrypt.DefaultCost)
	passenger.Password = string(hashedPassword)

	GetDB().Create(passenger)

	tk := &Token{UserId: passenger.ID, IsDriver: false}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	passenger.Password = ""

	response := u.Message(true, "User has been created")
	response["passenger"] = passenger
	response["token"] = tokenString

	return response
}

func PassengerLogin(email, password string) map[string]interface{} {

	passenger := &Passenger{}
	err := GetDB().Table("passengers").Where("email = ?", email).First(passenger).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passenger.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	passenger.Password = ""

	tk := &Token{UserId: passenger.ID, IsDriver: false}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	response := u.Message(true, "Logged In")
	response["passenger"] = passenger
	response["token"] = tokenString

	return response
}

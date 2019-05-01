package models

import (
	"fmt"
	"time"

	u "github.com/abdullahi/go-drive/utils"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type RideStatus string

const (
	Pending    RideStatus = "pending"
	Accepted   RideStatus = "accepted"
	Arrived    RideStatus = "arrived"
	DroppedOff RideStatus = "droppedoff"
	Cancelled  RideStatus = "cancelled"
	Failed     RideStatus = "failed"
)

func (e *RideStatus) Scan(value interface{}) error {
	*e = RideStatus(value.([]byte))
	return nil
}

func (e RideStatus) Value() (string, error) {
	return string(e), nil
}

type Ride struct {
	ID          string     `gorm:"primary_key;type:varchar(255);"`
	Driver      Driver     `json:"driver"`
	DriverID    string     `json:"-"`
	Passenger   Passenger  `json:"passenger"`
	PassengerID string     `json:"-"`
	PickUp      Location   `json:"pickup"`
	PickUpID    string     `json:"-"`
	DropOff     Location   `json:"dropoff"`
	DropOffID   string     `json:"-"`
	StartTime   time.Time  `json:"startTime"`
	EndTime     time.Time  `json:"endTime"`
	Status      RideStatus `json:"status" sql:"type:ENUM('pending','accepted','arrived','droppedoff', 'cancelled', 'failed')"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"-"`
}

func (ride *Ride) BeforeCreate(scope *gorm.Scope) error {
	u1 := uuid.Must(uuid.NewV4())
	scope.SetColumn("ID", u1.String())
	return nil
}

func (ride *Ride) Create(token *Token) map[string]interface{} {

	if token.IsDriver {
		return u.Message(false, "Driver unauthorized to create new trips")
	}

	temp := &Ride{}

	err := GetDB().Table("rides").Where("passenger_id = ? AND status = ?", token.UserId, Arrived).Or("passenger_id = ? AND status = ?", token.UserId, Arrived).Or("passenger_id = ? AND status = ?", token.UserId, Pending).First(&temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry")
	}

	fmt.Println(token.UserId)

	fmt.Println(temp.DriverID)

	if (token.IsDriver && temp.DriverID != "") || (token.IsDriver == false && temp.PassengerID != "") {
		return u.Message(false, "Ride is already in progress")
	}

	response := u.Message(true, "New Trip Created")
	ride.Status = Pending
	ride.StartTime = time.Now()
	ride.StartTime.Add(time.Minute)

	ride.EndTime = time.Now()
	ride.StartTime.Add(time.Minute * 5)

	GetDB().Create(&ride).Related(&ride.Passenger)

	response["ride"] = ride
	return response
}

func GetRide(driverId string) map[string]interface{} {
	ride := &Ride{}
	err := GetDB().Table("rides").Preload("Driver.Location").Preload("Passenger").Preload("PickUp").Preload("DropOff").Where("driver_id = ? AND status != ?", driverId, DroppedOff).First(&ride).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Could not find ride")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	response := u.Message(true, "Update ride status")
	response["ride"] = ride
	return response
}

func UpdateStatus(token Token, status string) map[string]interface{} {
	var err error
	if token.IsDriver {
		err = GetDB().Table("rides").Where("driver_id = ? AND (status != 'droppedoff' AND status != 'failed' AND status != 'cancelled')", token.UserId).Update("status", RideStatus(status)).Error
	} else {
		err = GetDB().Table("rides").Where("passenger_id = ? AND (status != 'droppedoff' AND status != 'failed' AND status != 'cancelled')", token.UserId).Update("status", RideStatus(status)).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Could not find ride")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	response := u.Message(true, "Update ride status")
	return response
}

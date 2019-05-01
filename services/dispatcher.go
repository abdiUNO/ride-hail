package services

import (
	"fmt"
	"math"

	"github.com/abdullahi/go-drive/models"
)

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return (2 * r * math.Asin(math.Sqrt(h))) * 0.000621
}

var FindDriver = func(lat float64, long float64) (*models.Driver, error) {
	var distance float64 = 10000.0
	drivers := models.GetDrivers()
	nearestDriver := &models.Driver{}

	for _, driver := range drivers {
		driverDistance := getDistance(lat, long, driver.Location.Latitude, driver.Location.Longitude)
		fmt.Println(driverDistance)
		if driverDistance < distance {
			nearestDriver = driver
		}
	}

	return nearestDriver, nil
}

package gogpslib

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func StrToInt(str string) int {
	i, err := strconv.Atoi(strings.Trim(str, " \t\n\r"))
	if err != nil {
		log.Printf("Failed to converting string to int: %s", err.Error())
		return -1
	}
	return i
}

func StrToFloat(str string) float64 {
	f, err := strconv.ParseFloat(strings.Trim(str, " \t\n\r"), 64)
	if err != nil {
		log.Printf("Failed to converting string to float: %s", err.Error())
		return -1.0
	}
	return f
}

func ParseDelphiDate(dt string) time.Time {
	dateSource := StrToFloat(dt)
	intpart, fracpart := math.Modf(dateSource)

	loc, _ := time.LoadLocation("UTC")
	delphiEpoch := time.Date(
		1899,
		time.December,
		30,
		0,
		0,
		0,
		0,
		loc)
	pointDate := delphiEpoch.AddDate(0, 0, int(intpart))
	return pointDate.Add(time.Second * time.Duration(fracpart*24*60*60))
}

func ToDelphiDate(t time.Time) float64 {
	return 0.0
}

package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	NairobiTZ      = "Africa/Nairobi"
	DateFormat     = "2006-01-02"
	DateFormat2    = "2006-Jan-02"
	DateTimeFormat = "2006-01-02 15:04:05"
)

var (
	loc *time.Location
)

func init() {
	var err error
	loc, err = time.LoadLocation(NairobiTZ)
	if err != nil {
		panic(err)
	}
}

func CurrentDate() string {
	return time.Now().In(loc).Format(DateFormat)
}

func CurrentDate2() string {
	return time.Now().In(loc).Format(DateFormat2)
}

func CurrentFullDate() string {
	return time.Now().In(loc).Format(DateTimeFormat)
}

func CurrentYear() string {
	return time.Now().In(loc).Format("2006")
}

func CurrentMonth() string {
	return time.Now().In(loc).Format("01")
}

func CurrentMonthName() string {
	return time.Now().In(loc).Format("Jan")
}

func CurrentDay() string {
	return time.Now().In(loc).Format("02")
}

func CurrentDayName() string {
	return time.Now().In(loc).Format("Mon")
}

// create a function takes in two params date and desired format date to specified format e.g YYYY-MM-DD and YYYY-MM-DD HH:MM:SS then return the formatted date. Ensure time is in Nairobi TZ

func FormatDate(inputDate string, desiredLength int) (string, error) {
	var layout string
	var outputLayout string

	switch desiredLength {
	case 19:
		layout = "2006-01-02T15:04:05Z" // For DATETIME
		outputLayout = "2006-01-02 15:04:05"
	case 10:
		layout = "2006-01-02" // For DATE
		outputLayout = "2006-01-02"
	default:
		return "", fmt.Errorf("invalid desired length")
	}

	// Parse the input date string to time.Time object
	date, err := time.Parse(layout, inputDate)
	if err != nil {
		return "", err
	}

	// Check if parsed date string has the expected length
	if len(date.Format(layout)) != desiredLength {
		return "", fmt.Errorf("invalid input date format")
	}

	// Format the date to the desired output format
	formattedDate := date.Format(outputLayout)

	return formattedDate, nil
}

func DatetimeFormatter(input string) string {
	return strings.Replace(input[:19], "T", " ", 1)
}

func DateFormatter(input string) string {
	return input[:10]
}

package customTypes

import (
	"strings"
	"time"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type MyTime time.Time

// Implement the json.Unmarshaler interface
func (t *MyTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		*t = MyTime(time.Time{})
		return nil
	}

	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		return err
	}

	parsed, err := time.ParseInLocation("2006-01-02", s, loc)
	if err != nil {
		return err
	}

	*t = MyTime(parsed)
	return nil
}

// Implement the json.Marshaler interface
func (t MyTime) MarshalJSON() ([]byte, error) {
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		return nil, err
	}

	formatted := time.Time(t).In(loc).Format("2006-01-02")
	return []byte(`"` + formatted + `"`), nil
}

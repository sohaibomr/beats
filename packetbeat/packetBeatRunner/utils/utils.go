package utils

import "time"

//StatusCode constants
const (
	StatusCodeOK  = 200
	StatusCode404 = 400
	StatusCode201 = 201
	WeekinMs      = 604800000
	DayinMs       = 86400000
)

// MakeTimestamp returns the current timestamp in millisecond
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

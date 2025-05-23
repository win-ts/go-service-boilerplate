package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// EncodeJSONtoBytes encodes an object to a JSON bytes array
func EncodeJSONtoBytes(obj any) string {
	data, err := json.Marshal(obj)
	if err != nil {
		slog.Error("[utils.EncodeJSONtoBytes] unable to encode JSON", slog.Any("error", err))
	}
	return string(data)
}

// DecodeJSONfromString decodes a JSON string to an object
func DecodeJSONfromString[T any](data string) *T {
	obj := new(T)
	if err := json.Unmarshal([]byte(data), obj); err != nil {
		slog.Error("[utils.DecodeJSONfromString] unable to decode JSON", slog.Any("error", err))
	}

	return obj
}

// Debug prints the object in a pretty format
func Debug(obj any) {
	raw, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(raw))
}

// LocalTime returns the current time in Asia/Bangkok timezone
func LocalTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

// ConvertStringTimetoTime converts a string time to time.Time
func ConvertStringTimetoTime(t string) time.Time {
	layout := "2006-01-02T15:04:05.999 -0700 MST"
	result, err := time.Parse(layout, t)
	if err != nil {
		slog.Error("[utils.convertStringTimetoTime] parse time failed", slog.Any("error", err))
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	return result.In(loc)
}

// ConvertStringDurationtoDuration converts a string duration to time.Duration
func ConvertStringDurationtoDuration(d string) time.Duration {
	result, err := time.ParseDuration(d)
	if err != nil {
		return 0
	}

	return result
}

// ConvertStringToInt converts a string to int
func ConvertStringToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return result
}

// ConvertStringToFloat64 converts a string to float64
func ConvertStringToFloat64(s string) float64 {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return result
}

// ConvertStringToBoolean converts a string to boolean
func ConvertStringToBoolean(s string) bool {
	result, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return result
}

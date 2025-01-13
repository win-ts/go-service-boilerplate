package middleware

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
)

const (
	maskChar      = "X"
	maxBase64Size = 256 * 1024
)

var sensitiveFields = map[string]func(string) string{
	"password":          maskStringAll,
	"email":             maskEmail,
	"mobile":            maskMobileNumber,
	"mobileno":          maskMobileNumber,
	"firstname":         maskStringExceptFirstAndLast,
	"lastname":          maskStringExceptFirstAndLast,
	"middlename":        maskStringExceptFirstAndLast,
	"fullname":          maskStringExceptFirstAndLast,
	"otp":               maskStringExceptFirstAndLast,
	"creditcard":        maskCardNumber,
	"mobileotpno":       maskMobileNumber,
	"nameth":            maskStringExceptFirstAndLast,
	"nameen":            maskStringExceptFirstAndLast,
	"customeremail":     maskEmail,
	"firstnameth":       maskStringExceptFirstAndLast,
	"lastnameth":        maskStringExceptFirstAndLast,
	"middlenameth":      maskStringExceptFirstAndLast,
	"firstnameen":       maskStringExceptFirstAndLast,
	"lastnameen":        maskStringExceptFirstAndLast,
	"middlenameen":      maskStringExceptFirstAndLast,
	"dateofbirth":       maskStringExceptFirstAndLast,
	"homephone":         maskMobileNumber,
	"mobilephone":       maskMobileNumber,
	"officephone":       maskMobileNumber,
	"smsmobilephone":    maskMobileNumber,
	"accesstoken":       maskStringAll,
	"phoneno":           maskStringExceptFirstAndLast,
	"requestername":     maskStringExceptFirstAndLast,
	"authorization":     maskStringAll,
	"x-api-key":         maskStringAll,
	"birthdate":         maskStringExceptFirstAndLast,
	"dathofbirth":       maskStringExceptFirstAndLast,
	"transactionnumber": maskStringExceptFirstAndLast,
	"data":              maskLargeBase64,
}

// MaskSensitiveData masks sensitive data in the given data used for logging
func MaskSensitiveData(key string, data any) any {
	storage := make(map[string]any)
	maskJSON(key, data, storage)
	return storage[key]
}

func maskJSON(key string, data any, storage map[string]any) {
	if data != nil {
		valType := reflect.TypeOf(data)
		if valType.Kind() == reflect.Ptr {
			valType = valType.Elem()
		}

		if valType.Kind() == reflect.Struct || valType.Kind() == reflect.Map {
			b, err := json.Marshal(data)
			if err == nil {
				data = string(b)
			}
		}
	}

	switch value := data.(type) {
	case map[string]any:
		innerStorage := make(map[string]any)
		maskMap(value, innerStorage)
		storage[key] = innerStorage
	case []any:
		maskArray(key, value, storage)
	case string:
		if !maskString(key, value, storage) {
			storage[key] = value
		}
	default:
		storage[key] = value
	}
}

func maskMap(data map[string]any, storage map[string]any) {
	for k, v := range data {
		maskJSON(k, v, storage)
	}
}

func maskArray(key string, data []any, storage map[string]any) {
	var resultArray []any
	for _, v := range data {
		if s, ok := v.(string); ok {
			resultArray = append(resultArray, maskStringArr(key, s))
		} else {
			innerStorage := make(map[string]any)
			maskJSON("", v, innerStorage)
			if str, ok := innerStorage[""]; ok {
				resultArray = append(resultArray, str)
			} else {
				resultArray = append(resultArray, innerStorage)
			}
		}
	}
	storage[key] = resultArray
}

func maskString(key string, value string, storage map[string]any) bool {
	if val, ok := isJSON(value); ok {
		innerStorage := make(map[string]any)
		maskMap(val, innerStorage)
		storage[key] = innerStorage
		return true
	} else if val, ok := isJSONArr(value); ok {
		maskArray(key, val, storage)
		return true
	} else if fn, ok := sensitiveFields[strings.ToLower(key)]; ok {
		storage[key] = fn(value)
		return true
	}
	return false
}

func maskStringArr(key string, value string) string {
	if fn, ok := sensitiveFields[strings.ToLower(key)]; ok {
		return fn(value)
	}

	return value
}

func isJSONArr(s string) ([]any, bool) {
	var arr []any
	return arr, json.Unmarshal([]byte(s), &arr) == nil
}

func isJSON(s string) (map[string]any, bool) {
	var m map[string]any
	return m, json.Unmarshal([]byte(s), &m) == nil
}

func maskStringAll(str string) string {
	if len(str) == 0 {
		return str
	}

	return strings.Repeat(maskChar, 6)
}

func maskEmail(str string) string {
	if !strings.ContainsAny(str, "@") {
		return str
	}

	strPart := strings.Split(str, "@")
	if len(strPart[0]) < 2 {
		return str
	}

	return strPart[0][:1] + strings.Repeat(maskChar, len(strPart[0])-2) + strPart[0][len(strPart[0])-1:] + "@" + strPart[1]
}

func maskMobileNumber(str string) string {
	if len(str) < 10 {
		return strings.Repeat(maskChar, 10)
	}
	return strings.Repeat(maskChar, len(str)-4) + str[len(str)-4:]
}

func maskStringExceptFirstAndLast(str string) string {
	if len(str) <= 2 {
		return str
	}
	return string(str[0]) + strings.Repeat(maskChar, len(str)-2) + string(str[len(str)-1])
}

func maskCardNumber(str string) string {
	if len(str) < 16 {
		str = strings.Repeat(maskChar, 16-len(str)) + str
	}
	return str[:6] + strings.Repeat(maskChar, 6) + str[6+6:]
}

func maskLargeBase64(str string) string {
	if len(str) < maxBase64Size {
		return str
	}

	_, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return str
	}

	return "LARGE_BASE64"
}

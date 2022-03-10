package helpers

import (
	"reflect"

	"github.com/gosimple/slug"
)

// Guard edits the given string.
// example: 'create $#% contact' -> 'create-contact'.
// @param string
// @param string
// return bool
func Guard(b string) string {
	return slug.Make(b)
}

// GuardArray edits the given string array.
// example: 'create $#% contact' -> 'create-contact'.
// @param []string
// return []string
func GuardArray(b []string) (guardArray []string) {
	for _, c := range b {
		guardArray = append(guardArray, slug.Make(c))
	}
	return
}

// IsInt is the given value an integer?
// @param interface{}
// return bool
func IsInt(value interface{}) bool {
	if reflect.TypeOf(value).Kind() == reflect.Int {
		return true
	}
	return false
}

// IsUInt is the given value an unsigned integer?
// @param interface{}
// return bool
func IsUInt(value interface{}) bool {
	if reflect.TypeOf(value).Kind() == reflect.Uint {
		return true
	}
	return false
}

// IsUIntArray is the given value an unsigned integer array?
// @param interface{}
// return bool
func IsUIntArray(value interface{}) bool {
	t := reflect.TypeOf(value)
	if !IsArray(value) {
		return false
	}
	if t.Elem().Kind() == reflect.Uint {
		return true
	}
	return false
}

// IsString is the given value an string?
// @param interface{}
// return bool
func IsString(value interface{}) bool {
	if reflect.TypeOf(value).Kind() == reflect.String {
		return true
	}
	return false
}

// IsStringArray is the given value an string array?
// @param interface{}
// return bool
func IsStringArray(value interface{}) bool {
	t := reflect.TypeOf(value)
	if !IsArray(value) {
		return false
	}
	if t.Elem().Kind() == reflect.String {
		return true
	}
	return false
}

// IsArray is the given value an array?
// @param interface{}
// return bool
func IsArray(value interface{}) bool {
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return false
	}
	return true
}

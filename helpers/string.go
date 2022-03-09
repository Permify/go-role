package helpers

import (
	`github.com/gosimple/slug`
	`reflect`
)

// Guard

func Guard(b string) string {
	return slug.Make(b)
}

func GuardArray(b []string) (guardArray []string) {
	for _, c := range b {
		guardArray = append(guardArray, slug.Make(c))
	}
	return
}

// Is

func IsInt(value interface{}) bool {
	if reflect.TypeOf(value).Kind() == reflect.Int {
		return true
	}
	return false
}

func IsUIntArray(value interface{}) bool {
	t := reflect.TypeOf(value)
	if !IsArray(value){
		return false
	}
	if t.Elem().Kind() == reflect.Uint {
		return true
	}
	return false
}

func IsString(value interface{}) bool {
	if reflect.TypeOf(value).Kind() == reflect.String {
		return true
	}
	return false
}

func IsStringArray(value interface{}) bool {
	t := reflect.TypeOf(value)
	if !IsArray(value){
		return false
	}
	if t.Elem().Kind() == reflect.String {
		return true
	}
	return false
}

func IsArray(value interface{}) bool {
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return false
	}
	return true
}

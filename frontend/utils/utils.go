package utils

import (
	"reflect"
)

func CheckType(received any, expects ...any) bool {
	receivedType := reflect.TypeOf(received)
	for _, e := range expects {
		expectedType := reflect.TypeOf(e)
		if receivedType == expectedType {
			return true
		}
	}
	return false
}

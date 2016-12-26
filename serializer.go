package main

import (
	"errors"
	"fmt"
	"reflect"
)

func Serialize(object interface{}) (string, error) {
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind == reflect.String {
		return fmt.Sprintf("\"%s\"", value), nil
	} else if kind == reflect.Int ||
		kind == reflect.Float64 {
		return fmt.Sprintf("%v", value), nil
	} else if kind == reflect.Slice {

	} else {
		return "", errors.New("No serializer for type '" + kind.String() + "'")
	}
}

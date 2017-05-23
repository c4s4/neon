package util

import (
	"fmt"
	"net"
	"reflect"
	"time"
	"unicode/utf8"
)

func ToList(object interface{}) ([]interface{}, error) {
	slice := reflect.ValueOf(object)
	if slice.Kind() == reflect.Slice {
		result := make([]interface{}, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			result[i] = slice.Index(i).Interface()
		}
		return result, nil
	} else {
		return nil, fmt.Errorf("must be a list")
	}
}

func ValueToInterface(value reflect.Value) interface{} {
	if value.IsValid() {
		kind := value.Kind()
		if kind == reflect.Slice {
			result := make([]interface{}, value.Len())
			for i := 0; i < value.Len(); i++ {
				result[i] = ValueToInterface(value.Index(i))
			}
			return result
		} else if kind == reflect.Map {
			result := make(map[interface{}]interface{})
			for _, key := range value.MapKeys() {
				keyInterface := ValueToInterface(key)
				valueInterface := ValueToInterface(value.MapIndex(key))
				result[keyInterface] = valueInterface
			}
			return result
		} else {
			return value.Interface()
		}
	} else {
		return nil
	}
}

func MaxLength(lines []string) int {
	length := 0
	for _, line := range lines {
		if utf8.RuneCountInString(line) > length {
			length = utf8.RuneCountInString(line)
		}
	}
	return length
}

func Singleton(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	go func() {
		for {
			listener.Accept()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return nil
}

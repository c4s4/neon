package build

import (
	"neon/util"
	"reflect"
	"fmt"
	"strings"
)

func ParseArgs(object *util.Object, args interface{}) (interface{}, error){
	st := reflect.TypeOf(args)
	value := reflect.ValueOf(args)
	if st.Kind() != reflect.Struct {
		return nil, fmt.Errorf("args must be a struct")
	}
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tagMandatory := field.Tag.Get("mandatory")
		if tagMandatory != "true" && tagMandatory != "false" {
			return nil, fmt.Errorf("tag 'mandatory' must be 'true' or 'false'")
		}
		mandatory := tagMandatory == "true"
		var name string
		tagName := field.Tag.Get("name")
		if tagName == "" {
			name = strings.ToLower(field.Name)
		} else {
			name = tagName
		}
		if object.HasField(name) {
			value.
		}
	}
	return args, nil
}

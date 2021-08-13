package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// From a type and a json name, return the file that contains is filed
func GetFieldName(tag string, key string, instance interface{}) (*string, error) {
	reflectType := reflect.TypeOf(instance)
	if reflectType.Kind() != reflect.Struct {
		return nil, errors.New(fmt.Sprintf("Instance is not a struct but it is %s", reflectType.Kind()))
	}
	for i := 0; i < reflectType.NumField(); i++ {
		reflectFiled := reflectType.Field(i)
		reflectValue := strings.Split(reflectFiled.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if reflectValue == tag {
			return &reflectFiled.Name, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Filed with json name %s not found", tag))
}

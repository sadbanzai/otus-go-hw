package hw09structvalidator

import (
	"errors"
	"reflect"
)

var ErrValueTypeIsNotStruct = errors.New("value type is not struct")

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrValueTypeIsNotStruct
	}

	return nil
}

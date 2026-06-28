package myerrors

import (
	"errors"
	"reflect"
)

func HasSameTypeAs(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}
	ptr := reflect.New(reflect.TypeOf(target))
	return errors.As(err, ptr.Interface())
}

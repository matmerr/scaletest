package utils

import (
	"reflect"
)

func GetPackagePath(data interface{}) (string, error) {
	t := reflect.TypeOf(data)
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		return ".", nil
	}
	return pkgPath, nil
}

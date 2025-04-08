package main

import (
	"fmt"
	"log"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	val := reflect.ValueOf(data)
	res := reflect.ValueOf(out)
	if res.Kind() != reflect.Ptr || res.IsNil() {
		return fmt.Errorf("type not pointer or nil")
	}

	if err := set(val, res.Elem()); err != nil {
		return err
	}

	return nil
}

func set(val, res reflect.Value) error {
	if val.Kind() == reflect.Interface {
		val = reflect.ValueOf(val.Interface())
	}

	log.Println("val:", val, "kind:", val.Kind(), "||", "res:", res, "kind:", res.Kind())

	switch res.Kind() {

	case reflect.Struct:
		if val.Kind() != reflect.Map || val.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("wrong type value in")
		}
		for _, key := range val.MapKeys() {
			field := res.FieldByName(key.String())
			if !field.IsValid() {
				continue
			} else {
				err := set(val.MapIndex(key), field)
				if err != nil {
					return err
				}
			}
		}

	case reflect.Slice:
		if val.Kind() != reflect.Slice {
			return fmt.Errorf("wrong type value in")
		}
		len := val.Len()
		res.Set(reflect.MakeSlice(res.Type(), len, len))
		for i := 0; i < len; i++ {
			err := set(val.Index(i), res.Index(i))
			if err != nil {
				return err
			}
		}

	case reflect.String:
		if val.Kind() != reflect.String {
			return fmt.Errorf("wrong type value in")
		}
		res.SetString(val.String())

	case reflect.Int:
		if val.Kind() != reflect.Float64 {
			return fmt.Errorf("wrong type value in")
		}
		res.SetInt(int64(val.Float()))

	case reflect.Bool:
		if val.Kind() != reflect.Bool {
			return fmt.Errorf("wrong type value in")
		}
		res.SetBool(val.Bool())

	default:
		return fmt.Errorf("unknow type")
	}

	return nil
}

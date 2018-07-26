package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/fatih/structs"
	"strings"
)

func SmartFillMapArray(target interface{}, data map[string][]string) error {
	req := structs.New(target)

	for _, field := range req.Fields() {
		name := field.Name()
		if value, ok := data[name]; ok && len(value) != 0 && len(value[0]) != 0 {
			if err := SmartFill(field, value[0]); err != nil {
				return err
			}

			continue
		}

		name = CamelCaseToUnderscore(field.Name())
		if value, ok := data[name]; ok && len(value) != 0 && len(value[0]) != 0 {
			if err := SmartFill(field, value[0]); err != nil {
				return err
			}

			continue
		}
	}
	return nil
}

func SmartFillMap(target interface{}, data map[string]string) error {
	req := structs.New(target)

	for _, field := range req.Fields() {
		name := field.Name()
		if value, ok := data[name]; ok && len(value) != 0 {
			if err := SmartFill(field, value); err != nil {
				return err
			}

			continue
		}

		name = CamelCaseToUnderscore(field.Name())
		if value, ok := data[name]; ok && len(value) != 0 {
			if err := SmartFill(field, value); err != nil {
				return err
			}

			continue
		}
	}
	return nil
}

func SmartFill(field *structs.Field, data string) error {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := fillInt(field, data); err != nil {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := fillUint(field, data); err != nil {
			return err
		}
	case reflect.Bool:
		if err := fillBool(field, data); err != nil {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if err := fillFloat(field, data); err != nil {
			return err
		}
	case reflect.String:
		return field.Set(data)
	}

	return nil
}

func fillInt(field *structs.Field, data string) error {
	if data, err := strconv.ParseInt(data, 10, 64); err != nil {
		return err
	} else {
		switch field.Kind() {
		case reflect.Int:
			return field.Set(int(data))
		case reflect.Int8:
			return field.Set(int8(data))
		case reflect.Int16:
			return field.Set(int16(data))
		case reflect.Int32:
			return field.Set(int32(data))
		case reflect.Int64:
			return field.Set(int64(data))
		default:
			return fmt.Errorf("smart convert fill int type error, field: %s", field)
		}
	}

	return nil
}

func fillUint(field *structs.Field, data string) error {
	if data, err := strconv.ParseUint(data, 10, 64); err != nil {
		return err
	} else {
		switch field.Kind() {
		case reflect.Uint:
			return field.Set(uint(data))
		case reflect.Uint8:
			return field.Set(uint8(data))
		case reflect.Uint16:
			return field.Set(uint16(data))
		case reflect.Uint32:
			return field.Set(uint32(data))
		case reflect.Uint64:
			return field.Set(uint64(data))
		default:
			return fmt.Errorf("smart convert fill uint type error, field: %s", field)
		}
	}

	return nil
}

func fillFloat(field *structs.Field, data string) error {
	if data, err := strconv.ParseFloat(data, 64); err != nil {
		return err
	} else {
		switch field.Kind() {
		case reflect.Float32:
			return field.Set(float32(data))
		case reflect.Float64:
			return field.Set(float64(data))
		default:
			return fmt.Errorf("smart convert fill uint type error, field: %s", field)
		}
	}

	return nil
}

func fillBool(field *structs.Field, data string) error {
	if data, err := strconv.ParseBool(data); err != nil {
		return err
	} else {
		return field.Set(data)
	}

	return nil
}

func StructsToMap(target interface{}) (*map[string]interface{}, error) {
	req := structs.New(target)
	var result = make(map[string]interface{})
	for _, field := range req.Fields() {
		name := field.Tag("json")
		nameList := strings.Split(name, ",")
		result[nameList[0]] = field.Value()
	}
	return &result, nil
}

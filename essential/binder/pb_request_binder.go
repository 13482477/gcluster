package binder

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	log "github.com/sirupsen/logrus"
	"poseidon/apollo/common/help"
)

type PbRequestBinder struct{}

func (b *PbRequestBinder) Bind(i interface{}, c echo.Context) (err error) {
	req := c.Request()
	var fillFields []string
	if req.ContentLength == 0 {
		if req.Method == echo.GET || req.Method == echo.DELETE {

			//有解密过
			d := c.Get("decrypted")

			if d != nil && d.(int) == 1 {
				ds := c.Get("d_data")
				var data map[string]interface{}
				log.Debug(" d_data", string(ds.([]byte)))
				err := json.Unmarshal(ds.([]byte), &data)
				if err != nil {
					return err
				}

				fakeData := make(map[string][]string)
				for k,v := range data {
					fakeData[k] = []string{help.GetAssertString(v)}
				}
				fillFields, err = b.bindData(i, fakeData)
			} else {
				fillFields, err = b.bindData(i, c.QueryParams())
			}

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		}
	} else {
		ctype := req.Header.Get(echo.HeaderContentType)
		switch {
		case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
			stringMap := make(map[string]interface{})
			newStringMap := make(map[string]interface{})

			d := c.Get("decrypted")

			log.Debug(" the decrypted data is ", d)

			if d != nil && d.(int) == 1 {
				ds := c.Get("d_data")
				err := json.Unmarshal(ds.([]byte), &stringMap)
				if err != nil {
					return err
				}
			} else {
				decoder := json.NewDecoder(req.Body)
				decoder.UseNumber()
				if err = decoder.Decode(&stringMap); err != nil {
					return err
				}
			}

			typ := reflect.TypeOf(i).Elem()
			sprops := proto.GetProperties(typ)
			var body []byte
			for _, prop := range sprops.Prop {
				orig := prop.OrigName
				camel := prop.OrigName
				if prop.JSONName != "" {
					camel = prop.JSONName
				}

				valOrig, okOrig := stringMap[orig]
				valCamel, okCamel := stringMap[camel]

				if okOrig {
					fillFields = append(fillFields, prop.Name)
					newStringMap[orig] = valOrig
				} else if okCamel {
					fillFields = append(fillFields, prop.Name)
					newStringMap[camel] = valCamel
				}
				body, err = json.Marshal(newStringMap)
				if err != nil {
					return nil
				}
			}
			if err = jsonpb.Unmarshal(bytes.NewReader(body), i.(proto.Message)); err != nil {
				return err
			}
		case strings.HasPrefix(ctype, echo.MIMEApplicationForm), strings.HasPrefix(ctype, echo.MIMEMultipartForm):
			params, err := c.FormParams()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			fillFields, err = b.bindData(i, params)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
		default:
			return echo.ErrUnsupportedMediaType
		}
	}

	requestFields := reflect.ValueOf(i).Elem().FieldByName("RequestFields")
	if requestFields.IsValid() && requestFields.CanSet() {
		requestFields.Set(reflect.ValueOf(fillFields))
	}
	return
}

func (b *PbRequestBinder) bindData(ptr interface{}, data map[string][]string) ([]string, error) {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()

	fillFields := make([]string, 0)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("Binding element must be a struct ")
	}

	sprops := proto.GetProperties(typ)

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind()
		prop := sprops.Prop[i]

		orig := prop.OrigName
		camel := prop.OrigName
		if prop.JSONName != "" {
			camel = prop.JSONName
		}

		vOrig, okOrig := data[orig]
		vCamel, okCamel := data[camel]

		if okOrig || okCamel {
			fillFields = append(fillFields, typeField.Name)
		} else {
			continue
		}
		var inputValue []string
		if okOrig {
			inputValue = vOrig
		}
		if okCamel {
			inputValue = vCamel
		}

		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok {
			if err != nil {
				return nil, err
			}
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ {
				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
					return nil, err
				}
			}
			val.Field(i).Set(slice)
		} else {
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
				return nil, err
			}
		}
	}
	return fillFields, nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	if ok, err := unmarshalField(valueKind, val, structField); ok {
		return err
	}

	switch valueKind {
	case reflect.Ptr:
		return setWithProperType(structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func unmarshalField(valueKind reflect.Kind, val string, field reflect.Value) (bool, error) {
	switch valueKind {
	case reflect.Ptr:
		return unmarshalFieldPtr(val, field)
	default:
		return unmarshalFieldNonPtr(val, field)
	}
}

func bindUnmarshaler(field reflect.Value) (echo.BindUnmarshaler, bool) {
	ptr := reflect.New(field.Type())
	if ptr.CanInterface() {
		iface := ptr.Interface()
		if unmarshaler, ok := iface.(echo.BindUnmarshaler); ok {
			return unmarshaler, ok
		}
	}
	return nil, false
}

func unmarshalFieldNonPtr(value string, field reflect.Value) (bool, error) {
	if unmarshaler, ok := bindUnmarshaler(field); ok {
		err := unmarshaler.UnmarshalParam(value)
		field.Set(reflect.ValueOf(unmarshaler).Elem())
		return true, err
	}
	return false, nil
}

func unmarshalFieldPtr(value string, field reflect.Value) (bool, error) {
	if field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
	}
	return unmarshalFieldNonPtr(value, field.Elem())
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

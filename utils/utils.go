package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	json.NewEncoder(w).Encode(data)
}

func UnmarshalJSON(data []byte, p interface{}) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	v := reflect.ValueOf(p).Elem()
	t := v.Type()

	var missing []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val, ok := m[field.Tag.Get("json")]
		delete(m, field.Tag.Get("json"))
		if !ok {
			missing = append(missing, field.Tag.Get("json"))
			continue
		}

		switch field.Type.Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			if reflect.TypeOf(val).Kind() == reflect.Float64 | reflect.Float32{
				v.Field(i).Set(reflect.ValueOf(val).Convert(field.Type))
			}
		default:
			v.Field(i).Set(reflect.ValueOf(val))
		}
	}

	if len(missing) > 0 {
		return errors.New("missing fields: " + strings.Join(missing, ", "))
	}

	if len(m) > 0 {
		extra := make([]string, 0, len(m))
		for field := range m {
			extra = append(extra, field)
		}
		return errors.New("unknown fields: " + strings.Join(extra, ", "))
	}

	return nil
}

func GetField(values []interface{}, field string) interface{}{
	return values[0].(map[string]interface{})[field]
}
package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// DBQueryBuilder ... Main struct for constructing query
type DBQueryBuilder struct {
	httpUtilFunc HTTPUtilityFunctions
}

// ParseQueryFilter ...
func (c DBQueryBuilder) ParseQueryFilter(query string) ([]QueryFilter, error) {
	var queryDict interface{}
	err := json.Unmarshal([]byte(query), &queryDict)
	if err == nil {
		m := queryDict.(map[string]interface{})
		fmt.Println(m)
	}
	return nil, err
}

//QueryStruct ...
type QueryStruct struct {
	filters []QueryFilter
}

//QueryFilter ...
type QueryFilter struct {
	name string
	op   string
	val  string
}

// FillStruct ...
func (s *QueryFilter) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetField ...
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

package generics

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var timeType = reflect.TypeOf(time.Now())

// Convert value to the type. This is a more aggressive conversion, that will
// recursively create new map and slice values as necessary. It doesn't
// handle circular references.
func Convert(value interface{}, t reflect.Type) (interface{}, error) {
	r := reflect.ValueOf(value)
	if r.Type().ConvertibleTo(t) {
		return r.Convert(t).Interface(), nil
	}
	if reflect.PtrTo(r.Type()) == t {
		return &value, nil
	}
	if r.Kind() == reflect.String && t == timeType {
		return ParseTime(value.(string))
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch value := value.(type) {
		case bool:
			if value {
				return 1, nil
			}
			return 0, nil
		case string:
			return strconv.Atoi(value)
		}
	case reflect.Float32, reflect.Float64:
		switch value := value.(type) {
		case int:
			return float64(value), nil
		case string:
			return strconv.ParseFloat(value, 64)
		}
	case reflect.Slice:
		if r.Kind() != reflect.Array && r.Kind() != reflect.Slice {
			break
		}
		out := reflect.MakeSlice(t, 0, r.Len())
		for i := 0; i < r.Len(); i++ {
			item, err := Convert(r.Index(i).Interface(), t.Elem())
			if err != nil {
				return nil, err
			}
			out = reflect.Append(out, reflect.ValueOf(item))
		}
		return out.Interface(), nil
	}
	return nil, genericErrorf("generic.Convert can't convert %#v<%s> to %v", value, r.Type(), t)
}

// MustConvert wraps Convert, but panics on error.
func MustConvert(value interface{}, t reflect.Type) interface{} {
	out, err := Convert(value, t)
	if err != nil {
		panic(err)
	}
	return out
}

// MustConvertItem converts item to conform to array, else panics.
func MustConvertItem(item interface{}, array []interface{}) interface{} {
	item, err := Convert(item, reflect.TypeOf(array).Elem())
	if err != nil {
		panic(fmt.Errorf("can't convert %#v to %s: %s", item, reflect.TypeOf(array).Elem(), err))
	}
	return item
}

var dateLayouts = []string{
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -4",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"January 2, 2006",
	"January 2 2006",
	"Jan 2, 2006",
	"Jan 2 2006",
}

// ParseTime tries a few heuristics to parse a date from a string
func ParseTime(value string) (time.Time, error) {
	if value == "now" {
		return time.Now(), nil
	}
	for _, layout := range dateLayouts {
		// fmt.Println(layout, time.Now().Format(layout), value)
		time, err := time.Parse(layout, value)
		if err == nil {
			return time, nil
		}
	}
	return time.Now(), genericErrorf("can't convert %s to a time", value)
}
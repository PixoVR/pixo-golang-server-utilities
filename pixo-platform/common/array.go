package common

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

func IsInSlice(slice interface{}, value interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("isInSlice: slice parameter must be a slice")
	}

	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(s.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

func RemoveNilPointers(slice interface{}) interface{} {
	if slice == nil {
		return nil
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		log.Error().Msgf("RemoveNullPointers only accepts slices, received %v", v.Kind())
		return nil
	}

	elemType := v.Type().Elem().Elem()
	newSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Elem()
		if elem.IsValid() {
			newSlice = reflect.Append(newSlice, elem)
		}
	}

	return newSlice.Interface()
}

func RemoveNilValues[T any](slice []T) []T {
	if slice == nil {
		return nil
	}

	newSlice := make([]T, 0, len(slice))
	for _, v := range slice {
		if any(v) != nil {
			newSlice = append(newSlice, v)
		}
	}

	return newSlice
}

func ToPointerArray(slice interface{}) interface{} {
	if slice == nil {
		return nil
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		log.Error().Msgf("ToPointerArray only accepts slices, received %v", v.Kind())
		return nil
	}

	elemType := reflect.PointerTo(v.Type().Elem())
	newSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		newSlice = reflect.Append(newSlice, elem.Addr())
	}

	return newSlice.Interface()
}

func GetDistinctValues[T any, V comparable](items []T, selector func(T) V) []V {
	distinctMap := make(map[V]struct{})

	for _, item := range items {
		key := selector(item)
		distinctMap[key] = struct{}{}
	}

	var distinctValues []V
	for key := range distinctMap {
		distinctValues = append(distinctValues, key)
	}

	return distinctValues
}

func ToStringArray[T any](input []T) []string {
	stringArray := make([]string, len(input))
	for i, v := range input {
		stringArray[i] = fmt.Sprintf("%v", v)
	}
	return stringArray
}

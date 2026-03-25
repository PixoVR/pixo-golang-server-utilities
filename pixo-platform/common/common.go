package common

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/99designs/gqlgen/graphql"
)

// Filter returns only the filtered elements
//
// example:
// Filter([]int{1, 2, 3, 4, 5}, func(i int) bool { return i > 2 })
func Filter[T any](data []T, filter func(T) bool) []T {
	var result []T
	for _, d := range data {
		if filter(d) {
			result = append(result, d)
		}
	}
	return result
}

func Take[T any](data []T, limit int) []T {
	if len(data) == 0 || limit <= 0 {
		return []T{}
	}
	if len(data) < limit {
		return data
	}
	return data[:limit]
}

func Find[T any](data []T, filter func(T) bool) *T {
	for _, d := range data {
		if filter(d) {
			return &d
		}
	}
	return nil
}

func Pluck[T any, K any](data []T, key func(T) K) []K {
	var result []K
	for _, d := range data {
		result = append(result, key(d))
	}
	return result
}

func Average[T any](data []T, key func(T) float64) float64 {
	if len(data) == 0 {
		return 0
	}
	var total float64
	for _, d := range data {
		total += key(d)
	}
	return total / float64(len(data))
}

func Contains[T comparable](data []T, item T) bool {
	return slices.Contains(data, item)
}

func GetFieldValue(obj interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct or a pointer to a struct")
	}

	field := v.FieldByName(fieldName)

	if !field.IsValid() {
		return nil, fmt.Errorf("no such field: %s in obj", fieldName)
	}

	return field.Interface(), nil
}

// TryGetParentGraphqlResultOfType Navigates through the context to find a previously resolved value of type T
func TryGetParentGraphqlResultOfType[T any](ctx context.Context) (*T, bool) {
	var result *T
	var ok bool

	for field := graphql.GetFieldContext(ctx); field != nil; field = field.Parent {
		if result, ok = field.Result.(*T); ok {
			return result, ok
		}
	}

	return nil, false
}

func Values[K comparable, T any](m map[K]T) []T {
	values := make([]T, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

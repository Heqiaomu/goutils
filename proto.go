package util

import (
	"fmt"
	"strconv"

	"reflect"

	"github.com/gogo/protobuf/types"
)

func ConvertToValue(v interface{}, from reflect.Type) (*types.Value, error) {

	// String type to Value_StringValue
	if from.Kind() == reflect.String {
		return &types.Value{Kind: &types.Value_StringValue{StringValue: reflect.ValueOf(v).String()}}, nil
	}

	// Int64 and Float64 type to Value_NumberValue
	if from.Kind() == reflect.Int64 || from.Kind() == reflect.Float64 {
		return &types.Value{Kind: &types.Value_NumberValue{NumberValue: reflect.ValueOf(v).Float()}}, nil
	}

	// String type and if its true or false
	if from.Kind() == reflect.String && (reflect.ValueOf(v).String() == "true" || reflect.ValueOf(v).String() == "false") || from.Kind() == reflect.Bool {
		boolVal, err := strconv.ParseBool(reflect.ValueOf(v).String())
		if err != nil {
			return nil, err
		}
		return &types.Value{Kind: &types.Value_BoolValue{BoolValue: boolVal}}, nil
	}

	// Slice type to Value_ListValue
	if from.Kind() == reflect.Slice {
		list := reflect.ValueOf(v)
		outputList := make([]*types.Value, 0)
		for i := 0; i < list.Len(); i++ {
			val, err := ConvertToValue(list.Index(i).Interface(), reflect.TypeOf(list.Index(i).Interface()))
			if err != nil {
				return nil, err
			}
			outputList = append(outputList, val)
		}
		return &types.Value{Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: outputList}}}, nil
	}

	// Struct type to Value_StructValue
	if from.Kind() == reflect.Struct {
		valStruct := reflect.ValueOf(v)
		outputMap := make(map[string]*types.Value)
		numOfField := valStruct.NumField()

		for i := 0; i < numOfField; i++ {
			val, err := ConvertToValue(valStruct.Field(i).Interface(), reflect.TypeOf(valStruct.Field(i).Interface()))
			if err != nil {
				return nil, err
			}
			outputMap[valStruct.Field(i).Type().Name()] = val
		}

		return &types.Value{Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: outputMap}}}, nil
	}

	return nil, nil
}

func ConvertValueToInterface(value types.Value) (interface{}, error) {

	// null value
	if _, ok := value.Kind.(*types.Value_NullValue); ok {
		return nil, nil
	}

	// string value to string
	if x, ok := value.Kind.(*types.Value_StringValue); ok {
		fmt.Println("Coming inside string for : ", x.StringValue)
		return x.StringValue, nil
	}

	// number value to float64
	if x, ok := value.Kind.(*types.Value_NumberValue); ok {
		return x.NumberValue, nil
	}

	// bool value to string
	if x, ok := value.Kind.(*types.Value_BoolValue); ok {
		return strconv.FormatBool(x.BoolValue), nil
	}

	// list value to slice
	if x, ok := value.Kind.(*types.Value_ListValue); ok {
		if x == nil || x.ListValue == nil || x.ListValue.Values == nil {
			return nil, nil
		}
		listValue := x.ListValue.Values
		if len(listValue) == 0 {
			return nil, nil
		}

		val, err := ConvertValueToInterface(*listValue[0])
		if err != nil {
			return nil, err
		}
		typ := reflect.TypeOf(val)
		outputList := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)

		for _, value := range listValue {
			if value != nil {
				val, err := ConvertValueToInterface(*value)
				if err != nil {
					return nil, err
				}
				outputList = reflect.Append(outputList, reflect.ValueOf(val))
			}
		}
		return outputList.Interface(), nil
	}

	// struct value to struct
	if x, ok := value.Kind.(*types.Value_StructValue); ok {

		if x == nil || x.StructValue == nil || x.StructValue.Fields == nil {
			return nil, nil
		}

		mapValue := x.StructValue.Fields

		var keyTyp reflect.Type
		var typ reflect.Type

		for key, value := range mapValue {
			if value != nil {
				val, err := ConvertValueToInterface(*value)
				if err != nil {
					return nil, err
				}
				keyTyp = reflect.TypeOf(key)
				typ = reflect.TypeOf(val)
				break
			} else {
				return nil, nil
			}
		}

		outputMap := reflect.MakeMap(reflect.MapOf(keyTyp, typ))

		for key, value := range mapValue {
			if value != nil {
				val, err := ConvertValueToInterface(*value)
				if err != nil {
					return nil, err
				}
				outputMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
			}
		}

		return outputMap.Interface(), nil
	}

	return nil, nil
}

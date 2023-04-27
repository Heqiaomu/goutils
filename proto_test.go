package util

import (
	"reflect"
	"testing"

	"github.com/gogo/protobuf/types"
)

var outputList []*types.Value
var outputMap map[string]*types.Value

func init() {
	outputList = make([]*types.Value, 0)
	outputList = append(outputList,
		&types.Value{Kind: &types.Value_NumberValue{NumberValue: reflect.ValueOf(float64(1)).Float()}},
		&types.Value{Kind: &types.Value_NumberValue{NumberValue: reflect.ValueOf(float64(2)).Float()}},
	)
	outputMap = make(map[string]*types.Value)
	outputMap[reflect.TypeOf("abc").Name()] = &types.Value{Kind: &types.Value_StringValue{StringValue: "abc"}}
}

type arg struct {
	V string `json:"v"`
}

// string int64/float64 bool slice struct
func TestConvertToValue(t *testing.T) {

	type args struct {
		v    interface{}
		from reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Value
		wantErr bool
	}{
		{
			"TestConvertToValue_1",
			args{
				"abc",
				reflect.TypeOf("abc"),
			},
			&types.Value{Kind: &types.Value_StringValue{StringValue: reflect.ValueOf("abc").String()}},
			false,
		},
		{
			"TestConvertToValue_2",
			args{
				float64(123),
				reflect.TypeOf(float64(123)),
			},
			&types.Value{Kind: &types.Value_NumberValue{NumberValue: reflect.ValueOf(float64(123)).Float()}},
			false,
		},
		// need to fix
		{
			"TestConvertToValue_3",
			args{
				"true",
				reflect.TypeOf(true),
			},
			&types.Value{Kind: &types.Value_BoolValue{BoolValue: true}},
			false,
		},
		{
			"TestConvertToValue_4",
			args{
				[]float64{1, 2},
				reflect.TypeOf([]float64{1, 2}),
			},
			&types.Value{Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: outputList}}},
			false,
		},
		{
			"TestConvertToValue_5",
			args{
				arg{"abc"},
				reflect.TypeOf(arg{"abc"}),
			},
			&types.Value{Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: outputMap}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToValue(tt.args.v, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertValueToInterface(t *testing.T) {
	type args struct {
		value types.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"TestConvertValueToInterface_1",
			args{
				types.Value{Kind: &types.Value_NullValue{NullValue: types.NullValue_NULL_VALUE}},
			},
			nil,
			false,
		},
		{
			"TestConvertValueToInterface_2",
			args{
				types.Value{Kind: &types.Value_StringValue{StringValue: reflect.ValueOf("abc").String()}},
			},
			"abc",
			false,
		},
		{
			"TestConvertValueToInterface_3",
			args{
				types.Value{Kind: &types.Value_NumberValue{NumberValue: reflect.ValueOf(float64(123)).Float()}},
			},
			float64(123),
			false,
		},
		{
			"TestConvertValueToInterface_4",
			args{
				types.Value{Kind: &types.Value_BoolValue{BoolValue: true}},
			},
			"true",
			false,
		},
		{
			"TestConvertValueToInterface_5",
			args{
				types.Value{Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: outputList}}},
			},
			[]float64{1, 2},
			false,
		},
		{
			"TestConvertValueToInterface_6",
			args{
				types.Value{Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: outputMap}}},
			},
			// arg{"abc"},
			map[string]string{"string": "abc"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertValueToInterface(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertValueToInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertValueToInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "slice in list",
			args: args{a: "a", list: []string{"a", "b", "c"}},
			want: true,
		},
		{
			name: "slice not int list",
			args: args{a: "d", list: []string{"a", "b", "c"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringInSliceConv(t *testing.T) {
	Convey("TestStringInSliceConv", t, func() {
		Convey("string in slice", func() {
			So(StringInSlice("a", []string{"a", "b", "c"}), ShouldBeTrue)
		})
		Convey("string not in slice", func() {
			So(StringInSlice("d", []string{"a", "b", "c"}), ShouldBeFalse)
		})
	})
}

func TestIntInSlice(t *testing.T) {
	type args struct {
		a    int
		list []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "slice in list",
			args: args{a: 1, list: []int{1, 2, 3}},
			want: true,
		},
		{
			name: "slice not int list",
			args: args{a: 4, list: []int{1, 2, 3}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("IntInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

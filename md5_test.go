package util

import "testing"

func TestMD5(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "testmd5",
			args: struct{ text string }{text: "hello"},
			want: "5d41402abc4b2a76b9719d911017c592",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5(tt.args.text); got != tt.want {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5Part(t *testing.T) {
	type args struct {
		text  string
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "testmd5part1",
			args: struct {
				text  string
				start int
				end   int
			}{text: "hello", start: 0, end: 1},
			want: "5",
		},
		{
			name: "testmd5part2",
			args: struct {
				text  string
				start int
				end   int
			}{text: "hello", start: 34, end: 0},
			want: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name: "testmd5part3",
			args: struct {
				text  string
				start int
				end   int
			}{text: "hello", start: 2, end: 10},
			want: "41402abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5Part(tt.args.text, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("MD5Part() = %v, want %v", got, tt.want)
			}
		})
	}
}
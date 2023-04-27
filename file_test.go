package util

import (
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// last case should invoke this function
func clear() {
	os.RemoveAll("a")
}

func TestExt(t *testing.T) {
	type args struct {
		filename  string
		checkList []string
	}
	tests := []struct {
		name    string
		args    args
		wantExt string
		wantErr bool
	}{
		{
			name:    "get ext success",
			args:    args{filename: "a.toml", checkList: []string{"toml", "yml", "ini"}},
			wantExt: "toml",
			wantErr: false,
		},
		{
			name:    "get ext failed",
			args:    args{filename: "a.exe", checkList: []string{"toml", "yml", "ini"}},
			wantExt: "exe",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt, err := Ext(tt.args.filename, tt.args.checkList...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExt != tt.wantExt {
				t.Errorf("Ext() = %v, want %v", gotExt, tt.wantExt)
			}
		})
	}
}

func TestExtConv(t *testing.T) {
	Convey("TestExtConv", t, func() {
		Convey("get ext success", func() {
			ext, err := Ext("a.toml", "toml", "yml", "ini")
			So(ext, ShouldEqual, "toml")
			So(err, ShouldBeNil)
		})
		Convey("get ext failed", func() {
			ext, err := Ext("a.exe", "toml", "yml", "ini")
			So(ext, ShouldEqual, "exe")
			So(err, ShouldNotBeNil)
		})
		Convey("extension valid check fail", func() {
			ext, err := Ext("a", "toml", "yml", "ini")
			So(ext, ShouldEqual, "")
			So(err, ShouldNotBeNil)
			ext, err = Ext("a.a", "toml", "yml", "ini")
			So(ext, ShouldEqual, "a")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAdvanceGenR(t *testing.T) {
	type args struct {
		data interface{}
		tplT []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				data: map[string]string{
					"A": "123",
				},
				tplT: []byte(`
----
{{indent 4 .A}}
`),
			},
			want: []byte(`
----
    123
`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AdvanceGenR(tt.args.data, tt.args.tplT)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdvanceGenR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdvanceGenR() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDirList(t *testing.T) {
	os.MkdirAll("a/b/c", os.ModePerm)
	f, _ := os.Create("a/b/c/d.txt")
	f.WriteString("this is test")

	type args struct {
		dirpath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"TestGetDirList_1",
			args{
				"a",
			},
			[]string{"a/b"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDirList(tt.args.dirpath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDirList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDirList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"TestRemoveDirContents",
			args{
				"a",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExist(tt.args.path); got != tt.want {
				t.Errorf("FileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGen(t *testing.T) {
	os.MkdirAll("a", os.ModePerm)
	f, _ := os.Create("a/testGen.tpl")
	f.WriteString("{{.A}}")
	type args struct {
		data     interface{}
		tpl      string
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestGen_1",
			args: args{
				data: map[string]string{
					"A": "123",
				},
				tpl:      "a/testGen.tpl",
				filepath: "a/testGen.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Gen(tt.args.data, tt.args.tpl, tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("Gen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenR(t *testing.T) {
	type args struct {
		data interface{}
		tplT []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				data: map[string]string{
					"A": "123",
				},
				tplT: []byte("{{.A}}"),
			},
			want:    []byte("123"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenR(tt.args.data, tt.args.tplT)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	type args struct {
		data []byte
		f    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteFile(tt.args.data, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreatedDirIfNotExist(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestCreatedDirIfNotExist_1",
			args{
				"a/TestCreatedDirIfNotExist",
			},
			false,
		},
		{
			"TestCreatedDirIfNotExist_2",
			args{
				"a",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreatedDirIfNotExist(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("CreatedDirIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoveDirContents(t *testing.T) {
	defer func() {
		clear()
	}()
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestRemoveDirContents",
			args{
				"a",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveDirContents(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("RemoveDirContents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

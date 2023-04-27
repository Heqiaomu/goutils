package compress

import (
	"os"
	"path"
	"testing"
)

const p = "./"

func init() {
	os.MkdirAll(path.Join(p, "testdata/before"), os.ModePerm)
	os.MkdirAll(path.Join(p, "testdata/after"), os.ModePerm)
}

// last case should invoke this function
func clear() {
	os.RemoveAll(path.Join(p, "testdata"))
}

func TestTardir(t *testing.T) {
	os.Create(path.Join(p, "testdata/before/untar.txt"))
	// defer func() {
	// 	os.RemoveAll(path.Join(p, "/testdata/after/test.tar.gz"))
	// 	os.RemoveAll(path.Join(p, "/testdata/before/untar.txt"))
	// }()

	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestTardir_1",
			args: args{
				path.Join(p, "/testdata/before/untar.txt"),
				path.Join(p, "/testdata/after/test.tar.gz"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Tardir(tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Tardir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUntar(t *testing.T) {
	type args struct {
		archive string
		dest    string
		force   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestUntar_1",
			args: args{
				path.Join(p, "/testdata/after/test.tar.gz"),
				path.Join(p, "/testdata/before"),
				true,
			},
		},
		{
			name: "TestUntar_2",
			args: args{
				path.Join(p, "/testdata/after/test.tar.gz"),
				path.Join(p, "/testdata/before"),
				false,
			},
			wantErr: true,
		},
		{
			name: "TestUntar_3",
			args: args{
				path.Join(p, "/testdata/after/test.tar.gz"),
				path.Join(p, "/testdata/before2"),
				false,
			},
			wantErr: false,
		},
		{
			name: "TestUntar_4",
			args: args{
				path.Join(p, "/testdata/after/test2.tar.gz"),
				path.Join(p, "/testdata/before2"),
				false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Untar(tt.args.archive, tt.args.dest, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("Untar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUntargz(t *testing.T) {
	Tardir(path.Join(p, "/testdata/before"), path.Join(p, "/testdata/after/dirgz.tar.gz"))
	type args struct {
		archive  string
		dest     string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestUntargz_1",
			args: args{
				path.Join(p, "/testdata/after/dirgz.tar.gz"),
				path.Join(p, "/testdata/before"),
				path.Join(p, "before"),
			},
			wantErr: false,
		},
		{
			name: "TestUntargz_2",
			args: args{
				path.Join(p, "/testdata/after/dirgz.tar.gz"),
				path.Join(p, "/testdata/before3"),
				"before3",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Untargz(tt.args.archive, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Untar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUntargzWithName(t *testing.T) {
	defer func() {
		clear()
	}()
	Tardir(path.Join(p, "/testdata/before"), path.Join(p, "/testdata/after/dir.tar.gz"))
	type args struct {
		archive  string
		dest     string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestUntargzWithName_1",
			args: args{
				path.Join(p, "/testdata/after/dir.tar.gz"),
				path.Join(p, "/testdata/before"),
				path.Join(p, "before"),
			},
			wantErr: true,
		},
		{
			name: "TestUntargzWithName_2",
			args: args{
				path.Join(p, "/testdata/after/dir.tar.gz"),
				path.Join(p, "/testdata/before4"),
				"before4",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UntargzWithName(tt.args.archive, tt.args.dest, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("Untar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

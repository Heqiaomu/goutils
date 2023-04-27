package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func TestHTTPDoDownloadFile(t *testing.T) {
	type args struct {
		filepath string
		url      string
	}
	defer func() {
		_ = os.Remove("./a.txt")
		_ = os.Remove("./b.txt")
	}()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test download file, file not exist",
			args: args{
				filepath: "./a.txt",
				url:      "http://console.develop.blocface.baas.hyperchain.cn/gateway/core/api/v1.0/services/download/agents/agent.sh",
			},
			wantErr: false,
		},
		{
			name: "test download file, file exist",
			args: args{
				filepath: "./b.txt",
				url:      "http://console.develop.blocface.baas.hyperchain.cn/gateway/core/api/v1.0/services/download/agents/agent.sh",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *http.Client
			p := gomonkey.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *http.Client, _ *http.Request) (*http.Response, error) {
				t.Log("mock response")
				return &http.Response{
					Status:     "200 OK",
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("nihao"))),
				}, nil
			})
			defer p.Reset()

			err := HTTPDoDownloadFile(tt.args.filepath, tt.args.url)
			if tt.wantErr {
				assert.NotNil(t, err)
				t.Log(err.Error())
			} else {
				assert.Nil(t, err)
				file, err := ioutil.ReadFile(tt.args.filepath)
				assert.Nil(t, err)
				assert.True(t, strings.Contains(string(file), "nihao"))
			}
		})
	}
}

//
//func TestUploadBigFile(t *testing.T) {
//	//f := func(t *testing.T) *gomonkey.Patches {
//	//	patches := gomonkey.NewPatches()
//	//	var c *http.Client
//	//	p := patches.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *http.Client, _ *http.Request) (*http.Response, error) {
//	//		t.Log("mock response")
//	//		return &http.Response{
//	//			Status:     "200 OK",
//	//			StatusCode: http.StatusOK,
//	//			Body:       ioutil.NopCloser(bytes.NewReader([]byte("nihao"))),
//	//		}, nil
//	//	})
//	//	defer p.Reset()
//	//	return patches
//	//}
//	t.Run("test upload big file", func(t *testing.T) {
//		//patches := f(t)
//		//defer patches.Reset()
//		err := UploadBigFile("httputil.go", "http://console.test.blocface.baas.hyperchain.cn/gateway/core/api/v1.0/services/storeFile/httputil.go")
//		assert.NotNil(t, err)
//	})
//}

func TestHTTPDoGet(t *testing.T) {
	mm := make([]map[string]string, 1)
	mm[0] = make(map[string]string)
	mm[0]["key"] = "value"
	type args struct {
		url     string
		headers []map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"TestHTTPDoGet_1",
			args{
				"https://www.hyperchain.cn",
				nil,
			},
			[]byte("mockdata"),
			false,
		},
		{
			"TestHTTPDoGet_2",
			args{
				"https://1.1.1.1",
				mm,
			},
			[]byte("mockdata"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HTTPDoGet(tt.args.url, tt.args.headers...)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPDoGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("HTTPDoGet() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestHTTPDoPost(t *testing.T) {
	mm := make([]map[string]string, 1)
	mm[0] = make(map[string]string)
	mm[0]["key"] = "value"
	type args struct {
		body    interface{}
		url     string
		headers []map[string]string
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{
			"TestHTTPDoPost_1",
			args{
				"this is body",
				"https://www.baidu.cn",
				mm,
			},
			[]byte("mockdata"),
			false,
		},
		{
			"TestHTTPDoPost_2",
			args{
				"this is body",
				"https://www.hyperchain.cn",
				mm,
			},
			[]byte("mockdata"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HTTPDoPost(tt.args.body, tt.args.url, tt.args.headers...)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPDoPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotData, tt.wantData) {
			// 	t.Errorf("HTTPDoPost() = %v, want %v", gotData, tt.wantData)
			// }
		})
	}
}

func TestHTTPDoUploadFile(t *testing.T) {
	type args struct {
		filepath string
		url      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestHTTPDoUploadFile_1",
			args{
				"md5.go",
				"https://hyperchain.cn",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HTTPDoUploadFile(tt.args.filepath, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("HTTPDoUploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

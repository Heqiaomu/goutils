package compress

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnzipWithFileName(t *testing.T) {
	//f := func(t *testing.T) *gomonkey.Patches {
	//	patches := gomonkey.NewPatches()
	//	var c *http.Client
	//	p := patches.ApplyMethod(reflect.TypeOf(c), "Do", func(_ *http.Client, _ *http.Request) (*http.Response, error) {
	//		t.Log("mock response")
	//		return &http.Response{
	//			Status:     "200 OK",
	//			StatusCode: http.StatusOK,
	//			Body:       ioutil.NopCloser(bytes.NewReader([]byte("nihao"))),
	//		}, nil
	//	})
	//	defer p.Reset()
	//	return patches
	//}
	t.Run("test upload big file", func(t *testing.T) {
		//patches := f(t)
		//defer patches.Reset()
		err := Zipdir("../compress", "../a.zip")
		assert.Nil(t, err)
		defer os.RemoveAll("../a.zip")
		err = UnzipWithFileName("../a.zip", "../", "a1")
		assert.Nil(t, err)
		defer os.RemoveAll("../a1")
		os.MkdirAll("zip", os.ModePerm)
		Unzip("../a.zip", "../")
		defer os.RemoveAll("zip")
	})
}

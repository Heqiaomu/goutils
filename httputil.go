//Package util util
package util

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	fp "path/filepath"
	"strings"
)

var (
	client = &http.Client{}
)

// HTTPDoGet http get
func HTTPDoGet(url string, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(headers) != 0 {
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do http get")
	} else if response == nil {
		return nil, errors.New("http response is nil")
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("fail to get, because http response code [%d], data [%s]", response.StatusCode, string(data))
	}
	return data, nil
}

// HTTPDoPost http post
func HTTPDoPost(body interface{}, url string, headers ...map[string]string) (data []byte, err error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return
	}
	if len(headers) != 0 {
		for _, header := range headers {
			for k, v := range header {
				request.Header.Set(k, v)
			}
		}
	}
	res, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "do http post")
	}
	if res == nil {
		return nil, errors.New("http response is nil")
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("fail to post, because http response code [%d], data [%s]", res.StatusCode, string(data))
	}
	return data, nil
}

// HTTPDoUploadFile http upload file
func HTTPDoUploadFile(filepath string, url string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	res, err := http.Post(url, "binary/octet-stream", file) // 第二个参数用来指定 "Content-Type"
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code : %d", res.StatusCode)
	}
	return nil
}

// HTTPDoDownloadFile download file through http
// saveFilepath : path + target file name, if target file not exist, create it
// url : remote download file url
func HTTPDoDownloadFile(saveFilepath, url string) error {
	filepath, err := fp.Abs(saveFilepath)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(req)
	if err != nil {
		return err
	} else if response == nil {
		return errors.New("http response is nil")
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("http response code : %d", response.StatusCode)
	}
	defer response.Body.Close()
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return err
	}
	_ = f.Sync()
	return nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// UploadBigFile upload big file
func UploadBigFile(file, url string) error {
	fph, err := fp.Abs(file)
	if err != nil {
		return errors.Wrap(err, "get absolute path")
	}
	pr, pw := io.Pipe()
	bw := multipart.NewWriter(pw)

	open, err := os.Open(fph)
	if err != nil {
		return errors.Wrapf(err, "open file [%s]", fph)
	}
	go func() {
		defer open.Close()
		_, fileName := fp.Split(file)
		fw1, err := bw.CreateFormFile("file", fileName)
		if err != nil {
			panic(err)
		}
		//h := make(textproto.MIMEHeader)
		//h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes(fileName)))
		//h.Set("Content-Type", "text/plain")
		//fw1, err := writer.CreatePart(h)
		//if err != nil {
		//	panic(err)
		//}
		_, _ = io.Copy(fw1, open)
		_ = bw.Close()
		_ = pw.Close()
	}()
	request, err := http.NewRequest(http.MethodPost, url, pr)
	if err != nil {
		return errors.Wrap(err, "new post request")
	}
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "do http request")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read response body")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to check http code, because http response code : %d data [%s]. please check server", resp.StatusCode, data)
	}
	return nil
}

package metrics

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type CustomResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.header
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	// implement it as per your requirement
	return 0, nil
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *CustomResponseWriter) Flush() {
}
func (w *CustomResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func NewCustomResponseWriter() *CustomResponseWriter {
	return &CustomResponseWriter{
		header: http.Header{},
	}
}

func TestNewResponseWriter(t *testing.T) {
	w := NewCustomResponseWriter()
	type args struct {
		w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want *responseWriter
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				w: w,
			},
			want: &responseWriter{w, http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponseWriter(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponseWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_responseWriter_WriteHeader(t *testing.T) {
	w := NewCustomResponseWriter()
	type fields struct {
		ResponseWriter http.ResponseWriter
		statusCode     int
	}
	type args struct {
		code int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			},
			args: args{code: http.StatusOK},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &responseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				statusCode:     tt.fields.statusCode,
			}
			rw.WriteHeader(tt.args.code)
		})
	}
}

func Test_responseWriter_Flush(t *testing.T) {
	w := NewCustomResponseWriter()
	type fields struct {
		ResponseWriter http.ResponseWriter
		statusCode     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &responseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				statusCode:     tt.fields.statusCode,
			}
			rw.Flush()
		})
	}
}

func Test_responseWriter_Hijack(t *testing.T) {
	w := NewCustomResponseWriter()
	type fields struct {
		ResponseWriter http.ResponseWriter
		statusCode     int
	}
	tests := []struct {
		name    string
		fields  fields
		want    net.Conn
		want1   *bufio.ReadWriter
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := &responseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				statusCode:     tt.fields.statusCode,
			}
			got, got1, err := rw.Hijack()
			if (err != nil) != tt.wantErr {
				t.Errorf("responseWriter.Hijack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("responseWriter.Hijack() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("responseWriter.Hijack() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPrometheusMiddleware(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	PrometheusMiddleware(next).ServeHTTP(rr, r)

	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

//Package util util
package util

import (
	"crypto/rsa"
	"testing"
)

func TestGeneratePublicKey(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want *rsa.PublicKey
	}{
		{
			"TestGeneratePublicKey_1",
			args{
				"key.pub",
			},
			&rsa.PublicKey{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := GeneratePublicKey(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("GeneratePublicKey() = %v, want %v", got, tt.want)
			// }
			t.Logf("GeneratePublicKey result: %+v", GeneratePublicKey(tt.args.filename))
		})
	}
}

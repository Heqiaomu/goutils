package util

import (
	"testing"
)

func TestGetIntranetIp(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			"TestGetIntranetIp",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetIntranetIp()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIntranetIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if got != tt.want {
			// 	t.Errorf("GetIntranetIp() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestGetUniqueIDFromMac(t *testing.T) {
	tests := []struct {
		name       string
		wantUnique string
		wantErr    bool
	}{
		{
			"TestGetIntranetIp",
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetUniqueIDFromMac()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUniqueIDFromMac() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if gotUnique != tt.wantUnique {
			// 	t.Errorf("GetUniqueIDFromMac() = %v, want %v", gotUnique, tt.wantUnique)
			// }
		})
	}
}

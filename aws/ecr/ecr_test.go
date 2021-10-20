package ecr

import (
	"reflect"
	"testing"
)

func TestListECR1(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListECR()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListECR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListECR() got = %v, want %v", got, tt.want)
			}
		})
	}
}

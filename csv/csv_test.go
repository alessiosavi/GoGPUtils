package csv

import (
	"reflect"
	"testing"
)

func TestReadCsv(t *testing.T) {
	type args struct {
		buf       []byte
		separator rune
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 [][]string
	}{
		{
			name: "testOK",
			args: args{
				buf: []byte(`HEADER1,HEADER2,HEADER3
data1,data2,data3
data4,data5,data6`),
				separator: ',',
			},
			want:  []string{"HEADER1", "HEADER2", "HEADER3"},
			want1: [][]string{{"data1", "data2", "data3"}, {"data4", "data5", "data6"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ReadCsv(tt.args.buf, tt.args.separator)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadCsv() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReadCsv() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

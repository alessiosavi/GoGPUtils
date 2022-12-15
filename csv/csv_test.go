package csvutils

import (
	processingutils "github.com/alessiosavi/GoGPUtils/files/processing"
	"os"
	"reflect"
	"testing"
)

func TestReadCSV(t *testing.T) {
	type args struct {
		buf       []byte
		separator rune
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want1   [][]string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "testOK",
			args: args{
				buf: []byte(`HEADER1,HEADER2,HEADER3
data1,data2,data3
data4,data5,data6`),
				separator: ',',
			},
			want:    []string{"HEADER1", "HEADER2", "HEADER3"},
			want1:   [][]string{{"data1", "data2", "data3"}, {"data4", "data5", "data6"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ReadCSV(tt.args.buf, tt.args.separator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadCSV() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReadCSV() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUTF(t *testing.T) {
	file, err := os.ReadFile("/tmp/Shipment_20220216.csv")
	if err != nil {
		panic(err)
	}
	file, err = processingutils.ToUTF8(file)
	if err != nil {
		return
	}
	headers, csvData, err := ReadCSV(file, ';')
	if err != nil {
		panic(err)
	}
	csvf, err := WriteCSV(headers, csvData, ';')
	if err != nil {
		panic(err)
	}
	os.WriteFile("/tmp/Shipment_20220216_test.csv", csvf, 0755)
}

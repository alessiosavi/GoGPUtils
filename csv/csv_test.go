package csvutils

import (
	"log"
	"slices"
	"testing"
)

func TestApply(t *testing.T) {
	fn := ApplyCSV(func(i int, v []string) []string {
		if i%2 == 0 {
			return append(v, "NINJA_PARI")
		} else {
			return append(v, "NINJA_DISPARI")
		}
	})

	type args struct {
		headers Headers
		csvData CSVData
	}
	tests := []struct {
		name     string
		args     args
		expected args
	}{
		{
			name: "OK1",
			args: args{
				headers: nil,
				csvData: [][]string{
					{
						"ROW_1_COL1", "ROW_1_COL2", "ROW_1_COL3",
					},
					{
						"ROW_2_COL1", "ROW_2_COL2", "ROW_2_COL3",
					},
				},
			},
			expected: args{
				headers: nil,
				csvData: [][]string{
					{
						"ROW_1_COL1", "ROW_1_COL2", "ROW_1_COL3", "NINJA_PARI",
					},
					{
						"ROW_2_COL1", "ROW_2_COL2", "ROW_2_COL3", "NINJA_DISPARI",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		res := tt.args.csvData.Apply(fn, false)
		for i := range res {
			if !slices.Equal(res[i], tt.expected.csvData[i]) {
				log.Println("INPLACE -> FALSE")
				log.Printf("EXPECTED:\n%v\n", tt.expected.csvData)
				log.Printf("FOUND:\n%v\n", res)
				t.Error()
			}

		}
		tt.args.csvData.Apply(fn, true)
		for i := range res {
			if !slices.Equal(tt.args.csvData[i], tt.expected.csvData[i]) {
				log.Println("INPLACE -> TRUE")
				log.Printf("EXPECTED:\n%v\n", tt.expected.csvData[i])
				log.Printf("FOUND:\n%v\n", tt.args.csvData[i])
				t.Error()
			}
		}
	}

}

func TestExplode(t *testing.T) {
	fn := ExplodeCSV(func(i int, strings []string) [][]string {
		var res [][]string
		if i%2 == 0 {
			res = append(res, append(strings, []string{"odd"}...))
		} else {
			res = append(res, append(strings, []string{"even1"}...))
			res = append(res, append(strings, []string{"even2"}...))
		}
		return res
	})

	type args struct {
		headers Headers
		csvData CSVData
	}
	tests := []struct {
		name     string
		args     args
		expected args
	}{
		{
			name: "OK1",
			args: args{
				headers: nil,
				csvData: [][]string{
					{"ROW_1_COL1", "ROW_1_COL2", "ROW_1_COL3"},
					{"ROW_2_COL1", "ROW_2_COL2", "ROW_2_COL3"},
				},
			},
			expected: args{
				headers: nil,
				csvData: [][]string{
					{"ROW_1_COL1", "ROW_1_COL2", "ROW_1_COL3", "odd"},
					{"ROW_2_COL1", "ROW_2_COL2", "ROW_2_COL3", "even1"},
					{"ROW_2_COL1", "ROW_2_COL2", "ROW_2_COL3", "even2"},
				},
			},
		},
	}

	for _, tt := range tests {
		res := tt.args.csvData.Explode(fn)
		for i := range res {
			if !slices.Equal(res[i], tt.expected.csvData[i]) {
				log.Println("INPLACE -> FALSE")
				log.Printf("EXPECTED:\n%v\n", tt.expected.csvData[i])
				log.Printf("FOUND:\n%v\n", res[i])
				t.Error()
			}
		}
	}

}

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
			got, got1, err := ReadCSV(tt.args.buf, tt.args.separator, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("ReadCSV() got = %v, want %v", got, tt.want)
			}
			for i := range got1 {
				if !slices.Equal(got1[i], tt.want1[i]) {
					t.Errorf("ReadCSV() got1 = %v, want %v", got1, tt.want1)
				}

			}

		})
	}
}

package sqs

import "testing"

func TestGetMessage(t *testing.T) {
	type args struct {
		queueName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				queueName: "qa-thb-sqs-data-input",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GetMessage(tt.args.queueName); (err != nil) != tt.wantErr {
				t.Errorf("GetMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

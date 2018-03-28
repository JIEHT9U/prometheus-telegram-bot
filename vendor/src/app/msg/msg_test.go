package msg

import (
	"io"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Alerts
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Parser(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Parser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

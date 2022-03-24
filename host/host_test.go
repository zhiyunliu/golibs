package host

import "testing"

func TestExtract(t *testing.T) {
	type args struct {
		hostPort string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "1. localhost ", args: args{hostPort: "localhost:1001"}, want: "::1001", wantErr: false},
		{name: "2. : ", args: args{hostPort: "::1001"}, want: "::1001", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Extract(tt.args.hostPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}

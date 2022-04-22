package mrt

import (
	"net/http"
	"testing"
)

func TestMRTService_GetUpcomingTimeTable(t *testing.T) {
	type fields struct {
		client http.Client
	}
	type args struct {
		station     string
		destination string
		number      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{http.Client{}},
			args:    args{"景美", "松山", 3},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MRTService{
				client: tt.fields.client,
			}
			got, err := s.GetUpcomingTimeTable(tt.args.station, tt.args.destination, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("MRTService.GetUpcomingTimeTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("got:", got)
		})
	}
}

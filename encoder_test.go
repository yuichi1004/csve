package csve

import (
	"bytes"
	"encoding/csv"
	"testing"
	"time"
)

func TestNewEncoder(t *testing.T) {
	// skip
}

func TestEncoder_Encode(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	type fields struct {
		Location *time.Location
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal case",
			fields: fields{
				Location: time.UTC,
			},
			args: args{
				&TestData{"str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Unix(1514129400, 0)},
			},
			want: "str,-1,-32,-64,32,64,3.2,6.4,2017-12-24T15:30:00\n",
		},
		{
			name: "normal case with locale",
			fields: fields{
				Location: loc,
			},
			args: args{
				&TestData{"str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Unix(1514097000, 0)},
			},
			want: "str,-1,-32,-64,32,64,3.2,6.4,2017-12-24T15:30:00\n",
		},
		{
			name: "pointer case",
			fields: fields{
				Location: time.UTC,
			},
			args: args{
				NewTestDataPtr("str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Unix(1514129400, 0)),
			},
			want: "str,-1,-32,-64,32,64,3.2,6.4,2017-12-24T15:30:00\n",
		},
		{
			name: "pointer nil case",
			fields: fields{
				Location: time.UTC,
			},
			args: args{
				TestDataPtr{},
			},
			want: ",,,,,,,,\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			e := &Encoder{
				Writer:   csv.NewWriter(buf),
				Location: tt.fields.Location,
			}
			if err := e.Encode(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Encoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
			e.Flush()
			if buf.String() != tt.want {
				t.Errorf("Encode() = %v, want %v", buf.String(), tt.want)
			}
		})
	}
}

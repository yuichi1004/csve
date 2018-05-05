package csve

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewEncoder(t *testing.T) {
	// skip
}

func TestEncoder_Encode(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	type fields struct {
		Location      *time.Location
		CustomEncoder CustomEncoder
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
			name: "normal case with custom encoder",
			fields: fields{
				CustomEncoder: func(e *Encoder, v reflect.Value, format string) (ok bool, raw string, err error) {
					val := v.Interface()
					var t time.Time
					if t, ok = val.(time.Time); ok {
						if t.IsZero() {
							return true, "N/A", nil
						}
					}
					return false, "", nil
				},
			},
			args: args{
				&TestData{"str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Time{}},
			},
			want: "str,-1,-32,-64,32,64,3.2,6.4,N/A\n",
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
				Writer:        csv.NewWriter(buf),
				Location:      tt.fields.Location,
				CustomEncoder: tt.fields.CustomEncoder,
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

func ExampleEncoder_Encode() {
	v := struct {
		ID      int       `csv:"0,id"`
		Name    string    `csv:"1,name"`
		Created time.Time `csv:"2,created,2006-01-02T15:04:05"`
	}{
		5, "Yuichi", time.Unix(1514129400, 0),
	}

	csvWriter := csv.NewWriter(os.Stdout)
	encoder, err := NewEncoder(csvWriter, false)
	if err != nil {
		fmt.Printf("failed to create encoder: %v\n", err)
		return
	}

	if err := encoder.Encode(&v); err != nil {
		fmt.Printf("failed to parse csv: %v\n", err)
		return
	}
	encoder.Flush()

	// Output: 5,Yuichi,2017-12-24T15:30:00
}

func ExampleCustomEncoder() {
	v := struct {
		ID      int       `csv:"0,id"`
		Name    string    `csv:"1,name"`
		Created time.Time `csv:"2,created,2006-01-02T15:04:05"`
	}{
		5, "Yuichi", time.Time{},
	}

	csvWriter := csv.NewWriter(os.Stdout)
	encoder, err := NewEncoder(csvWriter, false)
	if err != nil {
		fmt.Printf("failed to create encoder: %v\n", err)
		return
	}
	encoder.CustomEncoder = func(e *Encoder, v reflect.Value, format string) (ok bool, raw string, err error) {
		val := v.Interface()
		var t time.Time
		if t, ok = val.(time.Time); ok {
			if t.IsZero() {
				return true, "N/A", nil
			}
		}
		return false, "", nil
	}

	if err := encoder.Encode(&v); err != nil {
		fmt.Printf("failed to parse csv: %v\n", err)
		return
	}
	encoder.Flush()

	// Output: 5,Yuichi,N/A
}

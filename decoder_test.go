package csve

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
	"time"
)

type TestData struct {
	Str     string    `csv:"0,str"`
	Int     int       `csv:"1,int"`
	Int32   int32     `csv:"2,int32"`
	Int64   int64     `csv:"3,int64"`
	Uint32  uint32    `csv:"4,uint32"`
	Uint64  uint64    `csv:"5,uint64"`
	Float32 float32   `csv:"6,float32"`
	Float64 float64   `csv:"7,float64"`
	Time    time.Time `csv:"8,time,2006-01-02T15:04:05Z07:00"`
}

func TestMain(t *testing.T) {
}

func TestNewDecoder(t *testing.T) {
	type args struct {
		reader    *csv.Reader
		useHeader bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Decoder
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDecoder(tt.args.reader, tt.args.useHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	type fields struct {
		Reader *csv.Reader
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TestData
		wantErr bool
	}{
		{
			name: "normal case",
			fields: fields{
				Reader: csv.NewReader(strings.NewReader("str,-1,-32,-64,32,64,3.2,6.4,2017-12-24T15:30:00Z")),
			},
			args: args{
				v: &TestData{},
			},
			want: &TestData{
				"str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Unix(1514129400, 0).In(time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{
				Reader: tt.fields.Reader,
			}
			if err := d.Decode(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("Decode() = %v, want %v", tt.args.v, tt.want)
			}
		})
	}
}

func Test_getFields(t *testing.T) {

	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name       string
		args       args
		wantFields []field
		wantErr    bool
	}{
		{
			name: "normal case",
			args: args{
				reflect.TypeOf(&TestData{}),
			},
			wantFields: []field{
				{
					typ:        reflect.TypeOf(""),
					fieldindex: []int{0},
					fieldname:  "Str",
					csvname:    "str",
					csvindex:   0,
				},
				{
					typ:        reflect.TypeOf(int(0)),
					fieldindex: []int{1},
					fieldname:  "Int",
					csvname:    "int",
					csvindex:   1,
				},
				{
					typ:        reflect.TypeOf(int32(0)),
					fieldindex: []int{2},
					fieldname:  "Int32",
					csvname:    "int32",
					csvindex:   2,
				},
				{
					typ:        reflect.TypeOf(int64(0)),
					fieldindex: []int{3},
					fieldname:  "Int64",
					csvname:    "int64",
					csvindex:   3,
				},
				{
					typ:        reflect.TypeOf(uint32(0)),
					fieldindex: []int{4},
					fieldname:  "Uint32",
					csvname:    "uint32",
					csvindex:   4,
				},
				{
					typ:        reflect.TypeOf(uint64(0)),
					fieldindex: []int{5},
					fieldname:  "Uint64",
					csvname:    "uint64",
					csvindex:   5,
				},
				{
					typ:        reflect.TypeOf(float32(0)),
					fieldindex: []int{6},
					fieldname:  "Float32",
					csvname:    "float32",
					csvindex:   6,
				},
				{
					typ:        reflect.TypeOf(float64(0)),
					fieldindex: []int{7},
					fieldname:  "Float64",
					csvname:    "float64",
					csvindex:   7,
				},
				{
					typ:        reflect.TypeOf(time.Time{}),
					fieldindex: []int{8},
					fieldname:  "Time",
					csvname:    "time",
					csvindex:   8,
					csvformat:  "2006-01-02T15:04:05Z07:00",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFields, err := getFields(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := 0; i < len(gotFields); i++ {
				// ignore decoder to compare
				gotFields[i].dec = nil
			}
			if !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Errorf("getFields() = %v, want %v", gotFields, tt.wantFields)
			}
		})
	}
}

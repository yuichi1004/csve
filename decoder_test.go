package csve

import (
	"encoding/csv"
	"fmt"
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
	Time    time.Time `csv:"8,time,2006-01-02T15:04:05"`
}

type TestDataPtr struct {
	Str     *string    `csv:"0,str"`
	Int     *int       `csv:"1,int"`
	Int32   *int32     `csv:"2,int32"`
	Int64   *int64     `csv:"3,int64"`
	Uint32  *uint32    `csv:"4,uint32"`
	Uint64  *uint64    `csv:"5,uint64"`
	Float32 *float32   `csv:"6,float32"`
	Float64 *float64   `csv:"7,float64"`
	Time    *time.Time `csv:"8,time,2006-01-02T15:04:05"`
}

func normalizeTestData(v interface{}) interface{} {
	switch cv := v.(type) {
	case *TestData:
		v = cv
	case *TestDataPtr:
		td := new(TestData)
		if cv.Str != nil {
			td.Str = *cv.Str
		}
		if cv.Int != nil {
			td.Int = *cv.Int
		}
		if cv.Int32 != nil {
			td.Int32 = *cv.Int32
		}
		if cv.Int64 != nil {
			td.Int64 = *cv.Int64
		}
		if cv.Uint32 != nil {
			td.Uint32 = *cv.Uint32
		}
		if cv.Uint64 != nil {
			td.Uint64 = *cv.Uint64
		}
		if cv.Float32 != nil {
			td.Float32 = *cv.Float32
		}
		if cv.Float64 != nil {
			td.Float64 = *cv.Float64
		}
		if cv.Time != nil {
			td.Time = *cv.Time
		}
		v = td
	}
	return v
}

func TestNewDecoder(t *testing.T) {
	// skip
}

func TestDecoder_Decode(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	type fields struct {
		Reader   *csv.Reader
		Location *time.Location
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
				Reader:   csv.NewReader(strings.NewReader("str,-1,-32,-64,32,64,3.2,6.4,2017-12-24T15:30:00")),
				Location: time.UTC,
			},
			args: args{
				v: &TestData{},
			},
			want: &TestData{
				"str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Unix(1514129400, 0).In(time.UTC),
			},
		},
		{
			name: "normal case with locale",
			fields: fields{
				Reader:   csv.NewReader(strings.NewReader("str,-1,-32,-64,32,64,3.2,6.4,2017-12-24T15:30:00")),
				Location: loc,
			},
			args: args{
				v: &TestData{},
			},
			want: &TestData{
				"str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Unix(1514097000, 0).In(loc),
			},
		},
		{
			name: "pointer case",
			fields: fields{
				Reader:   csv.NewReader(strings.NewReader("str,-1,-32,-64,32,64,3.2,6.4,2017-12-24T15:30:00")),
				Location: time.UTC,
			},
			args: args{
				v: &TestDataPtr{},
			},
			want: &TestData{
				"str", -1, -32, -64, 32, 64, 3.2, 6.4, time.Unix(1514129400, 0).In(time.UTC),
			},
		},
		{
			name: "pointer nil case",
			fields: fields{
				Reader:   csv.NewReader(strings.NewReader(",,,,,,,,2017-12-24T15:30:00")),
				Location: time.UTC,
			},
			args: args{
				v: &TestDataPtr{},
			},
			want: &TestData{
				"", 0, 0, 0, 0, 0, 0, 0, time.Unix(1514129400, 0).In(time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decoder{
				Reader:   tt.fields.Reader,
				Location: tt.fields.Location,
			}
			if err := d.Decode(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(normalizeTestData(tt.args.v), tt.want) {
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
					csvformat:  "2006-01-02T15:04:05",
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

func ExampleDecoder_Decode() {
	v := struct {
		ID      int       `csv:"0,id"`
		Name    string    `csv:"1,name"`
		Created time.Time `csv:"2,created,2006-01-02T15:04:05"`
	}{}

	csvReader := csv.NewReader(strings.NewReader(`5,Yuichi,2017-12-24T15:30:00`))
	decoder, err := NewDecoder(csvReader, false)
	if err != nil {
		fmt.Printf("failed to create decoder: %v\n", err)
		return
	}

	if err := decoder.Decode(&v); err != nil {
		fmt.Printf("failed to parse csv: %v\n", err)
		return
	}
	fmt.Printf("ID:%d, Name:%s, Created:%v\n", v.ID, v.Name, v.Created)

	// Output: ID:5, Name:Yuichi, Created:2017-12-24 15:30:00 +0000 UTC
}

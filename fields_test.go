package csve

import (
	"reflect"
	"testing"
	"time"
)

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

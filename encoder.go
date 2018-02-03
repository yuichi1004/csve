package csve

import (
	"encoding/csv"
	"reflect"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

// Encoder writes values into csv writer.
type Encoder struct {
	*csv.Writer

	// Spcify location to be used encoding. If not specified, Encoder use time.UTC.
	Location *time.Location
}

// NewEncoder returns a new Encoder which encodes values into csv writer.
// If useHeader is true, Encoder writes csv header line.
// NOTE: useHeder is not implemented yet.
func NewEncoder(writer *csv.Writer, useHeader bool) (*Encoder, error) {
	return &Encoder{
		writer, time.UTC,
	}, nil
}

// Encode encodes value into csv writer.
func (e *Encoder) Encode(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return errors.New("invalid value type")
		}
		rv = rv.Elem()
	}

	fields, err := getFields(rv.Type())
	if err != nil {
		return err
	}

	encoded := make([]string, len(fields))
	for i, f := range fields {
		ref := rv.FieldByIndex(f.fieldindex)
		var err error
		encoded[i], err = f.enc(e, ref, f.csvformat)
		if err != nil {
			return errors.Errorf("field %s encode failed", f.fieldname)
		}
	}

	return e.Writer.Write(encoded)
}
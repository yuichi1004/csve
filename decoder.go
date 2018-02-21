package csve

import (
	"encoding/csv"
	"reflect"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

var (
	timeType = reflect.TypeOf(time.Time{})
)

// Decoder reads csv lines from upstream reader and decode the line.
type Decoder struct {
	*csv.Reader

	// Spcify location to be used decoding. If not specified, Decoder use time.UTC.
	Location *time.Location

	line int
}

// NewDecoder returns a NewDecoder which decodes values from reader.
// If useHeader is true, NewDecoder reads the header line and decode values
// based on the header.
// NOTE: useHader is not implemented yet.
func NewDecoder(reader *csv.Reader, useHeader bool) (*Decoder, error) {
	if useHeader {
		// TODO impelemnt header usage
	}

	return &Decoder{
		reader, time.UTC, 0,
	}, nil
}

// Decode reads csv line and decode values into v.
func (d *Decoder) Decode(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid value type")
	}

	fields, err := getFields(rv.Type())
	if err != nil {
		return err
	}

	cols, err := d.Read()
	if err != nil {
		return err
	}
	d.line++

	for _, f := range fields {
		ref := rv.Elem().FieldByIndex(f.fieldindex)

		var v string
		if f.csvindex < len(cols) {
			v = cols[f.csvindex]
		}

		if err := f.dec(d, ref, v, f.csvformat); err != nil {
			return errors.Errorf("field %s parse failed (line:%d)", f.fieldname, d.line)
		}
	}

	return nil
}

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

// CustomerDecoder inject your custom decode process.
// Return true as ok if this logic handles the decode, otherwise Decoder() will
// fallback to default decode process.
type CustomDecoder func(d *Decoder, v reflect.Value, raw, format string) (ok bool, err error)

// Decoder reads csv lines from upstream reader and decode the line.
type Decoder struct {
	*csv.Reader

	// Spcify location to be used decoding. If not specified, Decoder use time.UTC.
	Location *time.Location

	// Custom decoder to customize decoding process.
	CustomDecoder CustomDecoder

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
		reader, time.UTC, nil, 0,
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

		var ok bool
		if d.CustomDecoder != nil {
			ok, err = d.CustomDecoder(d, ref, v, f.csvformat)
			if err != nil {
				return errors.Errorf("field %s parse failed (line:%d)", f.fieldname, d.line)
			}
		}
		if !ok {
			if err := f.dec(d, ref, v, f.csvformat); err != nil {
				return errors.Errorf("field %s parse failed (line:%d)", f.fieldname, d.line)
			}
		}
	}

	return nil
}

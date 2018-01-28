package csve

import (
	"encoding/csv"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	timeType = reflect.TypeOf(time.Time{})
)

var fieldCache sync.Map

// Decoder reads csv lines from upstream reader and decode the line.
type Decoder struct {
	*csv.Reader

	// Spcify location to be used decoding. If not specified, Decoer use time.UTC.
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
		if err := f.dec(d, ref, cols[f.csvindex], f.csvformat); err != nil {
			return errors.Errorf("field %s parse failed (line:%d)", f.fieldname, d.line)
		}
	}

	return nil
}

type fieldDecoder func(d *Decoder, v reflect.Value, raw, format string) error

func intDecoder(d *Decoder, v reflect.Value, raw, format string) error {
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v.OverflowInt(n) {
		return err
	}
	v.SetInt(n)
	return nil
}

func uintDecoder(d *Decoder, v reflect.Value, raw, format string) error {
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || v.OverflowUint(n) {
		return err
	}
	v.SetUint(n)
	return nil
}

func floatDecoder(d *Decoder, v reflect.Value, raw, format string) error {
	n, err := strconv.ParseFloat(raw, 64)
	if err != nil || v.OverflowFloat(n) {
		return err
	}
	v.SetFloat(n)
	return nil
}

func stringDecoder(d *Decoder, v reflect.Value, raw, format string) error {
	v.SetString(raw)
	return nil
}

func timeDecoder(d *Decoder, v reflect.Value, raw, format string) error {
	t, err := time.ParseInLocation(format, raw, d.Location)
	if err != nil {
		return err
	}
	tv := reflect.ValueOf(t)
	v.Set(tv)
	return nil
}

type field struct {
	typ        reflect.Type
	dec        fieldDecoder
	fieldindex []int
	fieldname  string

	csvname   string
	csvindex  int
	csvformat string
}

func getFields(t reflect.Type) (fields []field, err error) {
	if t.Kind() == reflect.Ptr {
		return getFields(t.Elem())
	}

	if fields, ok := fieldCache.Load(t); ok {
		return fields.([]field), nil
	}

	fields = make([]field, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag, ok := f.Tag.Lookup("csv")
		if !ok {
			continue
		}

		tags := strings.Split(tag, ",")
		index, err := strconv.ParseInt(tags[0], 10, 64)
		if err != nil {
			index = -1
		}

		var format string
		if len(tags) >= 3 {
			format = tags[2]
		}

		var dec fieldDecoder
		dec, err = getFieldDecoder(f.Type)
		if err != nil {
			return nil, err
		}

		fields = append(fields, field{
			dec:        dec,
			typ:        f.Type,
			fieldname:  f.Name,
			fieldindex: f.Index,
			csvname:    tags[1],
			csvindex:   int(index),
			csvformat:  format,
		})
	}

	fieldCache.Store(t, fields)
	return
}

func getFieldDecoder(t reflect.Type) (dec fieldDecoder, err error) {
	switch t.Kind() {
	case reflect.Ptr:
		var rdec fieldDecoder
		rdec, err = getFieldDecoder(t.Elem())
		if err == nil {
			dec = func(d *Decoder, v reflect.Value, raw, format string) error {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
				return rdec(d, v, raw, format)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dec = intDecoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dec = uintDecoder
	case reflect.String:
		dec = stringDecoder
	case reflect.Float32, reflect.Float64:
		dec = floatDecoder
	case reflect.Struct:
		if t == timeType {
			dec = timeDecoder
		} else {
			err = errors.New("no field decoder found")
		}
	default:
		err = errors.New("no field decoder found")
	}
	return
}

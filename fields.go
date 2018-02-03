package csve

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var fieldCache sync.Map

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

type fieldEncoder func(e *Encoder, v reflect.Value, format string) (raw string, err error)

func vEncoder(e *Encoder, v reflect.Value, format string) (raw string, err error) {
	return fmt.Sprintf("%v", v), nil
}

func timeEncoder(e *Encoder, v reflect.Value, format string) (raw string, err error) {
	t := v.Interface().(time.Time)
	return t.In(e.Location).Format(format), nil
}

type field struct {
	typ        reflect.Type
	dec        fieldDecoder
	enc        fieldEncoder
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
		var enc fieldEncoder
		dec, enc, err = getFieldEncoder(f.Type)
		if err != nil {
			return nil, err
		}

		fields = append(fields, field{
			dec:        dec,
			enc:        enc,
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

func getFieldEncoder(t reflect.Type) (dec fieldDecoder, enc fieldEncoder, err error) {
	switch t.Kind() {
	case reflect.Ptr:
		var rdec fieldDecoder
		var renc fieldEncoder
		rdec, renc, err = getFieldEncoder(t.Elem())
		if err == nil {
			dec = func(d *Decoder, v reflect.Value, raw, format string) error {
				if raw == "" {
					v.Set(reflect.Zero(v.Type()))
					return nil
				} else {
					if v.IsNil() {
						v.Set(reflect.New(v.Type().Elem()))
					}
					v = v.Elem()
					return rdec(d, v, raw, format)
				}
			}
			enc = func(e *Encoder, v reflect.Value, format string) (string, error) {
				if v.IsNil() {
					return "", nil
				}
				v = v.Elem()
				return renc(e, v, format)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dec = intDecoder
		enc = vEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dec = uintDecoder
		enc = vEncoder
	case reflect.String:
		dec = stringDecoder
		enc = vEncoder
	case reflect.Float32, reflect.Float64:
		dec = floatDecoder
		enc = vEncoder
	case reflect.Struct:
		if t == timeType {
			dec = timeDecoder
			enc = timeEncoder
		} else {
			err = errors.New("no field decoder found")
		}
	default:
		err = errors.New("no field decoder found")
	}
	return
}

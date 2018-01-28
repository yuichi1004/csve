package csve

import (
	"encoding/csv"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

// decoder benchmark

type benchCsvSrcReader struct {
	data []byte
}

func newBenchCsvReader(c string) *csv.Reader {
	return csv.NewReader(&benchCsvSrcReader{[]byte(c)})
}

func (r *benchCsvSrcReader) Read(p []byte) (int, error) {
	copy(p, r.data)
	return len(r.data), nil
}

func BenchmarkDecode(b *testing.B) {
	type data struct {
		V1 string    `csv:"0,v1"`
		V2 int64     `csv:"1,v2"`
		V3 float64   `csv:"2,v3"`
		V4 time.Time `csv:"3,v4,2006-01-02T15:04:05Z07:00"`
	}
	r := newBenchCsvReader(`"str",1,2.0,2017-12-24T15:30:00Z` + "\n")
	dec, _ := NewDecoder(r, false)

	var d data
	for i := 0; i < b.N; i++ {
		dec.Decode(&d)
	}
}

func BenchmarkRaw(b *testing.B) {
	type data struct {
		V1 string
		V2 int64
		V3 float64
		V4 time.Time
	}
	r := newBenchCsvReader(`"str",1,2.0,2017-12-24T15:30:00Z` + "\n")

	var d data
	for i := 0; i < b.N; i++ {
		v, _ := r.Read()
		d.V1 = v[0]
		{
			v, _ := strconv.ParseInt(v[1], 10, 64)
			d.V2 = v
		}
		{
			v, _ := strconv.ParseFloat(v[2], 64)
			d.V3 = v
		}
		{
			v, err := time.Parse("2006-01-02T15:04:05Z07:00", v[3])
			if err != nil {
				b.Fatalf("faield to parse time: %v", err)
			}
			d.V4 = v
		}
	}
}

func BenchmarkJson(b *testing.B) {
	type data struct {
		V1 string  `json:"v1"`
		V2 int64   `json:"v2"`
		V3 float64 `json:"v3"`
		V4 string  `json:"v4"`
	}
	v := []byte(`{"v1":"str","v2":1,"v3":2.0,"v4":"2017-12-24T15:30:00Z"}`)
	var d data
	for i := 0; i < b.N; i++ {
		json.Unmarshal(v, &d)
		time.Parse("2006-01-02T15:04:05Z07:00", d.V4)
	}
}

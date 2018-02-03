# csve

CSV encode/decoder.

TODO: encoder needs to be implemented

# Decode

```go
type V struct {
    Name       string    `csv:"0,name"`
    ID         int       `csv:"1,id"`
    Registered time.Time `csv:"2,registered,2006-01-02T15:04:05Z07:00"`
}

var reader io.Reader
csvreader := csv.NewReader(csvreader)
decoder, err := csve.NewDecoder(csvreader)
if err != nil {
    panic(err)
}
var v V
if err := decoder.Decode(&v); err != nil {
    panic(err)
}
```

# Benchmark

csve has excellent performance comparing to standard encoding/json decoder.
It just a 50% overhead comparing to raw decoding code.

```
BenchmarkDecode-4        1000000              1345 ns/op             128 B/op          3 allocs/op
BenchmarkEncode-4        1000000              1105 ns/op             176 B/op          6 allocs/op
BenchmarkDecodePtr-4     1000000              1366 ns/op             128 B/op          3 allocs/op
BenchmarkRaw-4           2000000               973 ns/op              96 B/op          2 allocs/op
BenchmarkJson-4          1000000              2112 ns/op             336 B/op          6 allocs/op
```

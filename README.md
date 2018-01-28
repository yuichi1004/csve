# csve

CSV encode/decoder.

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


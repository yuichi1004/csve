package csve

// CsvReader defines interfce for encoding.csv.Reader.
type CsvReader interface {
	Read() (record []string, err error)
}

// CsvWriter defines interfce for encoding.csv.Writer.
type CsvWriter interface {
	Write(record []string) error
	Flush()
}

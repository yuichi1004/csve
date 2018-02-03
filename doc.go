/*
Package csve provide csv encoder and decoder.
csve reads struct tag to encode/decode value likewise encoding/json.

The following is the quick example for decoding csv value.

	v := struct {
		ID      int       `csv:"0,id"`
		Name    string    `csv:"1,name"`
		Created time.Time `csv:"2,created,2006-01-02T15:04:05"`
	}{}

	csvReader := csv.NewReader(strings.NewReader(`5,Yuichi,2017-12-24T15:30:00`))
	decoder, err := NewDecoder(csvReader, false)
	if err != nil {
		fmt.Printf("failed to create decoder: %v\n", err)
		return
	}

	if err := decoder.Decode(&v); err != nil {
		fmt.Printf("failed to parse csv: %v\n", err)
		return
	}
	fmt.Printf("ID:%d, Name:%s, Created:%v\n", v.ID, v.Name, v.Created)

	// Expected Output:
	// ID:5, Name:Yuichi, Created:2017-12-24 15:30:00 +0000 UTC

*/
package csve

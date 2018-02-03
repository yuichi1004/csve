package csve

import (
	"encoding/csv"
	"time"
)

type Encoder struct {
	*csv.Writer

	Location *time.Location
}

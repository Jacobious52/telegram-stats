package importer

import (
	"time"
)

type Table struct {
	Rows []Row
}

type Row struct {
	From string
	Date time.Time
	Text string
}

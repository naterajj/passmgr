package csv

import (
	"encoding/csv"
	"os"
	"passmgr/database"
	"time"
)

func ReadCSV(csvfile string) []database.Record {
	f, err := os.Open(csvfile)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	reader := csv.NewReader(f)
	csvrecords, err := reader.ReadAll()
	if err != nil {
		panic(err.Error())
	}

	records := []database.Record{}
	for i, r := range csvrecords {
		// skip header line
		if i == 0 {
			continue
		}
		now := time.Now().String()
		records = append(records, database.Record{
			0,
			r[0],
			r[1],
			r[2],
			r[3],
			now,
			now,
		})
	}
	return records
}

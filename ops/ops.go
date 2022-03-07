package ops

import (
	"fmt"
	"passmgr/csv"
	"passmgr/database"
	"passmgr/input"
	"strconv"
)

func ImportCSV(csvfile string, skipFirst bool, db database.Database) {
	inserted := db.Insert(csv.ReadCSV(csvfile), skipFirst)
	fmt.Println("Imported " + strconv.Itoa(inserted) + " records.")
}

func Search(host string, db database.Database) {
	records := db.Search(host)
	for _, r := range records {
		printRecord(r)
	}
}

func Update(id uint, db database.Database) {
	r := db.SearchByID(id)
	fmt.Println("Current record")
	printRecord(r)
	fmt.Println("Enter updated record.")
	nr := inputRecord()
	db.Update(id, nr)
	fmt.Println("Record updated")
}

func Insert(db database.Database) {
	r := inputRecord()
	record := []database.Record{r}
	inserted := db.Insert(record, false)
	if inserted == 1 {
		fmt.Println("Password added")
	}
}

func inputRecord() database.Record {
	r := database.Record{}
	r.Host = input.Read("Host: ")
	r.URL = input.Read("URL: ")
	r.Username = input.Read("Username: ")
	for {
		password, err := input.ReadSecret("Password (it won't echo)")
		if err == nil {
			r.Password = password
			break
		}
	}
	return r
}

func printRecord(r database.Record) {
	fmt.Printf("%d: %s %s\n", r.ID, r.Host, r.URL)
	fmt.Println("username: " + r.Username)
	fmt.Println("password: " + r.Password)
	fmt.Println("Last Updated: " + r.UpdatedAt)
}

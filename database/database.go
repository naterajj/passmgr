package database

import (
	"context"
	"database/sql"
	"fmt"
	"passmgr/symcrypto"
	"time"

	_ "modernc.org/sqlite"
)

type Database struct {
	db         *sql.DB
	passphrase string
	ctx        context.Context
}

type Record struct {
	ID        uint
	Host      string
	URL       string
	Username  string
	Password  string
	CreatedAt string
	UpdatedAt string
}

type SearchParams struct {
	Host     string
	Username string
}

const (
	insert_record  string = "INSERT INTO passwords (host, url, username, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	update_record  string = "UPDATE passwords SET password = ?, updated_at = ? WHERE id = ?"
	search_records string = "SELECT id, host, url, username, password, created_at, updated_at FROM passwords WHERE host LIKE ? ORDER by updated_at DESC"
	search_by_id   string = "SELECT id, host, url, username, password, created_at, updated_at FROM passwords WHERE ID = ?"
	create_table   string = `
CREATE TABLE IF NOT EXISTS passwords (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  host       TEXT NOT NULL,
  url        TEXT NOT NULL,
  username   TEXT,
  password   BLOB NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);`
	check_table string = "SELECT COUNT(name) FROM sqlite_master WHERE name='passwords'"
)

func NewDatabase(DBFile string, passphrase string) Database {
	db, err := sql.Open("sqlite", DBFile)
	if err != nil {
		panic(err)
	}

	database := Database{
		db:         db,
		passphrase: passphrase,
	}

	return database
}

func (d Database) Search(host string) []Record {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := d.db.PrepareContext(ctx, search_records)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query("%" + host)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var records []Record
	for rows.Next() {
		rec := Record{}
		var ciphertext sql.RawBytes
		if err := rows.Scan(&rec.ID, &rec.Host, &rec.URL, &rec.Username, &ciphertext, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			panic(err)
		}
		rec.Password = symcrypto.Decrypt(ciphertext, d.passphrase)
		records = append(records, rec)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return records
}

func (d Database) SearchByID(id uint) Record {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := d.db.PrepareContext(ctx, search_by_id)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	r := Record{}
	if rows.Next() {
		var ciphertext sql.RawBytes
		if err := rows.Scan(&r.ID, &r.Host, &r.URL, &r.Username, &ciphertext, &r.CreatedAt, &r.UpdatedAt); err != nil {
			panic(err)
		}
		r.Password = symcrypto.Decrypt(ciphertext, d.passphrase)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return r
}

func (d Database) Insert(recs []Record, skipFirst bool) int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	counter := 0
	for i, rec := range recs {
		if i == 0 && skipFirst {
			continue
		}
		ciphertext := symcrypto.Encrypt(rec.Password, d.passphrase)
		_, execErr := tx.ExecContext(ctx, insert_record, rec.Host, rec.URL, rec.Username, ciphertext, time.Now().Local(), time.Now().Local())
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				panic(fmt.Sprintf("Insert failed: %v, unable to rollback: %v\n", execErr, rollbackErr))
			}
			panic(fmt.Sprintf("Insert failed: %v", execErr))
		}
		counter = counter + 1
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
	return counter
}

func (d Database) Update(id uint, rec Record) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, execErr := tx.ExecContext(ctx, update_record, symcrypto.Encrypt(rec.Password, d.passphrase), time.Now().Local(), rec.ID)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			panic(fmt.Sprintf("Update failed: %v, unable to rollback: %v\n", execErr, rollbackErr))
		}
		panic(fmt.Sprintf("Update failed: %v", execErr))
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func (d Database) TableExists() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt, err := d.db.PrepareContext(ctx, check_table)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			panic(err)
		}
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return count == 1
}

func (d Database) SetupTable() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	_, execErr := tx.ExecContext(ctx, create_table)
	if execErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			panic(fmt.Sprintf("Setup failed: %v, unable to rollback: %v\n", execErr, rollbackErr))
		}
		panic(fmt.Sprintf("Setup failed: %v", execErr))
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

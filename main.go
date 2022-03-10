package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"passmgr/config"
	"passmgr/database"
	"passmgr/input"
	"passmgr/ops"
)

var (
	config_file string
	help        bool
	import_csv  string
	insert      bool
	search      string
	update      uint
	remove      uint
	version     bool
)

const VERSION string = "v0.0.1"

func main() {
	processFlags()

	conf, err := config.ReadConfig(config_file)
	if err != nil {
		log.Fatal(err)
	}

	var passphrase string
	for {
		passphrase, err = input.ReadSecret("Enter passphrase (it won't echo)")
		if err == nil {
			break
		}
	}

	db := prepareDB(conf, passphrase)
	defer enforceDBPermissions(conf)

	if search != "example.com" {
		ops.Search(search, db)
	} else if update != 0 {
		ops.Update(update, db)
	} else if remove != 0 {
		ops.Delete(remove, db)
	} else if insert {
		ops.Insert(db)
	} else if import_csv != "" {
		ops.ImportCSV(import_csv, true, db)
	}
}

func processFlags() {
	flag.StringVar(&config_file, "config", config.DefaultConfigPath(), "--config </path/to/config> (defaults to $HOME/.passmgr_config)")
	flag.StringVar(&import_csv, "import-csv", "", "--import-csv </path/to/file.csv>")
	flag.BoolVar(&insert, "insert", false, "--insert")
	flag.StringVar(&search, "search", "example.com", "--search <host>")
	flag.UintVar(&update, "update", 0, "--update <uint>")
	flag.UintVar(&remove, "delete", 0, "--delete <uint>")
	flag.BoolVar(&help, "help", false, "--help")
	flag.BoolVar(&version, "version", false, "--version")

	flag.Parse()

	if flag.NFlag() == 0 {
		usage(1)
	}

	if version {
		fmt.Println(os.Args[0] + " " + VERSION)
		os.Exit(0)
	}
	if help {
		usage(0)
	}
}

func enforceDBPermissions(conf config.Config) {
	if conf.EnforceDBPermissions {
		if _, err := os.Lstat(conf.Dbfile); err != nil {
			// dbfile doesn't exit yet
		} else if err := os.Chmod(conf.Dbfile, 0600); err != nil {
			panic(err)
		}
	}
}

func prepareDB(conf config.Config, passphrase string) database.Database {
	enforceDBPermissions(conf)

	db := database.NewDatabase(conf.Dbfile, passphrase)
	if !db.TableExists() {
		db.SetupTable()
	}

	return db
}

func usage(code int) {
	fmt.Println("passmgr --search <host> | --insert | --import-csv <path/to/csv> | --update <uint> | --config </path/to/config> | --help")
	os.Exit(code)
}

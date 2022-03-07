package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Dbfile               string `json:"dbfile"`
	EnforceDBPermissions bool   `json:"enforce-db-permissions"`
}

func DefaultConfigPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	filepath := dirname + string(os.PathSeparator) + ".passmgr_config"

	return filepath
}

func ReadConfig(filepath string) (Config, error) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	return config, nil
}

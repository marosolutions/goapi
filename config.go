package goapi

import (
	"encoding/json"
	"log"
	"os"
)

// Config ...
type Config struct {
	Service struct {
		Host   string `json:"host"`
		Port   string `json:"port"`
		Path   string `json:"path"`
		Method string `json:"method"`
	} `json:"service"`
	Cors struct {
		AllowOrigins string `json:"allow_origins"`
		AllowHeaders string `json:"allow_headers"`
	} `json:"cors"`
	Sort struct {
		DefaultSortOrder string   `json:"default_sort_order"`
		DefaultSortBy    string   `json:"default_sort_by"`
		ValidSortBy      []string `json:"valid_sort_by"`
	} `json:"sorting"`
	Pagination struct {
		DefaultPerPage string `json:"default_per_page"`
		MaxPerPage     string `json:"max_per_page"`
	} `json:"pagination"`
	Spanner struct {
		Project  string `json:"project"`
		Instance string `json:"instance"`
		Database string `json:"database"`
	} `json:"spanner"`
}

// ReadConfig ...
func ReadConfig() (config *Config) {
	env := os.Getenv("APP_ENV")
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	if env == "" {
		log.Fatal("Missing Envionment!")
	}
	filePath := "./config/" + env + ".json"

	configFile, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer configFile.Close()

	json.NewDecoder(configFile).Decode(&config)
	return
}

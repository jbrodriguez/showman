package main

import (
	"path/filepath"

	"github.com/namsral/flag"
)

// Config -
type Config struct {
	Version string `json:"version"`
}

// Settings -
type Settings struct {
	Config

	DataDir     string
	SourceDir   string
	DestDir     string
	TransferDir string
	ApiKey      string
	UserKey     string
	UserName    string
	LogToFile   bool

	Location string
}

// NewSettings -
func NewSettings(name, version, home string, locations []string) (*Settings, error) {
	var config, dataDir, sourceDir, destDir, transferDir, apiKey, userKey, userName string
	var logToFile bool

	location := SearchFile(name, locations)

	flag.StringVar(&config, "config", "", "config location")
	flag.StringVar(&dataDir, "datadir", filepath.Join(home, ".config", "showman"), "folder containing the user data")
	flag.StringVar(&sourceDir, "sourcedir", "", "folder containing the source media content")
	flag.StringVar(&destDir, "destdir", "", "where to move processed content")
	flag.StringVar(&transferDir, "transferdir", "", "where to move unprocessed content")
	flag.StringVar(&apiKey, "TVDB_APIKEY", "", "tvdb api key")
	flag.StringVar(&userKey, "TVDB_USERKEY", "", "tvdb user key")
	flag.StringVar(&userName, "TVDB_USERNAME", "", "tvdb user name")
	flag.BoolVar(&logToFile, "enableLogs", false, "true: logs to stdout and file; false: logs to stdout only")

	if found, _ := Exists(location); found {
		flag.Set("config", location)
	}
	flag.Parse()

	s := &Settings{}
	s.Version = version
	s.DataDir = dataDir
	s.SourceDir = sourceDir
	s.DestDir = destDir
	s.TransferDir = transferDir
	s.ApiKey = apiKey
	s.UserKey = userKey
	s.UserName = userName
	s.LogToFile = logToFile
	s.Location = location

	return s, nil
}

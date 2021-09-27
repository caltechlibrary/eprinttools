package eprinttools

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// Config holds a common configuration file used by eprinttools.
// Configuration file is expected to be in JSON format.
type Config struct {
	// Hostname for running service
	Hostname string `json:"hostname"`
	// EPrintURL points at EPrints 3.3.x Host
	EPrintURL string `json:"eprints_url"`
	// DbHost is the hostname and port running MySQL server for EPrints
	DbHost string `json:"db_host"`
	// DbUser is the username accessing MySQL server for EPrints
	DbUser string `json:"db_user"`
	// DbPassword is the password accessing MySQL server for EPrints
	DbPassword string `json:"db_password"`
	// Repositories is an array of string of DB names for EPrints repositories
	Repositories []string `json:"repositories"`
}

// NewConfig creates a new Config object setting values to defaults.
func NewConfig() *Config {
	cfg := new(Config)
	cfg.Hostname = "localhost:8484"
	cfg.EPrintURL = ""
	return cfg
}

func (cfg *Config) Load(fname string) error {
	if src, err := ioutil.ReadFile(fname); err != nil {
		return err
	} else {
		obj := map[string]string{}
		if err = json.Unmarshal(src, &obj); err != nil {
			return err
		}
		if value, ok := obj["hostname"]; ok == true {
			cfg.Hostname = fmt.Sprintf("%s", value)
		}
		if value, ok := obj["eprints_url"]; ok == true {
			cfg.EPrintURL = value
		}
	}
	return nil
}

func DisplayLicense(out io.Writer, appName string, license string) {
	fmt.Fprintf(out, strings.ReplaceAll(strings.ReplaceAll(license, "{appName}", appName), "{version}", Version))
}

func DisplayVersion(out io.Writer, appName string) {
	fmt.Fprintf(out, "\n%s %s\n", appName, Version)
}

func DisplayUsage(out io.Writer, appName string, flagSet *flag.FlagSet, description string, examples string, license string) {
	// Convert {appName} and {version} in description
	if description != "" {
		fmt.Fprintf(out, strings.ReplaceAll(description, "{appName}", appName))
	}
	flagSet.SetOutput(out)
	flagSet.PrintDefaults()

	if examples != "" {
		fmt.Fprintf(out, strings.ReplaceAll(examples, "{appName}", appName))
	}
	if license != "" {
		DisplayLicense(out, appName, license)
	}
	DisplayVersion(out, appName)
}

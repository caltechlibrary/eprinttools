package eprinttools

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	testDir := "testout"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		os.MkdirAll(testDir, 0775)
	}
	testSettings := path.Join(testDir, "settings-config-default.json")
	src := DefaultConfig()
	if err := os.WriteFile(testSettings, src, 0600); err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err := LoadConfig(testSettings)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestLoadConfig(t *testing.T) {
	dbUser := assertGetenvIsSet(t, "TEST_DB_USER")
	dbPassword := assertGetenvIsSet(t, "TEST_DB_PASSWORD")
	testDir := "testout"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		os.MkdirAll(testDir, 0775)
	}
	testSettings := path.Join(testDir, "settings-config.json")
	if _, err := os.Stat(testSettings); err != nil {
		src := []byte(fmtTxt(`{
    "jsonstore": "$DB_USER:$DB_PASSWORD@/test_repositories",
    "repositories": {
        "test_authors": {
            "dsn": "$DB_USER:$DB_PASSWORD@/test_authors",
            "base_url": "https://authors.library.example.edu",
            "write": false,
            "default_collection": "TestAUTHORS",
            "default_official_url": "https://resolver.library.example.edu",
            "default_rights": "No commercial reproduction, distribution, display or performance rights in this work are provided.",
            "default_refereed": "TRUE",
            "default_status": "archive",
            "strip_tags": true
        },
        "test_thesis": {
            "dsn": "$DB_USER:$DB_PASSWORD@/test_thesis",
            "base_url": "https://thesis.library.example.edu",
            "write": false,
            "default_collection": "TestTHESIS",
            "default_official_url": "https://resolver.library.example.edu",
            "default_rights": "No commercial reproduction, distribution, display or performance rights in this work are provided.",
            "default_refereed": "TRUE",
            "default_status": "archive",
            "strip_tags": true
        }
    }
}`, dbUser, dbPassword, "", "", ""))
		if err := os.WriteFile(testSettings, src, 0664); err != nil {
			t.Errorf("Failed to generate %q, %s", testSettings, err)
			t.FailNow()
		}
	}
	cfg, err := LoadConfig(testSettings)
	if err != nil {
		t.Errorf("LoadConfig(%q) failed, %s", testSettings, err)
	}
	if cfg == nil {
		t.Errorf("Configuration is nil")
	}
	expected := fmt.Sprintf("%s:%s@/test_repositories", dbUser, dbPassword)
	if expected != cfg.JSONStore {
		t.Errorf("Expected %q, got %q", expected, cfg.JSONStore)
	}
}

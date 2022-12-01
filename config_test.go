package eprinttools

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	t.Errorf("TestDefaultConfig() not implemented")
}

func TestLoadConfig(t *testing.T) {
	dbUser := os.Getenv("TEST_DB_USER")
	if dbUser == "" {
		t.Errorf("Expected TEST_DB_USER to be defined to run tests")
		t.FailNow()
	}
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	if dbPassword == "" {
		t.Errorf("Expected TEST_DB_PASSWORD to be defined to run tests")
		t.FailNow()
	}
	testDir := "testout"
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		os.MkdirAll(testDir, 0775)
	}
	testSettings := path.Join(testDir, "test_settings.json")
	if _, err := os.Stat(testSettings); err != nil {
		src := []byte(strings.ReplaceAll(strings.ReplaceAll(`{
    "jsonstore": "{DB_USER}:{DB_PASSWORD}@/test_repositories",
	"aggregation": "{DB_USER}:{DB_PASSWORD}@/test_aggregations",
    "repositories": {
        "test_authors": {
            "dsn": "{DB_USER}:{DB_PASSWORD}@/test_authors",
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
            "dsn": "{DB_USER}:{DB_PASSWORD}@/test_thesis",
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
}`, "{DB_USER}", dbUser), "{DB_PASSWORD}", dbPassword))
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
	expected := fmt.Sprintf("%s:%q@/test_repositories", dbUser, dbPassword)
	if expected != cfg.JSONStore {
		t.Errorf("Expected %q, got %q", expected, cfg.JSONStore)
	}
	expected = fmt.Sprintf("%s:%s@/test_aggregations", dbUser, dbPassword)
	if expected != cfg.AggregationStore {
		t.Errorf("Expected %q, got %q", expected, cfg.JSONStore)
	}
}

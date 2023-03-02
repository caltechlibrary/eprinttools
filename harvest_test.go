package eprinttools

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

// Test01Harvester tests the harvest run method
func Test01Harvester(t *testing.T) {
	dbUser := assertGetenvIsSet(t, "TEST_DB_USER")
	dbPassword := assertGetenvIsSet(t, "TEST_DB_PASSWORD")
	baseURL := assertGetenvIsSet(t, "TEST_BASE_URL")
	collection := assertGetenvIsSet(t, "TEST_COLLECTION")
	cName := strings.ToLower(collection)

	dName := "testout"
	fName := path.Join(dName, "settings.json")
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		os.MkdirAll(dName, 0775)
	}
	src := []byte(fmtTxt(`{
    "jsonstore": "$DB_USER:$DB_PASSWORD@/test_collections",
    "eprint_repositories": {
        "$REPO_NAME": {
            "dsn": "$DB_USER:$DB_PASSWORD@/$REPO_NAME",
            "base_url": "$BASE_URL",
            "write": true,
            "default_collection": "$COLLECTION",
            "default_official_url": "https://resolver.library.example.edu",
            "default_rights": "No commercial reproduction, distribution, display or performance rights in this work are provided.",
            "default_refereed": "TRUE",
            "default_status": "inbox",
            "strip_tags": true,
			"public_only": true
    	}
	},
	"project_dir": "testout",
	"htdocs": "testout/htdocs"
}`, dbUser, dbPassword, baseURL, cName, collection))
	if err := os.WriteFile(fName, src, 0600); err != nil {
		t.Error(err)
		t.FailNow()
	}
	// Initialize DB Schema
	txt, err := HarvesterDBSchema(fName)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	sqlFName := path.Join(dName, "test_db.sql")
	if err := os.WriteFile(sqlFName, []byte(txt), 0664); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := initSqlDB(sqlFName); err != nil {
		t.Error(err)
		t.FailNow()
	}
	today := time.Now()
	startYear := today.Year() - 1

	start := fmt.Sprintf("%d-01-01", startYear)
	if err := RunHarvester(fName, start, "", "", false, true); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// Test02Datasets tests the Run datasets method
func Test02Datasets(t *testing.T) {
	dName := "testout"
	fName := path.Join(dName, "settings.json")
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		t.Errorf("could not find %q, %s", fName, err)
		t.FailNow()
	}
	if err := RunDatasets(fName, "", "", true); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// Test03Genfeed tests the genfeeds run method
func Test03Genfeeds(t *testing.T) {
	dName := "testout"
	fName := path.Join(dName, "settings.json")
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		t.Errorf("could not find %q, %s", fName, err)
		t.FailNow()
	}
	if err := RunGenfeeds(fName, true); err != nil {
		t.Error(err)
		t.FailNow()
	}
}


package eprinttools

// Tests genfeed.go funcs
import (
	"testing"
	"os"
	"path"
	"strings"
)

func TestGenerateDatasets(t *testing.T) {
	dbUser := assertGetenvIsSet(t, "TEST_DB_USER")
	dbPassword := assertGetenvIsSet(t, "TEST_DB_PASSWORD")
	baseURL := assertGetenvIsSet(t, "TEST_BASE_URL")
	collection := assertGetenvIsSet(t, "TEST_COLLECTION")
	cName := strings.ToLower(collection)

	dName := "testout"
	fName := path.Join(dName, "test_settings-genfeeds.json")
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		os.MkdirAll(dName, 0775)
		src := []byte(fmtTxt(`{
    "jsonstore": "$DB_USER:$DB_PASSWORD@/test_collections",
    "repositories": {
        "$REPO_NAME": {
            "dsn": "$DB_USER:$DB_PASSWORD@/$REPO_NAME",
            "base_url": "$BASE_URL",
            "write": true,
            "default_collection": "$COLLECTION",
            "default_official_url": "https://resolver.library.example.edu",
            "default_rights": "No commercial reproduction, distribution, display or performance rights in this work are provided.",
            "default_refereed": "TRUE",
            "default_status": "inbox",
            "strip_tags": true
    	},
	}
}`, dbUser, dbPassword, baseURL, cName, collection))
		if err := os.WriteFile(fName, src, 0600); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	t.Errorf("TestGenerateDatasets() not fully implemented")
}

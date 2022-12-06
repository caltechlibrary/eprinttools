package eprinttools

import (
	"testing"
	"os"
	"strings"
)

func fmtTxt(txt, dbUser string, dbPassword string, baseURL string, cName string, collection string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(txt, "$DB_USER", dbUser), "$DB_PASSWORD", dbPassword), "$BASE_URL", baseURL), "$REPO_NAME", cName), "$COLLECTION", collection)
}


func assertGetenvIsSet(t *testing.T, vName string) string {
	s := os.Getenv(vName)
	if s == "" {
		t.Errorf("$%s not set in environment, skipping test", vName)
		t.FailNow()
	}
	return s
}


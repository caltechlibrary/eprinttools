package eprinttools

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// initSqlDB takes a configuraiton and SQL schema file
// initializing the repository
func initSqlDB(sqlFName string) error {
	cmd := exec.Command("mysql")
	fName, err := filepath.Abs(sqlFName)
	if err != nil {
		return err
	}
	cmd.Stdin = strings.NewReader(fmt.Sprintf("source %s", fName))
	stdoutStderr, err := cmd.CombinedOutput()
	log.Printf("%s\n", stdoutStderr)
	if err != nil {
		return err
	}
	return nil
}

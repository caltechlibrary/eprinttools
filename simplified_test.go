package eprinttools

import (
	//"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

/* NOTE: It is expected harvested and sanitized test EPrints records
exist in testdata as test_eprint1.xml and test_eprint2.xml.
*/

func TestRecordFromEPrint(t *testing.T) {
	for i, name := range []string{"test_eprint1.xml", "test_eprint2.xml", "test_eprint-embargoed.xml"} {
		fName := path.Join("testdata", name)
		src, err := ioutil.ReadFile(fName)
		if err != nil {
			t.Errorf("Failed to read %q, %s", fName, err)
			t.FailNow()
		}
		eprints := new(EPrints)
		err = xml.Unmarshal(src, &eprints)
		if err != nil {
			t.Errorf("Failed to unmarshal %q, %s", fName, err)
			t.FailNow()
		}
		if len(eprints.EPrint) == 0 {
			t.Errorf("Expected at least 1 test record in %q", fName)
			t.FailNow()
		}
		for j, eprint := range eprints.EPrint {
			if rec, err := CrosswalkEPrintToRecord(eprint); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR in crosswalk:\n%s\n", rec.ToString())
				t.Errorf("CrosswalkEPrintToRecord() failed (%d:%d), %s", i, j, err)
			}
		}
	}
}

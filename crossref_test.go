package eprinttools

import (
	"io/ioutil"
	"testing"
)

func TestIssue35(t *testing.T) {
	doi := `10.1093/mnras/stab2505`
	testFile := `srctest/issue-35.json`
	src, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf(`Missing %q, %s`, testFile, err)
		t.FailNow()
	}
	obj := make(map[string]interface{})
	err = jsonDecode(src, &obj)
	if err != nil {
		t.Errorf(`Failed to unmarhsal %q, %s`, testFile, err)
		t.FailNow()
	}
	eprint, err := CrossRefWorksToEPrint(obj)
	if err != nil {
		t.Errorf(`Expected to crosswalk CrossRef for %q, %s`, doi, err)
		t.FailNow()
	}
	if eprint == nil {
		t.Errorf(`Should have a EPrint data structure populated.`)
		t.FailNow()
	}
	if eprint.Abstract == `` {
		t.Errorf(`Expected to have Abstract populated for %q`, doi)
	}
	if eprint.Subjects == nil || eprint.Subjects.Length() == 0 {
		t.Errorf(`Should have subjects populated for %q`, doi)
	}
	if eprint.Subjects.Length() != 2 {
		t.Errorf(`Expected 2 subjects for doi %q, got %d`, doi, eprint.Subjects.Length())
	}
}

func TestIssue36(t *testing.T) {
	doi := `10.1093/mnras/stab2505`
	testFile := `srctest/issue-35.json`
	src, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf(`Missing %q, %s`, testFile, err)
		t.FailNow()
	}
	obj := make(map[string]interface{})
	err = jsonDecode(src, &obj)
	if err != nil {
		t.Errorf(`Failed to unmarhsal %q, %s`, testFile, err)
		t.FailNow()
	}
	eprint, err := CrossRefWorksToEPrint(obj)
	if err != nil {
		t.Errorf(`Expected to crosswalk CrossRef for %q, %s`, doi, err)
		t.FailNow()
	}
	if eprint == nil {
		t.Errorf(`Should have a EPrint data structure populated.`)
		t.FailNow()
	}
	if eprint.Funders == nil || eprint.Funders.Length() == 0 {
		t.Errorf(`Missing funders for %q`, doi)
		t.FailNow()
	}
	// Expected two National Science foundation grant numbers
	cnt := 0 // fount the grant numbers found for NSF
	for _, item := range eprint.Funders.Items {
		if item.Agency == `National Science Foundation` && item.GrantNumber != `` {
			cnt++
		}
	}
	if cnt != 2 {
		t.Errorf(`Execpted 2 grant number entry for NSF, got %d`, cnt)
	}
}

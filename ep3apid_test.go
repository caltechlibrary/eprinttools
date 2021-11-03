package eprinttools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

func runServiceForTest(t *testing.T, fName string) {
	if err := RunExtendedAPI("Go test", fName); err != nil {
		t.Errorf(`RunExtendedAPI("Go test", %s) failed, %s`, fName, err)
	}
}

func runShutdownForTest(t *testing.T, appName string) {
	if Shutdown(appName, "simulated signal") != 0 {
		t.Errorf("Expected zero return for Shutdown()")
	}
}

func httpGet(u string) ([]byte, error) {
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func httpPost(u string, contentType string, data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	res, err := http.Post(u, contentType, buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	src, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return src, fmt.Errorf(`HTTP status code %d`, res.StatusCode)
	}
	return src, nil
}

func checkForHelpPages(hostname string, repoID string, route string) error {
	var u string
	switch {
	case repoID == `` && route == ``:
		u = fmt.Sprintf(`http://%s`, config.Hostname)
	case repoID == `` && route != ``:
		u = fmt.Sprintf(`http://%s/%s`, config.Hostname, route)
	case route == ``:
		u = fmt.Sprintf(`http://%s/%s`, config.Hostname, repoID)
	default:
		u = fmt.Sprintf(`http://%s/%s/%s`, config.Hostname, repoID, route)
	}
	src, err := httpGet(u)
	if err != nil {
		return fmt.Errorf("client error, %s", err)
	}
	if len(src) == 0 {
		return fmt.Errorf("expected src, got empty byte array for %s", u)
	}
	return nil
}

func runWriteTest(t *testing.T, repoID string, repo *DataSource, route string) {
	baseURL := fmt.Sprintf(`http://%s`, config.Hostname)
	repo, ok := config.Repositories[repoID]
	if ok == false || repo.Write == false {
		t.Errorf(`%s not configured for writing`, repoID)
		t.FailNow()
	}
	// Test writes to a Authors like database, e.g lemurAuthors
	testFile := path.Join(`srctest`, fmt.Sprintf(`%sEPrint.xml`, repoID))
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf(`Could not find %q, %s`, testFile, err)
		t.FailNow()
	}
	src, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf(`Cound not read %q, %s`, testFile, err)
		t.FailNow()
	}
	//t.Logf(`Read eprint XML %q: %s`, testFile, src)
	u := fmt.Sprintf(`%s/%s/eprint-import`, baseURL, repoID)
	src, err = httpPost(u, `application/xml`, src)
	if err != nil {
		t.Logf(`%s`, src)
		t.Errorf(`Post failed, %q, %s`, u, err)
		t.FailNow()
	}
	t.Logf(`Post returned: %s`, src)
	ids := []int{}
	if err := json.Unmarshal(src, &ids); err != nil {
		t.Errorf(`Failed to unmarshal post results fo %q, %s`, u, err)
		t.FailNow()
	}
	if len(ids) == 0 || ids[0] == 0 {
		t.Errorf(`Expected non zero id in ids list`)
		t.FailNow()
	}
}

func runReadTests(t *testing.T, repoID string, route string) {
	for repoID, dsn := range config.Repositories {
		baseURL := fmt.Sprintf(`http://%s`, config.Hostname)
		u := fmt.Sprintf(`%s/repository/%s`, baseURL, repoID)
		src, err := httpGet(u)
		if err != nil {
			t.Errorf(`Failed %s, %s`, u, err)
			t.FailNow()
		}
		repository := map[string][]string{}
		if err := json.Unmarshal(src, &repository); err != nil {
			t.Errorf(`Failed unmarshal %s, %s`, u, err)
			t.FailNow()
		}
		u = fmt.Sprintf(`%s/%s/keys`, baseURL, repoID)
		src, err = httpGet(u)
		if err != nil {
			t.Errorf(`Failed %s, %s`, u, err)
			t.FailNow()
		}
		keys := []int{}
		if err := json.Unmarshal(src, &keys); err != nil {
			t.Errorf(`Failed %s, %s`, u, err)
			t.FailNow()
		}
		if dsn.Write == false && len(keys) == 0 {
			t.Errorf(`Failed %s, expected greater than zero keys`, u)
			t.FailNow()
		}
	}
}

func checkRepoStructure(repoID string) error {
	u := fmt.Sprintf(`http://%s/repository/%s`, config.Hostname, repoID)
	src, err := httpGet(u)
	if err != nil {
		return err
	}
	if len(src) == 0 {
		return fmt.Errorf(`Expected JSON content for /repository/%s, got none`, repoID)
	}
	//FIXME: make sure we get back a valid JSON structure.
	return nil
}

func runClientForTest(t *testing.T, appName string, settings string) {
	// Run client tests
	const wait = 1
	t.Logf("Starting client test sequence in %d seconds", wait)
	time.Sleep(time.Second * wait)
	if err := Reload("Go test reload", "simulated signal", settings); err != nil {
		t.Errorf("Reload(), %s", err)
	}
	if len(config.Routes) == 0 {
		t.Errorf(`Expected some routes for test`)
	}
	// Check for repositories help
	if err := checkForHelpPages(config.Hostname, ``, ``); err != nil {
		t.Error(err)
	}
	if err := checkForHelpPages(config.Hostname, ``, `repositories`); err != nil {
		t.Error(err)
	}
	if err := checkForHelpPages(config.Hostname, ``, `repository`); err != nil {
		t.Error(err)
	}
	// Let's test for help pages
	for repoID, routes := range config.Routes {
		t.Logf(`Testing routes for %s`, repoID)
		if err := checkRepoStructure(repoID); err != nil {
			t.Errorf(`%s /repository/%s failed, %s`, repoID, repoID, err)
		}
		for route := range routes {
			if err := checkForHelpPages(config.Hostname, repoID, route); err != nil {
				t.Errorf(`%s, route %s, %s`, repoID, route, err)
			}
			switch route {
			case `eprint-import`:
				if repo, ok := config.Repositories[repoID]; ok && repo.Write {
					runWriteTest(t, repoID, repo, route)
				}
			default:
				runReadTests(t, repoID, route)
			}
		}
	}
}

func TestExtendedAPI(t *testing.T) {
	appName := `Extend API simulation`
	settings := `test-settings.json`
	if _, err := os.Stat(settings); os.IsNotExist(err) {
		t.Skipf(`Can't find %q, %s`, settings, err)
		t.SkipNow()
	}
	// Startup a test instance of the API.
	if err := InitExtendedAPI(settings); err != nil {
		t.Errorf(`Can't init extended API, %s`, err)
		t.FailNow()
	}
	// Run a test of the extended API to test.
	go func() {
		runServiceForTest(t, settings)
	}()
	defer func() {
		runShutdownForTest(t, appName)
	}()
	runClientForTest(t, appName, settings)
}

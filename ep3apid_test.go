package eprinttools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

func runServiceForTest(t *testing.T, api *EP3API, fName string) {
	if err := api.RunExtendedAPI("Go test", fName); err != nil {
		t.Errorf(`RunExtendedAPI("Go test", %s) failed, %s`, fName, err)
	}
}

func runShutdownForTest(t *testing.T, api *EP3API, appName string) {
	if api.Shutdown(appName, "simulated signal") != 0 {
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

func checkForHelpPages(api *EP3API, hostname string, repoID string, route string) error {
	var u string
	switch {
	case repoID == `` && route == ``:
		u = fmt.Sprintf(`http://%s`, api.Config.Hostname)
	case repoID == `` && route != ``:
		u = fmt.Sprintf(`http://%s/%s`, api.Config.Hostname, route)
	case route == ``:
		u = fmt.Sprintf(`http://%s/%s`, api.Config.Hostname, repoID)
	default:
		u = fmt.Sprintf(`http://%s/%s/%s`, api.Config.Hostname, repoID, route)
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

func runWriteTest(t *testing.T, api *EP3API, repoID string, repo *DataSource, route string) {
	baseURL := fmt.Sprintf(`http://%s`, api.Config.Hostname)
	repo, ok := api.Config.Repositories[repoID]
	if ok == false || repo.Write == false {
		t.Errorf(`%s not configured for writing`, repoID)
		t.FailNow()
	}
	// Test writes to a Authors like database, e.g lemurAuthors
	// NOTE: I used doi2eprintxml to create two new EPrint XML files for import testing.
	for i := 1; i <= 25; i++ {
		testFile := path.Join(`srctest`, fmt.Sprintf(`%s-import-api-%d.xml`, repoID, i))
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
		userID := 1
		u := fmt.Sprintf(`%s/%s/eprint-import/%d`, baseURL, repoID, userID)
		src, err = httpPost(u, `application/xml`, src)
		if err != nil {
			t.Logf(`%s`, src)
			t.Errorf(`Post failed, %q, %s`, u, err)
			t.FailNow()
		}
		t.Logf(`Post returned: %s`, src)
		ids := []int{}
		if err := jsonDecode(src, &ids); err != nil {
			t.Errorf(`Failed to unmarshal post results fo %q, %s`, u, err)
			t.FailNow()
		}
		if len(ids) == 0 || ids[0] == 0 {
			t.Errorf(`Expected non zero id in ids list`)
			t.FailNow()
		}
	}
}

func runReadTests(t *testing.T, api *EP3API, repoID string, route string) {
	for repoID, dsn := range api.Config.Repositories {
		baseURL := fmt.Sprintf(`http://%s`, api.Config.Hostname)
		u := fmt.Sprintf(`%s/repository/%s`, baseURL, repoID)
		src, err := httpGet(u)
		if err != nil {
			t.Errorf(`Failed %s, %s`, u, err)
			t.FailNow()
		}
		repository := map[string][]string{}
		if err := jsonDecode(src, &repository); err != nil {
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
		if err := jsonDecode(src, &keys); err != nil {
			t.Errorf(`Failed %s, %s`, u, err)
			t.FailNow()
		}
		if dsn.Write == false && len(keys) == 0 {
			t.Errorf(`Failed %s, expected greater than zero keys`, u)
			t.FailNow()
		}
	}
}

func checkRepoStructure(api *EP3API, repoID string) error {
	u := fmt.Sprintf(`http://%s/repository/%s`, api.Config.Hostname, repoID)
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

func runClientForTest(t *testing.T, api *EP3API, appName string, settings string) {
	// Run client tests
	const wait = 1
	t.Logf("Starting client test sequence in %d seconds", wait)
	time.Sleep(time.Second * wait)
	if err := api.Reload("Go test reload", "simulated signal", settings); err != nil {
		t.Errorf("Reload(), %s", err)
	}
	if len(api.Config.Routes) == 0 {
		t.Errorf(`Expected some routes for test`)
	}
	// Check for repositories help
	if err := checkForHelpPages(api, api.Config.Hostname, ``, ``); err != nil {
		t.Error(err)
	}
	if err := checkForHelpPages(api, api.Config.Hostname, ``, `repositories`); err != nil {
		t.Error(err)
	}
	if err := checkForHelpPages(api, api.Config.Hostname, ``, `repository`); err != nil {
		t.Error(err)
	}
	// Let's test for help pages
	for repoID, routes := range api.Config.Routes {
		t.Logf(`Testing routes for %s`, repoID)
		if err := checkRepoStructure(api, repoID); err != nil {
			t.Errorf(`%s /repository/%s failed, %s`, repoID, repoID, err)
		}
		for route := range routes {
			if err := checkForHelpPages(api, api.Config.Hostname, repoID, route); err != nil {
				t.Errorf(`%s, route %s, %s`, repoID, route, err)
			}
			switch route {
			case `eprint-import`:
				if repo, ok := api.Config.Repositories[repoID]; ok && repo.Write {
					runWriteTest(t, api, repoID, repo, route)
				}
			default:
				runReadTests(t, api, repoID, route)
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
	api := new(EP3API)
	// Startup a test instance of the API.
	if err := api.InitExtendedAPI(settings); err != nil {
		t.Errorf(`Can't init extended API, %s`, err)
		t.FailNow()
	}
	if api.Config == nil {
		t.Errorf("API not configured")
		t.FailNow()
	}
	if api.Config.Logfile == `` {
		t.Errorf(`expected logging to file for tests`)
		t.FailNow()
	}
	// Run a test of the extended API to test.
	go func() {
		runServiceForTest(t, api, settings)
	}()
	defer func() {
		runShutdownForTest(t, api, appName)
	}()
	runClientForTest(t, api, appName, settings)
}

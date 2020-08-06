package rc

import (
	"encoding/json"
	"fmt"
	"testing"
)

// Testing Rest Client access against ORCID REST API
func TestRestAPI(t *testing.T) {
	apiURL := "https://pub.sandbox.orcid.org"
	clientID := "APP-01XX65MXBF79VJGF"
	clientSecret := "3a87028d-c84c-4d5f-8ad5-38a93181c9e1"
	testORCID := "0000-0003-0900-6903"

	// Test setup
	api, err := New(apiURL, OAuth, clientID, clientSecret)
	if err != nil {
		t.Errorf("Can't create API, %s", err)
		t.FailNow()
	}
	if api == nil {
		t.Errorf("API shouldn't be nil")
		t.FailNow()
	}

	// Test internal login method
	err = api.oAuthLogin()
	if err != nil {
		t.Errorf("Can't authenticate, %s", err)
		t.FailNow()
	}
	api.token = ""

	// Test request method
	src, err := api.Request("get", "/v2.0/"+testORCID+"/record", map[string]string{})
	if err != nil {
		t.Errorf("request profile failed, %s", err)
		t.FailNow()
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(src, &data); err != nil {
		t.Errorf("Can't unmashall JSON response, %s", err)
		t.FailNow()
	}

	if _, ok := data["orcid-identifier"]; ok != true {
		t.Errorf("missing orcid-identifier")
		t.FailNow()
	}
}

// TestURLStringAuth
func TestURLStringAuth(t *testing.T) {
	clientID := "APP-01XX65MXBF79VJGF"
	clientSecret := "3a87028d-c84c-4d5f-8ad5-38a93181c9e1"
	apiURL := fmt.Sprintf("https://%s:%s@pub.sandbox.orcid.org", clientID, clientSecret)
	testORCID := "0000-0003-0900-6903"

	// Test setup
	api, err := New(apiURL, OAuth, "", "")
	if err != nil {
		t.Errorf("Can't create API, %s", err)
		t.FailNow()
	}
	if api == nil {
		t.Errorf("API shouldn't be nil")
		t.FailNow()
	}

	// Test internal login method
	err = api.oAuthLogin()
	if err != nil {
		t.Errorf("Can't authenticate, %s", err)
		t.FailNow()
	}
	api.token = ""

	// Test request method
	src, err := api.Request("get", "/v2.0/"+testORCID+"/record", map[string]string{})
	if err != nil {
		t.Errorf("request profile failed, %s", err)
		t.FailNow()
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(src, &data); err != nil {
		t.Errorf("Can't unmashall JSON response, %s", err)
		t.FailNow()
	}

	if _, ok := data["orcid-identifier"]; ok != true {
		t.Errorf("missing orcid-identifier")
		t.FailNow()
	}
}

package eprinttools

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// GetEPrint fetches an EPrint record via the EPrint REST API.
func GetEPrint(baseURL string, eprintID int) (*EPrints, error) {
	var (
		username string
		password string
		auth     string
		src      []byte
	)
	endPoint := fmt.Sprintf("%s/rest/eprint/%d.xml", baseURL, eprintID)
	u, err := url.Parse(endPoint)
	if err != nil {
		return nil, fmt.Errorf("%q, %s,", endPoint, err)
	}
	username, password, auth = "", "", "basic"
	if userinfo := u.User; userinfo != nil {
		username = userinfo.Username()
		if secret, isSet := userinfo.Password(); isSet {
			password = secret
		}
	}

	// NOTE: We build our client request object so we can
	// set authentication if necessary.
	req, err := http.NewRequest("GET", endPoint, nil)
	switch strings.ToLower(auth) {
	case "basic":
		req.SetBasicAuth(username, password)
	}
	req.Header.Set("User-Agent", Version)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		src, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%s for %s", res.Status, endPoint)
	}
	if len(bytes.TrimSpace(src)) == 0 {
		return nil, fmt.Errorf("No data")
	}

	data := new(EPrints)
	err = xml.Unmarshal(src, &data)
	if err != nil {
		return nil, err
	}
	for _, e := range data.EPrint {
		e.SyntheticFields()
	}
	return data, nil
}

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

// GetEPrint fetches a single EPrint record via the EPrint REST API.
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

// GetKeys returns a list of eprint record ids from the EPrints REST API
func GetKeys(baseURL string, authType int, username string, secret string) ([]string, error) {
	var (
		results []string
	)

	workURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if workURL.Path == "" {
		workURL.Path = path.Join("rest", "eprint") + "/"
	} else {
		p := workURL.Path
		workURL.Path = path.Join(p, "rest", "eprint") + "/"
	}
	// Switch to use Rest Client Wrapper
	rest, err := rc.New(workURL.String(), authType, username, secret)
	if err != nil {
		return nil, fmt.Errorf("requesting %s, %s", workURL.String(), err)
	}
	content, err := rest.Request("GET", workURL.Path, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("requested %s, %s", workURL.String(), err)
	}
	eIDs := new(ePrintIDs)
	err = xml.Unmarshal(content, &eIDs)
	if err != nil {
		return nil, err
	}
	// Build a list of Unique IDs in a map, then convert unique querys to results array
	m := make(map[string]bool)
	for _, val := range eIDs.IDs {
		if strings.HasSuffix(val, ".xml") == true {
			eprintID := strings.TrimSuffix(val, ".xml")
			if _, hasID := m[eprintID]; hasID == false {
				// Save the new ID found
				m[eprintID] = true
				// Only store Unique IDs in result
				results = append(results, eprintID)
			}
		}
	}
	return results, nil
}

// GetModifiedKeys returns a list of eprint record ids from the EPrints REST API that match the modification date range
func GetModifiedKeys(baseURL string, authType int, username string, secret string, start time.Time, end time.Time, verbose bool) ([]string, error) {
	var (
		results []string
	)
	// need to calculate the base restDocPath
	workURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if workURL.Path == "" {
		workURL.Path = path.Join("rest", "eprint") + "/"
	} else {
		p := workURL.Path
		workURL.Path = path.Join(p, "rest", "eprint") + "/"
	}
	restDocPath := workURL.Path
	// Switch to use Rest Client Wrapper
	rest, err := rc.New(baseURL, authType, username, secret)
	if err != nil {
		return nil, fmt.Errorf("requesting %s, %s", baseURL, err)
	}

	// Pass baseURL to GetKeys(), get key list then filter for modified times.
	pid := os.Getpid()
	keys, err := GetKeys(baseURL, authType, username, secret)
	// NOTE: consecutiveFailedCount tracks repeated failures
	// e.g. You need to authenticate with the server to get
	// modified information.
	consecutiveFailedCount := 0
	for _, key := range keys {
		// form a request to the REST API for just the modified date
		docPath := path.Join(restDocPath, key, "lastmod.txt")
		lastModified, err := rest.Request("GET", docPath, map[string]string{})
		if err != nil {
			if verbose == true {
				log.Printf("(pid: %d) request failed, %s", pid, err)
			}
			consecutiveFailedCount++
			if consecutiveFailedCount >= maxConsecutiveFailedRequests {
				return results, err
			}
		} else {
			consecutiveFailedCount = 0
			datestring := fmt.Sprintf("%s", lastModified)
			if len(datestring) > 9 {
				datestring = datestring[0:10]
			}
			// Parse the modified date and compare to our range
			if dt, err := time.Parse("2006-01-02", datestring); err == nil && dt.Unix() >= start.Unix() && dt.Unix() <= end.Unix() {
				// If range is OK then add the key to results
				results = append(results, key)
			}
		}
	}
	return results, nil
}

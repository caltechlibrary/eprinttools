//
// Package eprinttools is a collection of structures, functions and programs// for working with the EPrints XML and EPrints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package eprinttools

/**
 * ep3api.go defines an extended EPrints web server.
 */

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type EP3API struct {
	Config *Config
	Log    *log.Logger
}

//
// Handle parameters that may continue URL
func joinArgs(args []string) string {
	if len(args) > 0 && (strings.HasPrefix(args[0], `http:`) ||
		strings.HasPrefix(args[0], `https:`)) {
		args[0] = args[0] + `/`
	}
	return strings.Join(args, `/`)
}

//
// Date/Time expansion
//
func expandAproxDate(dt string, roundDown bool) string {
	switch len(dt) {
	case 4:
		if roundDown {
			dt += "-01-01"
		} else {
			dt += "-12-31"
		}
	case 7:
		if roundDown {
			dt += "-01 00:00:00"
		} else {
			//FIXME: need to handle 28/29/30/31 day mounths
			dt += "-31"
		}
	}
	return dt
}

//
// Package functions take the collected data and package into an HTTP
// response.
//

func (api *EP3API) packageIntIDs(w http.ResponseWriter, repoID string, values []int, err error) (int, error) {
	if err != nil {
		api.Log.Printf("ERROR: (%s) query error, %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	src, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		api.Log.Printf("ERROR: marshal error (%q), %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func (api *EP3API) packageStringIDs(w http.ResponseWriter, repoID string, values []string, err error) (int, error) {
	if err != nil {
		api.Log.Printf("ERROR: (%s) query error, %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	src, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		api.Log.Printf("ERROR: marshal error (%q), %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func (api *EP3API) packageObject(w http.ResponseWriter, repoID string, obj interface{}, err error) (int, error) {
	if err != nil {
		api.Log.Printf("ERROR: (%s) query error, %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	src, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		api.Log.Printf("ERROR: marshal error (%q), %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func (api *EP3API) packageDocument(w http.ResponseWriter, src string) (int, error) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, src)
	return 200, nil
}

func (api *EP3API) packageJSON(w http.ResponseWriter, repoID string, src []byte, err error) (int, error) {
	if err != nil {
		api.Log.Printf("ERROR: (%s) package JSON error, %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

func (api *EP3API) packageXML(w http.ResponseWriter, repoID string, src []byte, err error) (int, error) {
	if err != nil {
		api.Log.Printf("ERROR: (%s) package JSON error, %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintln(w, `<?xml version="1.0" encoding="utf-8"?>`)
	fmt.Fprintf(w, "%s", src)
	return 200, nil
}

// Given a POST Request, read the Request body and
// populate an EPrints structure or return an error.
func unpackageEPrintsPOST(r *http.Request) (*EPrints, error) {
	var (
		eprints *EPrints
		err     error
		src     []byte
	)
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		if src, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, fmt.Errorf("failed to read, %s", err)
		}
		err = jsonDecode(src, &eprints)
	case "application/xml":
		if src, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, fmt.Errorf("failed to read, %s", err)
		}
		err = xml.Unmarshal(src, &eprints)
	default:
		return nil, fmt.Errorf("%s not supported", contentType)
	}
	return eprints, err
}

//
// EP3API End Points
//

func (api *EP3API) versionEndPoint(w http.ResponseWriter, r *http.Request) (int, error) {
	return api.packageDocument(w, fmt.Sprintf("EPrint 3.3 extended API, %s", Version))
}

//
// Expose the repositories available
//

func (api *EP3API) repositoriesEndPoint(w http.ResponseWriter, r *http.Request) (int, error) {
	if strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, repositoryDocument())
	}
	repositories := []string{}
	for repository, _ := range api.Config.Repositories {
		repositories = append(repositories, repository)
	}
	src, err := json.MarshalIndent(repositories, "", "  ")
	if err != nil {
		return 500, fmt.Errorf("internal server error, %s", err)
	}
	return api.packageDocument(w, string(src))
}

func (api *EP3API) repositoryEndPoint(w http.ResponseWriter, r *http.Request, repoID string) (int, error) {
	if repoID == "" || strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, repositoryDocument())
	}
	if api.Config.Connections == nil {
		api.Log.Printf("Database connections not configured.")
		return 500, fmt.Errorf("internal server error")
	}
	if _, ok := api.Config.Connections[repoID]; !ok {
		api.Log.Printf("Database connections not configured for %s", repoID)
		return 404, fmt.Errorf("not found")
	}
	data, err := GetTablesAndColumns(api.Config, repoID)
	if err != nil {
		api.Log.Printf("GetTablesAndColumn(%q), %s", repoID, err)
		return 500, fmt.Errorf("internal server error")
	}
	src, err := json.MarshalIndent(data, "", "    ")
	return api.packageJSON(w, repoID, src, err)
}

//
// End Point for user information
//
func (api *EP3API) usernamesEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, userDocument(repoID))
	}
	usernames, err := GetUsernames(api.Config, repoID)
	return api.packageStringIDs(w, repoID, usernames, err)
}

func (api *EP3API) lookupUserIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) != 1 || strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, userDocument(repoID))
	}
	ids, err := GetUserID(api.Config, repoID, args[0])
	return api.packageIntIDs(w, repoID, ids, err)
}

func (api *EP3API) userEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) != 1 || strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, userDocument(repoID))
	}
	userid, err := strconv.Atoi(args[0])
	if err == nil {
		user, err := GetUserBy(api.Config, repoID, `userid`, userid)
		if user.HideEMail {
			user.EMail = ``
		}
		return api.packageObject(w, repoID, user, err)
	}
	user, err := GetUserBy(api.Config, repoID, `username`, args[0])
	if user.HideEMail {
		user.EMail = ``
	}
	return api.packageObject(w, repoID, user, err)
}

//
// End Point handles (route as defined `/{REPO_ID}/{END-POINT}/{ARGS}`)
//

func (api *EP3API) keysEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) != 0 || strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, keysDocument(repoID))
	}
	eprintIDs, err := GetAllEPrintIDs(api.Config, repoID)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) createdEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, createdDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("bad request")
	}
	var err error
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		start, err = time.Parse(timestamp, args[0])
		if err != nil {
			return 400, fmt.Errorf("bad request, (start) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("bad request, (end) %s", err)
		}
	}
	eprintIDs, err := GetEPrintIDsInTimestampRange(api.Config, repoID, "datestamp", start.Format(timestamp), end.Format(timestamp))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) updatedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, updatedDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("bad request")
	}
	var err error
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		start, err = time.Parse(timestamp, args[0])
		if err != nil {
			return 400, fmt.Errorf("bad request, (start) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("bad request, (end) %s", err)
		}
	}
	eprintIDs, err := GetEPrintIDsInTimestampRange(api.Config, repoID, "lastmod", start.Format(timestamp), end.Format(timestamp))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) deletedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, deletedDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("bad request")
	}
	var err error
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		start, err = time.Parse(timestamp, args[0])
		if err != nil {
			return 400, fmt.Errorf("bad request, (start) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		end, err = time.Parse(timestamp, args[1])
		if err != nil {
			return 400, fmt.Errorf("bad request, (end) %s", err)
		}
	}
	eprintIDs, err := GetEPrintIDsWithStatus(api.Config, repoID, "deletion", start.Format(timestamp), end.Format(timestamp))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) pubdateEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, pubdateDocument(repoID))
	}
	if (len(args) < 1) || (len(args) > 2) {
		return 400, fmt.Errorf("bad request")
	}
	var (
		err error
		dt  string
	)
	end := time.Now()
	start := time.Now()
	if (len(args) > 0) && (args[0] != "now") {
		dt = expandAproxDate(args[0], true)
		start, err = time.Parse(datestamp, dt)
		if err != nil {
			return 400, fmt.Errorf("bad request, (start date) %s", err)
		}
	}
	if (len(args) > 1) && (args[1] != "now") {
		dt = expandAproxDate(args[1], false)
		end, err = time.Parse(datestamp, dt)
		if err != nil {
			return 400, fmt.Errorf("bad request, (end date) %s", err)
		}
	}
	eprintIDs, err := GetEPrintIDsForDateType(api.Config, repoID, "published", start.Format(datestamp), end.Format(datestamp))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

//
// Person or Organizational end points
//

func (api *EP3API) creatorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, creatorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonOrOrgIDs(api.Config, repoID, "creators")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForPersonOrOrgID(api.Config, repoID, "creators", args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) creatorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, creatorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonNames(api.Config, repoID, "creators_name")
		return api.packageStringIDs(w, repoID, values, err)
	}
	family, given := args[0], ``
	if len(args) == 2 {
		given = args[1]
	}
	if len(args) > 2 {
		return 400, fmt.Errorf("bad request")
	}
	eprintIDs, err := GetEPrintIDsForPersonName(api.Config, repoID, "creators_name", family, given)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) creatorORCIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, creatorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllORCIDs(api.Config, repoID)
		return api.packageStringIDs(w, repoID, values, err)
	}
	//FIXME: Should validate ORCID format ...
	eprintIDs, err := GetEPrintIDsForORCID(api.Config, repoID, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) editorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, editorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonOrOrgIDs(api.Config, repoID, "editors")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForPersonOrOrgID(api.Config, repoID, "editors", args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) editorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, editorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonNames(api.Config, repoID, "editors_name")
		return api.packageStringIDs(w, repoID, values, err)
	}
	family, given := args[0], ``
	if len(args) == 2 {
		given = args[1]
	}
	if len(args) > 2 {
		return 400, fmt.Errorf("bad request")
	}
	eprintIDs, err := GetEPrintIDsForPersonName(api.Config, repoID, "editors_name", family, given)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) contributorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, contributorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonOrOrgIDs(api.Config, repoID, "contributors")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForPersonOrOrgID(api.Config, repoID, "contributors", args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) contributorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, contributorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonNames(api.Config, repoID, "contributors_name")
		return api.packageStringIDs(w, repoID, values, err)
	}
	family, given := args[0], ``
	if len(args) == 2 {
		given = args[1]
	}
	if len(args) > 2 {
		return 400, fmt.Errorf("bad request")
	}
	eprintIDs, err := GetEPrintIDsForPersonName(api.Config, repoID, "contributors_name", family, given)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) advisorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, advisorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonOrOrgIDs(api.Config, repoID, "thesis_advisor")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForPersonOrOrgID(api.Config, repoID, "thesis_advisor", args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) advisorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, advisorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonNames(api.Config, repoID, "thesis_advisor_name")
		return api.packageStringIDs(w, repoID, values, err)
	}
	family, given := args[0], ``
	if len(args) == 2 {
		given = args[1]
	}
	if len(args) > 2 {
		return 400, fmt.Errorf("bad request")
	}
	eprintIDs, err := GetEPrintIDsForPersonName(api.Config, repoID, "thesis_advisor_name", family, given)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) committeeIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, committeeDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonOrOrgIDs(api.Config, repoID, "thesis_committee")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForPersonOrOrgID(api.Config, repoID, "thesis_committee", args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) committeeNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, committeeDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllPersonNames(api.Config, repoID, "thesis_committee_name")
		return api.packageStringIDs(w, repoID, values, err)
	}
	family, given := args[0], ``
	if len(args) == 2 {
		given = args[1]
	}
	if len(args) > 2 {
		return 400, fmt.Errorf("bad request")
	}
	eprintIDs, err := GetEPrintIDsForPersonName(api.Config, repoID, "thesis_committee_name", family, given)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) corpCreatorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, corpCreatorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllItems(api.Config, repoID, "corp_creators_id")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForItem(api.Config, repoID, "corp_creators_id", args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) corpCreatorNameEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, corpCreatorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllItems(api.Config, repoID, "corp_creators_name")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForItem(api.Config, repoID, "corp_creators_name", joinArgs(args))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) corpCreatorURIEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, corpCreatorDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllItems(api.Config, repoID, "corp_creators_uri")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForItem(api.Config, repoID, "corp_creators_uri", joinArgs(args))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) groupIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, groupDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllItems(api.Config, repoID, "local_group")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForItem(api.Config, repoID, "local_group", args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) funderIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, funderDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllItems(api.Config, repoID, "funders_agency")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForItem(api.Config, repoID, "funders_agency", joinArgs(args))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) grantNumberEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, funderDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllItems(api.Config, repoID, "funders_grant_number")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForItem(api.Config, repoID, "funders_grant_number", joinArgs(args))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api EP3API) patentAssigneeEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, patentAssigneeDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllItems(api.Config, repoID, "patent_assignee")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForItem(api.Config, repoID, "patent_assignee", joinArgs(args))
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api EP3API) yearEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, yearDocument(repoID))
	}
	if len(args) == 0 {
		years, err := GetAllYears(api.Config, repoID)
		return api.packageIntIDs(w, repoID, years, err)
	}
	year, err := strconv.Atoi(args[0])
	if err != nil {
		return api.packageIntIDs(w, repoID, []int{}, err)
	}
	eprintIDs, err := GetEPrintIDsForYear(api.Config, repoID, year)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

//
// Unique identifiers (e.g. doi, issn, isbn) end points
//
func (api *EP3API) doiEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, doiDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "doi")
		return api.packageStringIDs(w, repoID, values, err)
	}
	doi := joinArgs(args)
	//FIXME: Should validate DOI format ...
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, "doi", doi)
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) pubmedIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, pubmedIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "pmid")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, `pmid`, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) pubmedCentralIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, pubmedCentralIDDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "pmc_id")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, `pmc_id`, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) issnEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, issnDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "issn")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, `issn`, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) isbnEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, isbnDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "isbn")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, `isbn`, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) patentApplicantEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, patentApplicantDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "patent_applicant")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, `patent_applicant`, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) patentNumberEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, patentNumberDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "patent_number")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, `patent_number`, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

func (api *EP3API) patentClassificationEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if strings.HasSuffix(r.URL.Path, `/help`) {
		return api.packageDocument(w, patentClassificationDocument(repoID))
	}
	if len(args) == 0 {
		values, err := GetAllUniqueID(api.Config, repoID, "patent_classification")
		return api.packageStringIDs(w, repoID, values, err)
	}
	eprintIDs, err := GetEPrintIDsForUniqueID(api.Config, repoID, `patent_classification`, args[0])
	return api.packageIntIDs(w, repoID, eprintIDs, err)
}

//
// Record End Point is experimental and may not make it to the
// release version of eprinttools. It accepts a EPrint ID and returns
// a simplfied JSON object.
//
func (api *EP3API) recordEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, recordDocument(repoID))
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("bad request")
	}
	ds, ok := api.Config.Repositories[repoID]
	if !ok {
		api.Log.Printf("Data Source not found for %q", repoID)
		return 404, fmt.Errorf("not found")
	}
	eprintID, err := strconv.Atoi(args[0])
	if err != nil {
		return 400, fmt.Errorf("bad request, eprint id invalid, %s", err)
	}

	eprint, err := SQLReadEPrint(api.Config, repoID, ds.BaseURL, eprintID)
	if err != nil {
		api.Log.Printf("SQLReadEPrint Error: %s\n", err)
		return 404, fmt.Errorf("not found")
	}
	//FIXME: this should just be a simple JSON from SQL ...
	simple, err := CrosswalkEPrintToRecord(eprint)
	if err != nil {
		return 500, fmt.Errorf("internal server error")
	}
	src, err := json.MarshalIndent(simple, "", "    ")
	return api.packageJSON(w, repoID, src, err)
}

//
// EPrint XML End Point is an experimental read end point provided
// in the extended EPrint API.  It reads EPrint data structures
// based on SQL calls to the MySQL database for a given EPrints
// repository. It accepts two content types - "application/xml"
// which returns EPrints XML or "application/json" which returns
// a JSON version of the EPrints XML.
//
func (api *EP3API) eprintEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if len(args) == 0 || repoID == "" || strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, eprintReadWriteDocument(repoID))
	}
	contentType := r.Header.Get("Content-Type")
	if r.Method != "GET" {
		return 405, fmt.Errorf("method not allowed")
	}
	if len(args) != 1 {
		return 400, fmt.Errorf("bad request, missing eprint id in path")
	}
	eprintID, err := strconv.Atoi(args[0])
	if err != nil {
		return 400, fmt.Errorf("bad request, eprint id (%s) %q not valid", repoID, eprintID)
	}
	ds, ok := api.Config.Repositories[repoID]
	if !ok {
		api.Log.Printf("Data Source not found for %q", repoID)
		return 404, fmt.Errorf("not found")
	}
	eprint, err := SQLReadEPrint(api.Config, repoID, ds.BaseURL, eprintID)
	if err != nil {
		return 404, fmt.Errorf("%s", err)

	}
	switch contentType {
	case "application/json":
		src, err := json.MarshalIndent(eprint, "", "    ")
		return api.packageJSON(w, repoID, src, err)
	case "application/xml":
		eprints := NewEPrints()
		eprints.XMLNS = `http://eprints.org/ep2/data/2.0`
		eprints.Append(eprint)
		src, err := xml.MarshalIndent(eprints, "", "  ")
		return api.packageXML(w, repoID, src, err)
	case "":
		eprints := NewEPrints()
		eprints.XMLNS = `http://eprints.org/ep2/data/2.0`
		eprints.Append(eprint)
		src, err := xml.MarshalIndent(eprints, "", "  ")
		return api.packageXML(w, repoID, src, err)
	default:
		return 415, fmt.Errorf("unsupported media type, %q", contentType)
	}
}

//
// EPrint Import End Point is an experimental write end point provided
// in the extended EPrint API.  It Accepts EPrints XML or the JSON
// expression of the EPrint XML and creates new EPrints medata records
//
// This end requires a POST method.  It accepts POST content encoded
// as either "application/json" or "application/xml".
//
// NOTE: this end point has to be enabled in the settings.json file
// defining the individual repository support. "write" needs to be
// set to true.
func (api *EP3API) eprintImportEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if repoID == "" {
		repoID = `{REPO_ID}`
		return api.packageDocument(w, eprintReadWriteDocument(repoID))
	}
	if r.Method == "GET" || strings.HasSuffix(r.URL.Path, "/help") {
		return api.packageDocument(w, eprintReadWriteDocument(repoID))
	}
	writeAccess := false
	dataSource, ok := api.Config.Repositories[repoID]
	if ok == true {
		writeAccess = dataSource.Write
	} else {
		api.Log.Printf("Data Source not found for %q", repoID)
		return 404, fmt.Errorf("not found")
	}
	if r.Method != "POST" || writeAccess == false {
		return 405, fmt.Errorf("method not allowed %q", r.Method)
	}
	// Check to see if we have application/xml or application/json
	// Get data from post
	eprints, err := unpackageEPrintsPOST(r)
	if err != nil {
		return 400, fmt.Errorf("bad request, POST failed (%s), %s", repoID, err)
	}
	for _, eprint := range eprints.EPrint {
		eprint.EPrintStatus = `buffer`
	}

	ids := []int{}
	ids, err = ImportEPrints(api.Config, repoID, dataSource, eprints)
	if err != nil {
		return 400, fmt.Errorf("bad request, create EPrint failed, %s", err)
	}
	return api.packageIntIDs(w, repoID, ids, err)
}

// The following define the API as a service handling errors,
// routes and logging.
//
func (api *EP3API) logRequest(r *http.Request, status int, err error) {
	q := r.URL.Query()
	errStr := "OK"
	if err != nil {
		errStr = fmt.Sprintf("%s", err)
	}
	if len(q) > 0 {
		api.Log.Printf("%s %s RemoteAddr: %s UserAgent: %s Query: %+v Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, errStr)
	} else {
		api.Log.Printf("%s %s RemoteAddr: %s UserAgent: %s Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, errStr)
	}
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, fmt.Sprintf(`%s`, err), statusCode)
	//fmt.Fprintf(w, `ERROR: %d %s`, statusCode, err)
}

func (api *EP3API) routeEndPoints(w http.ResponseWriter, r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	// Create args by removing empty strings from path parts
	repoID, endPoint, args := "", "", []string{}
	for _, arg := range parts {
		if arg != "" {
			// FIXME: URL decode that path part
			args = append(args, arg)
		}
	}
	if len(args) == 0 {
		return api.packageDocument(w, readmeDocument())
	}
	if len(args) == 1 {
		return api.packageDocument(w, strings.ReplaceAll(readmeDocument(), "{REPO_ID}", args[0]))
	}
	// Expected URL structure of `/<REPO_ID>/<END_POINT>/<ARGS>`
	if len(args) == 2 {
		repoID, endPoint, args = args[0], args[1], []string{}
	} else {
		repoID, endPoint, args = args[0], args[1], args[2:]
	}
	if routes, hasRepo := api.Config.Routes[repoID]; hasRepo == true {
		// Confirm we have a route
		if fn, hasRoute := routes[endPoint]; hasRoute == true {
			// Confirm we have a DB connection
			if _, hasConnection := api.Config.Connections[repoID]; hasConnection == true {
				return fn(w, r, repoID, args)
			}
		}
	}
	return 404, fmt.Errorf("Not Found")
}

func (api *EP3API) routeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		statusCode int
	)
	if r.Method != "GET" && r.Method != "POST" {
		statusCode, err = 405, fmt.Errorf("method not allowed, %q", r.Method)
		handleError(w, statusCode, err)
	} else {
		switch {
		case r.URL.Path == "/version":
			statusCode, err = api.versionEndPoint(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		case r.URL.Path == "/favicon.ico":
			statusCode, err = 200, nil
			fmt.Fprintf(w, "")
			//statusCode, err = 404, fmt.Errorf("Not Found")
			//handleError(w, statusCode, err)
		case strings.HasPrefix(r.URL.Path, "/repositories"):
			statusCode, err = api.repositoriesEndPoint(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		case strings.HasPrefix(r.URL.Path, "/repository"):
			repoID := strings.TrimPrefix(r.URL.Path, "/repository/")
			statusCode, err = api.repositoryEndPoint(w, r, repoID)
			if err != nil {
				handleError(w, statusCode, err)
			}
		default:
			statusCode, err = api.routeEndPoints(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		}
	}
	api.logRequest(r, statusCode, err)
}

// hasColumn takes a table map (key = table name, value is
// array of column names) and return true if table and column found.
func hasColumn(tableMap map[string][]string, tableName string, columnName string) bool {
	if columns, ok := tableMap[tableName]; ok {
		for _, value := range columns {
			if value == columnName {
				return true
			}
		}
	}
	return false
}

// Shutdown shutdowns the EPrints extended API web service started with
// RunExtendedAPI.
func (api *EP3API) Shutdown(appName string, sigName string) int {
	exitCode := 0
	pid := os.Getpid()
	api.Log.Printf(`Received signal %s`, sigName)
	api.Log.Printf(`Closing database connections %s pid: %d`, appName, pid)
	if err := CloseConnections(api.Config); err != nil {
		exitCode = 1
	}
	api.Log.Printf(`Shutdown completed %s pid: %d exit code: %d `, appName, pid, exitCode)
	return exitCode
}

// Reload performs a Shutdown and an init after re-reading
// in the settings.json file.
func (api *EP3API) Reload(appName string, sigName string, settings string) error {
	exitCode := api.Shutdown(appName, sigName)
	if exitCode != 0 {
		return fmt.Errorf("Reload failed, could not shutdown the current processes")
	}
	api.Log.Printf("Restarting %s using %s", appName, settings)
	return api.InitExtendedAPI(settings)
}

func (api *EP3API) InitExtendedAPI(settings string) error {
	var err error
	// NOTE: This reads the settings file and creates a global
	// config object.
	api.Config, err = LoadConfig(settings)
	if err != nil {
		return fmt.Errorf("Failed to load %q, %s", settings, err)
	}
	if api.Config == nil {
		return fmt.Errorf("Missing configuration")
	}
	/* Setup logging */
	if api.Config.Logfile == `` {
		api.Log = log.Default()
	} else {
		// Append or create a new log file
		lp, err := os.OpenFile(api.Config.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		api.Log = log.New(lp, ``, log.LstdFlags)
	}
	if api.Config.Hostname == "" {
		return fmt.Errorf("Hostings hostname for service")
	}
	if api.Config.Repositories == nil || len(api.Config.Repositories) < 1 {
		return fmt.Errorf(`Missing "repositories" configuration`)
	}
	if err := OpenConnections(api.Config); err != nil {
		return fmt.Errorf(`Failed to open database connections, %s`, err)
	}
	if api.Config.Routes == nil {
		api.Config.Routes = map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){}
	}
	// This is a map standard endpoints and end point handlers.
	// NOTE: Eventually this should evolving into a registration
	// style design pattern based on unique fields in repositories.
	routes := map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){
		// Standard fields
		"keys":              api.keysEndPoint,
		"created":           api.createdEndPoint,
		"updated":           api.updatedEndPoint,
		"deleted":           api.deletedEndPoint,
		"pubdate":           api.pubdateEndPoint,
		"doi":               api.doiEndPoint,
		"record":            api.recordEndPoint,
		"eprint":            api.eprintEndPoint,
		"eprint-import":     api.eprintImportEndPoint,
		"creator-id":        api.creatorIDEndPoint,
		"creator-orcid":     api.creatorORCIDEndPoint,
		"editor-id":         api.editorIDEndPoint,
		"contributor-id":    api.contributorIDEndPoint,
		"advisor-id":        api.advisorIDEndPoint,
		"committee-id":      api.committeeIDEndPoint,
		"corp-creator-id":   api.corpCreatorIDEndPoint,
		"group-id":          api.groupIDEndPoint,
		"funder-id":         api.funderIDEndPoint,
		"grant-number":      api.grantNumberEndPoint,
		"creator-name":      api.creatorNameEndPoint,
		"editor-name":       api.editorNameEndPoint,
		"contributor-name":  api.contributorNameEndPoint,
		"advisor-name":      api.advisorNameEndPoint,
		"committee-name":    api.committeeNameEndPoint,
		"corp-creator-name": api.corpCreatorNameEndPoint,
		"corp-creator-uri":  api.corpCreatorURIEndPoint,
		"issn":              api.issnEndPoint,
		"isbn":              api.isbnEndPoint,
		"year":              api.yearEndPoint,
		"usernames":         api.usernamesEndPoint,
		"lookup-userid":     api.lookupUserIDEndPoint,
		"user":              api.userEndPoint,
	}

	/* NOTE: We need a DB connection to MySQL for each
	   EPrints repository supported by the API
	   for access to MySQL */
	if err := OpenConnections(api.Config); err != nil {
		return err
	}

	for repoID, dataSource := range api.Config.Repositories {
		// Add routes (end points) for the target repository
		for route, fn := range routes {
			if api.Config.Routes[repoID] == nil {
				api.Config.Routes[repoID] = map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){}
			}
			api.Config.Routes[repoID][route] = fn
		}
		// NOTE: make sure each end point is supported by repository
		// e.g. CaltechTHESIS doens't have "patent_number",
		// "patent_classification", "parent_assignee", "pmc_id",
		// or "pmid".
		if hasColumn(dataSource.TableMap, "eprint", "pmc_id") {
			api.Config.Routes[repoID]["pmcid"] = api.pubmedCentralIDEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "pmid") {
			api.Config.Routes[repoID]["pmid"] = api.pubmedIDEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "patent_applicant") {
			api.Config.Routes[repoID]["patent-applicant"] = api.patentApplicantEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "patent_number") {
			api.Config.Routes[repoID]["patent-number"] = api.patentNumberEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint", "patent_classification") {
			api.Config.Routes[repoID]["patent-classification"] = api.patentClassificationEndPoint
		}
		if hasColumn(dataSource.TableMap, "eprint_patent_assignee", "patent_assignee") {
			api.Config.Routes[repoID]["patent-assignee"] = api.patentAssigneeEndPoint
		}
	}
	return nil
}

func (api *EP3API) RunExtendedAPI(appName string, settings string) error {
	/* Setup web server */
	api.Log.Printf(`
%s %s

Using configuration %s

Process id: %d

EPrints 3.3.x Extended API

Listening on http://%s

Press ctl-c to terminate.
`, appName, Version, settings, os.Getpid(), api.Config.Hostname)

	/* Listen for Ctr-C */
	processControl := make(chan os.Signal, 1)
	signal.Notify(processControl, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		sig := <-processControl
		switch sig {
		case syscall.SIGINT:
			os.Exit(api.Shutdown(appName, sig.String()))
		case syscall.SIGTERM:
			os.Exit(api.Shutdown(appName, sig.String()))
		case syscall.SIGHUP:
			if err := api.Reload(appName, sig.String(), settings); err != nil {
				api.Log.Println(err)
				os.Exit(1)
			}
		default:
			os.Exit(api.Shutdown(appName, sig.String()))
		}
	}()

	http.HandleFunc("/", api.routeHandler)
	return http.ListenAndServe(api.Config.Hostname, nil)
}

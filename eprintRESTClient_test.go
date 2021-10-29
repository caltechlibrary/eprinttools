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

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
)

func configRestClient() (*Config, error) {
	if _, err := os.Stat(`test-settings.json`); os.IsNotExist(err) {
		return nil, err
	}
	src, err := ioutil.ReadFile(`test-settings.json`)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = json.Unmarshal(src, &config)
	if err != nil {
		return nil, err
	}
	configuredClientFound := false
	for _, dsn := range config.Repositories {
		if dsn.RestAPI != `` {
			configuredClientFound = true
		}
	}
	if !configuredClientFound {
		return config, fmt.Errorf(`RestAPI not defined for any repositories`)
	}
	return config, err
}

func TestRestClient(t *testing.T) {
	config, err := configRestClient()
	if err != nil {
		t.Skipf(`Skipping TestRestClient(), %s`, err)
		t.SkipNow()
	}
	for repoID, dataSource := range config.Repositories {
		eprintURL := dataSource.RestAPI
		if eprintURL == `` {
			t.Logf(`RestAPI not configured for %s`, repoID)
		} else {
			keys, err := GetKeys(eprintURL)
			if err != nil {
				t.Errorf(`  GetKeys(%q) returned an error, %s`, eprintURL, err)
				t.FailNow()
			}
			if len(keys) < 1 {
				t.Errorf(`  Expected some keys form GetKeys(%q)`, eprintURL)
				t.FailNow()
			}
			// Take a sample of 25 keys
			sampleSize := 25
			if len(keys) > sampleSize {
				rand.Shuffle(len(keys), func(i, j int) {
					keys[i], keys[j] = keys[j], keys[i]
				})
				keys = keys[0:sampleSize]
				sort.Ints(keys)
			}
			t.Logf("Running test with %d randomly selected keys\n", len(keys))

			spinner := "._-+xX#*#Xx+-_."
			t.Logf("  Testing GetEPrint(baseURL, key) %q ", repoID)
			for i, key := range keys {
				t.Logf("\r%s", string(spinner[i%len(spinner)]))

				//NOTE: we need to check ep and raw if we don't and an error
				rec, err := GetEPrint(eprintURL, key)
				if err != nil {
					sErr := fmt.Sprintf("%s", err)
					// NOTE: We should get an error for 401's, or when
					// we try to retrieve something with eprint_status of
					// buffer, deletion, and inbox.
					if strings.HasPrefix(sErr, "401") == false &&
						strings.HasSuffix(sErr, "buffer") == false &&
						strings.HasSuffix(sErr, "deletion") == false &&
						strings.HasSuffix(sErr, "inbox") == false {
						t.Errorf("%d GetEPrint(%q, %d) -> %q", i, eprintURL, key, err)
						t.FailNow()
					}
				}
				if len(rec.EPrint) == 1 {
					if rec.EPrint[0].ID == "" {
						t.Errorf("Expected a eprint ID for %q Key: %d", repoID, key)
					}
					if rec.EPrint[0].EPrintStatus == "" {
						t.Errorf("Expected a EPrintStatus for %q Key: %d", repoID, key)
					}
					if rec.EPrint[0].EPrintStatus == "published" {
						if rec.EPrint[0].Title == "" {
							t.Errorf("Expected a Title for %q Key: %d", repoID, key)
						}
					}
				} else {
					t.Errorf("Expected to have one eprint from %q key: %d", repoID, key)
				}
			}
		}
		t.Log("\r \n")
	}
}

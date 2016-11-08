//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
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
package epgo

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// cfg is configuration to access the EPrints REST API for tests
var (
	cfg Config
)

func TestMain(m *testing.M) {
	cfg.MergeEnv("EPGO", "API_URL", "")
	cfg.MergeEnv("EPGO", "DBNAME", "")
	cfg.MergeEnv("EPGO", "BLEVE", "")
	cfg.MergeEnv("EPGO", "HTDOCS", "")
	cfg.MergeEnv("EPGO", "TEMPLATE_PATH", "")
	cfg.MergeEnv("EPGO", "SITE_URL", "")
	os.Exit(m.Run())
}

func TestMergeEnv(t *testing.T) {
	var tCfg Config

	for _, term := range []string{"API_URL", "DBNAME", "BLEVE", "HTDOCS", "TEMPLATE_PATH", "SITE_URL"} {
		if err := tCfg.MergeEnv("EPGO", term, "test_"+term); err != nil {
			if strings.Compare(tCfg.Get(term), "test_"+term) != 0 {
				t.Error(fmt.Sprintf("%s_%s error %s", "EPGO", term, err))
				t.FailNow()
			}
		}
	}
}

func TestHarvest(t *testing.T) {
	api, err := New(cfg)
	if err != nil {
		t.Errorf("Cannot create a new API instance %q", err)
		t.FailNow()
	}
	api.DBName = "test.boltdb"
	api.Htdocs = "testsite"
	_, err = os.Stat(api.Htdocs)
	if err != nil && os.IsNotExist(err) == true {
		err = os.Mkdir(api.Htdocs, 0770)
		if err != nil {
			t.Errorf("Cannot create %s %q", api.Htdocs, err)
			t.FailNow()
		}
	}

	err = api.ExportEPrints(5)
	if err != nil {
		t.Errorf("Cannot harvest for test site %q", err)
		t.FailNow()
	}

	err = api.BuildPages(5, "Recently Published", "recently-published", func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
		return api.GetPublishedRecords(start, count, direction)
	})
	if err != nil {
		t.Errorf("Cannot build test site %q", err)
		t.FailNow()
	}
	err = api.BuildPages(5, "Recent Articles", "recent-articles", func(api *EPrintsAPI, start, count, direction int) ([]*Record, error) {
		return api.GetPublishedArticles(start, count, direction)
	})
}

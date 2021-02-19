//
// Package eprinttools is a collection of structures and functions for working with the E-Prints REST API
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
	"os"
	"testing"
	"time"
)

func TestListEPrintsURI(t *testing.T) {
	eprintURL := os.Getenv("EPRINT_URL")
	if len(eprintURL) == 0 {
		t.Log("Skipping TestListEPrintsURI(), environment not set")
		t.SkipNow()
	}
	suppressNote := true

	api, err := New(eprintURL, suppressNote, "", "", "")
	if err != nil {
		t.Errorf("Failed to create new api, %s", err)
		t.FailNow()
	}

	uris, err := api.ListEPrintsURI()
	if err != nil {
		t.Errorf("listEPrintsURI() %s", err)
	}
	if len(uris) == 0 {
		t.Errorf("ListEPrintsURI() returned zero uri")
		t.FailNow()
	}

	start, _ := time.Parse("2006-01-02", "2018-01-01")
	end, _ := time.Parse("2006-01-02", "2018-12-31")
	uris, err = api.ListModifiedEPrintsURI(start, end, true)
	if err != nil {
		t.Errorf("ListModifiedEPrintsURI(%q, %q, true) %s", start.String(), end.String(), err)
	}
	if len(uris) == 0 {
		t.Errorf("Expected uris for list modified from %s to %s", start.String(), end.String())
	}
}

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
package eprinttools

import (
	//"encoding/json"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	// Caltech Library packages
	"github.com/caltechlibrary/simplified"
)

/* NOTE: It is expected harvested and sanitized test EPrints records
exist in testdata as test_eprint1.xml and test_eprint2.xml.
*/

func TestRecordFromEPrint(t *testing.T) {
	for i, name := range []string{"test_eprint1.xml", "test_eprint2.xml", "test_eprint4.xml"} {
		fName := path.Join("testdata", name)
		src, err := ioutil.ReadFile(fName)
		if err != nil {
			t.Errorf("Failed to read %q, %s", fName, err)
			t.FailNow()
		}
		eprints := NewEPrints()
		err = xml.Unmarshal(src, &eprints)
		if err != nil {
			t.Errorf("Failed to unmarshal %q, %s", fName, err)
			t.FailNow()
		}
		if len(eprints.EPrint) == 0 {
			t.Errorf("Expected at least 1 test record in %q", fName)
			t.FailNow()
		}
		for j, eprint := range eprints.EPrint {
			rec := new(simplified.Record)
			if err := CrosswalkEPrintToRecord(eprint, rec); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR in crosswalk:\n%s\n", rec.ToString())
				t.Errorf("CrosswalkEPrintToRecord() failed (%d:%d), %s", i, j, err)
			}
		}
	}
}

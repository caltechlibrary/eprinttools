//
// py/libeprinttools.go implements a C shared library for implementing a eprinttools module in Python3
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
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
package main

// #cgo pkg-config: python-3.6
// #define Py_LIMITED_API
// #include <Python.h>
import (
	"C"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/eprinttools"
	"github.com/caltechlibrary/rc"
)

var (
	verbose    = false
	errorValue error
)

//export is_verbose
func is_verbose() C.int {
	if verbose == true {
		return C.int(1)
	}
	return C.int(0)
}

//export verbose_on
func verbose_on() {
	verbose = true
}

//export verbose_off
func verbose_off() {
	verbose = false
}

func error_dispatch(err error, s string, values ...interface{}) {
	errorValue = err
	if verbose == true {
		log.Printf(s, values...)
	}
}

//export error_message
func error_message() *C.char {
	if errorValue != nil {
		s := fmt.Sprintf("%s", errorValue)
		errorValue = nil
		return C.CString(s)
	}
	return C.CString("")
}

//export version
func version() *C.char {
	return C.CString(eprinttools.Version)
}

func apiCfg(m map[string]string) (string, int, string, string) {
	uri, username, password := "", "", ""
	authType := rc.AuthNone
	if _, ok := m["url"]; ok == true {
		uri = m["url"]
	}
	if _, ok := m["auth_type"]; ok == true {
		switch m["auth_type"] {
		case "BasicAuth":
			authType = rc.BasicAuth
		case "OAuth":
			authType = rc.OAuth
		case "Shibboleth":
			authType = rc.Shibboleth
		default:
			authType = rc.AuthNone
		}
	}
	if _, ok := m["username"]; ok == true {
		username = m["username"]
	}
	if _, ok := m["password"]; ok == true {
		password = m["password"]
	}
	return uri, authType, username, password
}

func dsCfg(m map[string]string) string {
	if cName, ok := m["dataset"]; ok == true {
		return cName
	}
	return ""
}

//export get_keys
func get_keys(cfg *C.char) *C.char {
	errorValue = nil
	m := map[string]string{}
	src := []byte(C.GoString(cfg))
	err := json.Unmarshal(src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.CString("")
	}
	base_url, authType, username, password := apiCfg(m)
	if base_url == "" {
		error_dispatch(err, "Missing url for config %s", src)
		return C.CString("")
	}

	keys, err := eprinttools.GetKeys(base_url, authType, username, password)
	if err != nil {
		error_dispatch(err, "Can't GetKeys(), %s", err)
		return C.CString("")
	}
	src, err = json.Marshal(keys)
	if err != nil {
		error_dispatch(err, "Can't marshal keys, %s", err)
		return C.CString("")
	}
	return C.CString(fmt.Sprintf("%s", src))
}

//export get_modified_keys
func get_modified_keys(cfg *C.char, cStart *C.char, cEnd *C.char) *C.char {
	errorValue = nil
	m := map[string]string{}
	src := []byte(C.GoString(cfg))
	err := json.Unmarshal(src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.CString("")
	}
	base_url, authType, username, password := apiCfg(m)
	if base_url == "" {
		error_dispatch(err, "Missing url for config %s", src)
		return C.CString("")
	}
	sStart := C.GoString(cStart)
	sEnd := C.GoString(cEnd)
	start, err := time.Parse("2006-01-02", sStart)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}
	end, err := time.Parse("2006-01-02", sEnd)
	if err != nil {
		error_dispatch(err, "%s", err)
		return C.CString("")
	}

	keys, err := eprinttools.GetModifiedKeys(base_url, authType, username, password, start, end)
	if err != nil {
		error_dispatch(err, "Can't GetModifiedKeys(), %s", err)
		return C.CString("")
	}
	src, err = json.Marshal(keys)
	if err != nil {
		error_dispatch(err, "Can't marshal keys, %s", err)
		return C.CString("")
	}
	return C.CString(fmt.Sprintf("%s", src))
}

//export get_metadata
func get_metadata(cfg *C.char, cKey *C.char, cSave C.int) *C.char {
	save := false
	if int(cSave) == 1 {
		save = true
	}
	key := C.GoString(cKey)
	m := map[string]string{}
	src := []byte(C.GoString(cfg))
	err := json.Unmarshal(src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.CString("")
	}
	base_url, authType, username, password := apiCfg(m)
	eprint, xml_src, err := eprinttools.GetEPrints(base_url, authType, username, password, key)
	if err != nil {
		error_dispatch(err, "can't get eprint %s, %s", key, err)
		return C.CString("")
	}
	src, err = json.Marshal(eprint)
	if err != nil {
		error_dispatch(err, "can't marshal eprint %s, %s", key, err)
		return C.CString("")
	}

	if save {
		collectionName := dsCfg(m)
		c, err := dataset.Open(collectionName)
		if err != nil {
			error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
			return C.CString(fmt.Sprintf("%s", src))
		}
		c.Close()
		if c.HasKey(key) {
			if err := c.UpdateJSON(key, src); err != nil {
				error_dispatch(err, "can't save %s to %s, %s", key, collectionName, err)
				return C.CString(fmt.Sprintf("%s", src))
			}
		} else {
			if err := c.CreateJSON(key, src); err != nil {
				error_dispatch(err, "can't save %s to %s, %s", key, collectionName, err)
				return C.CString(fmt.Sprintf("%s", src))
			}
		}
		if err := c.AttachFile(key, key+".xml", bytes.NewReader(xml_src)); err != nil {
			error_dispatch(err, "can't attach %s.xml to %s in %s, %s", key, key, collectionName, err)
			return C.CString(fmt.Sprintf("%s", src))
		}
	}
	return C.CString(fmt.Sprintf("%s", src))
}

func main() {}

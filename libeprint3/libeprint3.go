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

import (
	"C"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/eprinttools"
	"github.com/caltechlibrary/rc"
)

var (
	verbose    = false
	errorValue error
	xmlBuffer  []byte
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
	keys, err := eprinttools.GetModifiedKeys(base_url, authType, username, password, start, end, verbose)
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
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
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
	src, err := json.Marshal(eprint)
	if err != nil {
		error_dispatch(err, "can't marshal eprint %s, %s", key, err)
		return C.CString("")
	}

	if save {
		xmlBuffer = xml_src
	}
	return C.CString(fmt.Sprintf("%s", src))
}

//export get_buffered_xml
func get_buffered_xml() *C.char {
	src := fmt.Sprintf("%s", xmlBuffer)
	xmlBuffer = nil
	return C.CString(src)
}

//export get_eprint_xml
func get_eprint_xml(cfg *C.char, cKey *C.char) C.int {
	key := C.GoString(cKey)
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.int(0)
	}
	base_url, authType, username, password := apiCfg(m)
	eprint, xml_src, err := eprinttools.GetEPrints(base_url, authType, username, password, key)
	if err != nil {
		error_dispatch(err, "can't get eprint %s, %s", key, err)
		return C.int(0)
	}
	src, err := json.Marshal(eprint)
	if err != nil {
		error_dispatch(err, "can't marshal eprint %s, %s", key, err)
		return C.int(0)
	}

	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't save key %s, ", key, err)
		return C.int(0)
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()
	if c.HasKey(key) {
		if err := c.UpdateJSON(key, src); err != nil {
			error_dispatch(err, "can't save %s to %s, %s", key, collectionName, err)
			return C.int(0)
		}
	} else {
		if err := c.CreateJSON(key, src); err != nil {
			error_dispatch(err, "can't save %s to %s, %s", key, collectionName, err)
			return C.int(0)
		}
	}
	if err := c.AttachFile(key, key+".xml", bytes.NewReader(xml_src)); err != nil {
		error_dispatch(err, "can't attach %s.xml to %s in %s, %s", key, key, collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// create_record adds an eprint structured record to a dataset collection. It takes a configuration (a map containing a 'DATASET' key/value pair),
// a key (e.g. "20134", the eprint record id as a string) and a record which is the JSON representation of the EPrint XML
//
//export create_record
func create_record(cfg *C.char, cKey *C.char, cRecord *C.char) C.int {
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.int(0)
	}
	key := C.GoString(cKey)
	src := []byte(C.GoString(cRecord))
	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't create key %s, ", key, err)
		return C.int(0)
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	if err := c.CreateJSON(key, src); err != nil {
		error_dispatch(err, "can't create %s in %s, %s", key, collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// read_record returns an eprint structured record from a dataset collection. It takes a configuration (a map containing a 'DATASET' key/value pair),
// a key (e.g. "20134", the eprint record id as a string) and returns a JSON representation of the EPrint XML found in the collection.
//
//export read_record
func read_record(cfg *C.char, cKey *C.char) *C.char {
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.CString("")
	}
	key := C.GoString(cKey)
	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't read key %s, ", key, err)
		return C.CString("")
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()

	src, err := c.ReadJSON(key)
	if err != nil {
		error_dispatch(err, "can't read %s from %s, %s", key, collectionName, err)
		return C.CString("")
	}
	s := fmt.Sprintf("%s", src)
	return C.CString(s)
}

// update_record replaces an eprint structured record in a dataset collection. It takes a configuration (a map containing a 'DATASET' key/value pair),
// a key (e.g. "20134", the eprint record id as a string) and a record which is the JSON represetnation of the EPrint XML
//
//export update_record
func update_record(cfg *C.char, cKey *C.char, cRecord *C.char) C.int {
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.int(0)
	}
	key := C.GoString(cKey)
	src := []byte(C.GoString(cRecord))
	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't update key %s, ", key, err)
		return C.int(0)
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	if err := c.UpdateJSON(key, src); err != nil {
		error_dispatch(err, "can't update %s in %s, %s", key, collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// delete_record removes an eprint structured record from a dataset collection. It takes a configuration (a map containing a 'DATASET' key/value pair),
// a key (e.g. "20134", the eprint record id as a string).
//
//export delete_record
func delete_record(cfg *C.char, cKey *C.char) C.int {
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.int(0)
	}
	key := C.GoString(cKey)
	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't delete key %s, ", key, err)
		return C.int(0)
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()

	if err := c.Delete(key); err != nil {
		error_dispatch(err, "can't delete %s in %s, %s", key, collectionName, err)
		return C.int(0)
	}
	return C.int(1)
}

// record_keys returns a list of keys in a dataset collection. It takes a configuration (a map containing a 'DATASET' key/value pair)
//
func record_keys(cfg *C.char) *C.char {
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.CString("")
	}
	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't get list of keys, ", err)
		return C.CString("")
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()
	keys := c.Keys()
	src, err := json.Marshal(keys)
	if err != nil {
		error_dispatch(err, "failed to marshal list of keys from collection %q, %s", collectionName, err)
		return C.CString("")
	}
	s := fmt.Sprintf("%s", src)
	return C.CString(s)
}

// has_record returns 1 if record is found in a dataset collection, 0 otherwise. It takes a configuration (a map containing a 'DATASET' key/value pair),
// and a key (e.g. "20134", the eprint record id as a string).
//
//export has_record
func has_record(cfg *C.char, cKey *C.char) C.int {
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.int(0)
	}
	key := C.GoString(cKey)
	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't get list of keys, ", err)
		return C.int(0)
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.int(0)
	}
	defer c.Close()
	if c.HasKey(key) == true {
		return C.int(1)
	}
	return C.int(0)
}

// records_to_eprint_xml takes a list of keys matching those stored in a dataset collection and returns an EPrint XML document containing them.
func records_to_eprint_xml(cfg *C.char, cKeys *C.char) *C.char {
	m := map[string]string{}
	cfg_src := []byte(C.GoString(cfg))
	err := json.Unmarshal(cfg_src, &m)
	if err != nil {
		error_dispatch(err, "can't unmarshal config, %s", err)
		return C.CString("")
	}
	collectionName := dsCfg(m)
	if collectionName == "" {
		err := fmt.Errorf("collection name is an empty string")
		error_dispatch(err, "can't get list of keys, ", err)
		return C.CString("")
	}
	c, err := dataset.Open(collectionName)
	if err != nil {
		error_dispatch(err, "failed to open collection %q, %s", collectionName, err)
		return C.CString("")
	}
	defer c.Close()
	src := []byte(C.GoString(cKeys))
	keyList := []string{}
	if err := json.Unmarshal(src, &keyList); err != nil {
		error_dispatch(err, "failed to unmarshal keys for %q, %s", collectionName, err)
		return C.CString("")
	}
	record_src := []string{}
	for _, key := range keyList {
		src, err := c.ReadJSON(key)
		if err != nil {
			error_dispatch(err, "can't read %s from %s, %s", key, collectionName, err)
			return C.CString("")
		}
		eprint := new(eprinttools.EPrint)
		err = json.Unmarshal(src, &eprint)
		if err != nil {
			error_dispatch(err, "can't unmarshal %s from %s, %s", key, collectionName, err)
			return C.CString("")
		}
		src, err = xml.MarshalIndent(eprint, "", "    ")
		if err != nil {
			error_dispatch(err, "can't marshal %s from %s, %s", key, collectionName, err)
			return C.CString("")
		}
		record_src = append(record_src, fmt.Sprintf("%s", src))
	}
	return C.CString(fmt.Sprintf(`<eprint>
%s
</eprint>`, strings.Join(record_src, "\n")))
}

func main() {}

//
// Package eprinttools is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
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
package harvest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	// CaltechLibrary packages
	"github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/eprinttools"
	"github.com/caltechlibrary/rc"
)

var (
	// ExportEPrintDocs if true  include document files as an
	// attachment when calling the export funcs
	ExportEPrintDocs = false

	// EPrintsExportBatchSize sets the summary output frequency when exporting content from E-Prints
	EPrintsExportBatchSize = 1000
)

type byURI []string

func (s byURI) Len() int {
	return len(s)
}

func (s byURI) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byURI) Less(i, j int) bool {
	var (
		a1  int
		a2  int
		err error
	)
	s1 := strings.TrimSuffix(path.Base(s[i]), path.Ext(s[i]))
	s2 := strings.TrimSuffix(path.Base(s[j]), path.Ext(s[j]))
	a1, err = strconv.Atoi(s1)
	if err != nil {
		return false
	}
	a2, err = strconv.Atoi(s2)
	if err != nil {
		return false
	}
	//NOTE: We're creating a descending sort, so a1 should be larger than a2
	return a1 > a2
}

// handleAttachments will attach the EPrintsXML
func handleAttachments(api *eprinttools.EPrintsAPI, c *dataset.Collection, key string, rec *eprinttools.EPrint, xmlSrc []byte) error {
	if ExportEPrintDocs == false {
		c.AttachFile(key, key+".xml", bytes.NewReader(xmlSrc))
		return nil
	}
	rest, err := rc.New(api.URL.String(), api.AuthType, api.Username, api.Secret)
	if err != nil {
		return err
	}
	// Create a temp directory for our harvested files.
	// Harvest documents to tdir, then attach individual files
	// as attachments.
	if rec.Documents != nil && rec.Documents.Length() > 0 {
		tdir, err := ioutil.TempDir(c.Name, "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tdir)
		fNames := []string{}
		fName := path.Join(tdir, key+".xml")
		err = ioutil.WriteFile(fName, xmlSrc, 0666)
		if err != nil {
			return err
		}
		fNames = append(fNames, fName)
		for i := 0; i < rec.Documents.Length(); i++ {
			doc := rec.Documents.IndexOf(i)
			if doc.Files != nil && len(doc.Files) > 0 {
				for _, f := range doc.Files {
					u, err := url.Parse(f.URL)
					if err != nil {
						continue
					}
					src, err := rest.Request("GET", u.Path, map[string]string{})
					if err != nil {
						log.Println(err)
						continue
					}
					fName = path.Join(tdir, f.Filename)
					err = ioutil.WriteFile(fName, src, 0666)
					if err != nil {
						log.Println(err)
						continue
					}
					fNames = append(fNames, fName)
				}
				if len(fNames) > 0 {
					c.AttachFiles(key, fNames...)
				}
			}
		}
	}
	return nil
}

// ExportEPrintsKeyList export a list of eprints from a list of keys
func ExportEPrintsKeyList(api *eprinttools.EPrintsAPI, keys []string, saveKeys string, verbose bool) error {
	var (
		exportedKeys []string
		err          error
		src          []byte
	)

	c, err := dataset.Open(api.Dataset)
	if err != nil {
		return fmt.Errorf("ExportEPrintsKeyList() %s, %s", api.Dataset, err)
	}
	defer c.Close()

	uris := []string{}
	for _, key := range keys {
		uri := fmt.Sprintf("/rest/eprint/%s.xml", strings.TrimSpace(key))
		uris = append(uris, uri)
	}

	pid := os.Getpid()
	uriCount := len(uris)
	count := uriCount
	j := 0 // success count
	k := 0 // error count
	if verbose == true {
		log.Printf("(pid: %d) Exporting %d of %d uris", pid, count, uriCount)
	}
	for i := 0; i < uriCount && i < count; i++ {
		uri := uris[i]
		rec, xmlSrc, err := api.GetEPrint(uri)
		if err != nil {
			if strings.HasPrefix(err.Error(), "WARNING") {
				id := fmt.Sprintf("%d", rec.EPrintID)
				if c.HasKey(id) {
					if err2 := c.Delete(id); err2 == nil {
						log.Printf("(pid: %d) Pruning, %s\n", pid, err)
					} else {
						log.Printf("(pid: %d) Failed to prune %d, %s", pid, id, err2)
					}
				} else {
					log.Printf("(pid: %d) Skipping, %s\n", pid, err)
				}
			} else {
				log.Printf("(pid: %d) Failed, %s\n", pid, err)
			}
			k++
		} else {
			key := fmt.Sprintf("%d", rec.EPrintID)
			src, err = json.Marshal(rec)
			if err != nil {
				log.Printf("(pid: %d) can't marshal key %s, %s", pid, key, err)
			} else {
				// NOTE: Check to see if we're doing an update or create
				if c.HasKey(key) == true {
					err = c.UpdateJSON(key, src)
				} else {
					err = c.CreateJSON(key, src)
				}
			}
			if err == nil {
				if len(saveKeys) > 0 {
					exportedKeys = append(exportedKeys, key)
				}
				// We've exported a record successfully, now update select lists
				j++
			} else {
				if verbose == true {
					log.Printf("(pid: %d) Failed to save eprint %s (%s) to %s, %s\n", pid, key, uri, api.Dataset, err)
				}
				k++
			}
			//NOTE: If ExportEPrintsDocs == true then attach docs and files
			handleAttachments(api, c, key, rec, xmlSrc)
		}
		if verbose == true && (i%EPrintsExportBatchSize) == 0 {
			log.Printf("(pid: %d) %d/%d uri processed, %d exported, %d unexported", pid, i+1, count, j, k)
		}
	}
	if verbose == true {
		log.Printf("(pid: %d) %d/%d uri processed, %d exported, %d unexported", pid, len(uris), count, j, k)
	}
	if len(saveKeys) > 0 {
		if err := ioutil.WriteFile(saveKeys, []byte(strings.Join(exportedKeys, "\n")), 0664); err != nil {
			return fmt.Errorf("Failed to export %s, %s", saveKeys, err)
		}
	}
	return nil
}

// ExportEPrints from highest ID to lowest for cnt. Saves each record in a DB and indexes published ones
func ExportEPrints(api *eprinttools.EPrintsAPI, count int, saveKeys string, verbose bool) error {
	var (
		exportedKeys []string
		err          error
		src          []byte
	)

	c, err := dataset.Open(api.Dataset)
	if err != nil {
		return fmt.Errorf("ExportEPrints() %s, %s", api.Dataset, err)
	}
	defer c.Close()

	uris, err := api.ListEPrintsURI()
	if err != nil {
		return fmt.Errorf("Export %s failed, %s", api.URL.String(), err)
	}

	// NOTE: I am sorting the URI by decscending ID number so that the
	// newest articles are exported first
	sort.Sort(byURI(uris))

	pid := os.Getpid()
	uriCount := len(uris)
	if count < 0 {
		count = uriCount
	}
	j := 0 // success count
	k := 0 // error count
	if verbose == true {
		log.Printf("(pid: %d) Exporting %d of %d uris", pid, count, uriCount)
	}
	i := 0
	for i = 0; i < uriCount && i < count; i++ {
		uri := uris[i]
		rec, xmlSrc, err := api.GetEPrint(uri)
		if err != nil {
			if strings.HasPrefix(err.Error(), "WARNING") {
				id := fmt.Sprintf("%d", rec.EPrintID)
				if c.HasKey(id) {
					if err2 := c.Delete(id); err2 == nil {
						log.Printf("(pid: %d) Pruning, %s\n", pid, err)
					} else {
						log.Printf("(pid: %d) Failed to prune %d, %s", pid, id, err2)
					}
				} else {
					log.Printf("(pid: %d) Skipping, %s\n", pid, err)
				}
			} else {
				log.Printf("(pid: %d) Failed, %s\n", pid, err)
			}
			k++
		} else {
			key := fmt.Sprintf("%d", rec.EPrintID)
			src, err = json.Marshal(rec)
			if err != nil {
				log.Printf("(pid: %d) Can't marshal key %s, %s", pid, key, err)
			} else {
				// NOTE: Check to see if we're doing an update or create
				if c.HasKey(key) == true {
					err = c.UpdateJSON(key, src)
				} else {
					err = c.CreateJSON(key, src)
				}
			}
			if err == nil {
				if len(saveKeys) > 0 {
					exportedKeys = append(exportedKeys, key)
				}
				// We've exported a record successfully, now update select lists
				j++
			} else {
				if verbose == true {
					log.Printf("(pid: %d) Failed to save eprint %s (%s) to %s, %s\n", pid, key, uri, api.Dataset, err)
				}
				k++
			}
			//NOTE: If ExportEPrintsDocs == true then attach docs and files
			handleAttachments(api, c, key, rec, xmlSrc)
		}
		if verbose == true && (i%EPrintsExportBatchSize) == 0 {
			log.Printf("(pid: %d) %d/%d uri processed, %d exported, %d unexported", pid, i+1, count, j, k)
		}
	}
	if verbose == true {
		log.Printf("(pid: %d) %d/%d uri processed, %d exported, %d unexported", pid, i, count, j, k)
	}
	if len(saveKeys) > 0 {
		if err := ioutil.WriteFile(saveKeys, []byte(strings.Join(exportedKeys, "\n")), 0664); err != nil {
			return fmt.Errorf("Failed to export %s, %s", saveKeys, err)
		}
	}
	return nil
}

// ExportModifiedEPrints returns a list of ids modified in one or between the start, end times
func ExportModifiedEPrints(api *eprinttools.EPrintsAPI, start, end time.Time, saveKeys string, verbose bool) error {
	var (
		exportedKeys []string
		err          error
		src          []byte
	)

	pid := os.Getpid()

	c, err := dataset.Open(api.Dataset)
	if err != nil {
		return fmt.Errorf("ExportModifiedEPrints() %s, %s", api.Dataset, err)
	}
	defer c.Close()

	uris, err := api.ListModifiedEPrintsURI(start, end, verbose)
	if err != nil {
		return fmt.Errorf("Export modified from %s to %s failed, %s", start, end, err)
	}

	// NOTE: I am sorting the URI by decscending ID number so that the
	// newest articles are exported first
	sort.Sort(byURI(uris))

	count := len(uris)
	j := 0 // success count
	k := 0 // error count
	if verbose == true {
		log.Printf("(pid: %d) Exporting %d uris", pid, count)
	}
	i := 0
	for i = 0; i < count; i++ {
		uri := uris[i]
		rec, xmlSrc, err := api.GetEPrint(uri)
		if err != nil {
			if verbose == true {
				log.Printf("(pid: %d) Failed to get %s, %s\n", pid, uri, err)
			}
			k++
		} else {
			key := fmt.Sprintf("%d", rec.EPrintID)
			src, err = json.Marshal(rec)
			if err != nil {
				log.Printf("(pid: %d) Can't marshal key %s, %s", pid, key, err)
			} else {
				// NOTE: Check to see if we're doing an update or create
				if c.HasKey(key) == true {
					err = c.UpdateJSON(key, src)
				} else {
					err = c.CreateJSON(key, src)
				}
			}
			if err == nil {
				if len(saveKeys) > 0 {
					exportedKeys = append(exportedKeys, key)
				}
				// We've exported a record successfully, now update select lists
				j++
			} else {
				if verbose == true {
					log.Printf("(pid: %d) Failed to save eprint %s (%s) to %s, %s\n", pid, key, uri, api.Dataset, err)
				}
				k++
			}
			//NOTE: If ExportEPrintsDocs == true then attach docs and files
			handleAttachments(api, c, key, rec, xmlSrc)
		}
		if verbose == true && (i%EPrintsExportBatchSize) == 0 {
			log.Printf("(pid: %d) %d/%d uri processed, %d exported, %d unexported", pid, i+1, count, j, k)
		}
	}
	if verbose == true {
		log.Printf("(pid: %d) %d/%d uri processed, %d exported, %d unexported", pid, i, count, j, k)
	}
	if len(saveKeys) > 0 {
		if err := ioutil.WriteFile(saveKeys, []byte(strings.Join(exportedKeys, "\n")), 0664); err != nil {
			return fmt.Errorf("Failed to export %s, %s", saveKeys, err)
		}
	}
	return nil
}

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
package eprinttools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	// CaltechLibrary packages
	"github.com/caltechlibrary/dataset"
)

// ExportEPrintsKeyList export a list of eprints from a list of keys
func (api *EPrintsAPI) ExportEPrintsKeyList(keys []string, saveKeys string, verbose bool) error {
	var (
		exportedKeys []string
		err          error
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

	uriCount := len(uris)
	count := uriCount
	j := 0 // success count
	k := 0 // error count
	if verbose == true {
		log.Printf("Exporting %d of %d uris", count, uriCount)
	}
	for i := 0; i < uriCount && i < count; i++ {
		uri := uris[i]
		rec, xmlSrc, err := api.GetEPrint(uri)
		if err != nil {
			log.Printf("Failed, %s\n", err)
			k++
		} else {
			key := fmt.Sprintf("%d", rec.EPrintID)
			// NOTE: Check to see if we're doing an update or create
			if c.HasKey(key) == true {
				err = c.UpdateFrom(key, rec)
			} else {
				err = c.CreateFrom(key, rec)
			}
			if err == nil {
				if len(saveKeys) > 0 {
					exportedKeys = append(exportedKeys, key)
				}
				// We've exported a record successfully, now update select lists
				j++
			} else {
				if verbose == true {
					log.Printf("Failed to save eprint %s (%s) to %s, %s\n", key, uri, api.Dataset, err)
				}
				k++
			}
			c.AttachFile(key, key+".xml", bytes.NewReader(xmlSrc))
		}
		if verbose == true && (i%EPrintsExportBatchSize) == 0 {
			log.Printf("%d/%d uri processed, %d exported, %d unexported", i+1, count, j, k)
		}
	}
	if verbose == true {
		log.Printf("%d/%d uri processed, %d exported, %d unexported", len(uris), count, j, k)
	}
	if len(saveKeys) > 0 {
		if err := ioutil.WriteFile(saveKeys, []byte(strings.Join(exportedKeys, "\n")), 0664); err != nil {
			return fmt.Errorf("Failed to export %s, %s", saveKeys, err)
		}
	}
	return nil
}

// ExportEPrints from highest ID to lowest for cnt. Saves each record in a DB and indexes published ones
func (api *EPrintsAPI) ExportEPrints(count int, saveKeys string, verbose bool) error {
	var (
		exportedKeys []string
		err          error
	)

	c, err := dataset.InitCollection(api.Dataset)
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

	uriCount := len(uris)
	if count < 0 {
		count = uriCount
	}
	j := 0 // success count
	k := 0 // error count
	if verbose == true {
		log.Printf("Exporting %d of %d uris", count, uriCount)
	}
	for i := 0; i < uriCount && i < count; i++ {
		uri := uris[i]
		rec, xmlSrc, err := api.GetEPrint(uri)
		if err != nil {
			log.Printf("Failed, %s\n", err)
			k++
		} else {
			key := rec.ID
			// NOTE: Check to see if we're doing an update or create
			if c.HasKey(key) == true {
				err = c.UpdateFrom(key, rec)
			} else {
				err = c.CreateFrom(key, rec)
			}
			if err == nil {
				if len(saveKeys) > 0 {
					exportedKeys = append(exportedKeys, key)
				}
				// We've exported a record successfully, now update select lists
				j++
			} else {
				if verbose == true {
					log.Printf("Failed to save eprint %s (%s) to %s, %s\n", key, uri, api.Dataset, err)
				}
				k++
			}
			c.AttachFile(key, key+".xml", bytes.NewReader(xmlSrc))
		}
		if verbose == true && (i%EPrintsExportBatchSize) == 0 {
			log.Printf("%d/%d uri processed, %d exported, %d unexported", i+1, count, j, k)
		}
	}
	if verbose == true {
		log.Printf("%d/%d uri processed, %d exported, %d unexported", len(uris), count, j, k)
	}
	if len(saveKeys) > 0 {
		if err := ioutil.WriteFile(saveKeys, []byte(strings.Join(exportedKeys, "\n")), 0664); err != nil {
			return fmt.Errorf("Failed to export %s, %s", saveKeys, err)
		}
	}
	return nil
}

// ExportModifiedEPrints returns a list of ids modified in one or between the start, end times
func (api *EPrintsAPI) ExportModifiedEPrints(start, end time.Time, saveKeys string, verbose bool) error {
	var (
		exportedKeys []string
		err          error
	)

	c, err := dataset.InitCollection(api.Dataset)
	if err != nil {
		return fmt.Errorf("ExportModifiedEPrints() %s, %s", api.Dataset, err)
	}
	defer c.Close()

	uris, err := api.ListModifiedEPrintURI(start, end, verbose)
	if err != nil {
		return fmt.Errorf("Export modified %s to %s failed, %s", start, end, err)
	}

	// NOTE: I am sorting the URI by decscending ID number so that the
	// newest articles are exported first
	sort.Sort(byURI(uris))

	count := len(uris)
	j := 0 // success count
	k := 0 // error count
	if verbose == true {
		log.Printf("Exporting %d uris", count)
	}
	for i := 0; i < count; i++ {
		uri := uris[i]
		rec, xmlSrc, err := api.GetEPrint(uri)
		if err != nil {
			if verbose == true {
				log.Printf("Failed to get %s, %s\n", uri, err)
			}
			k++
		} else {
			key := rec.ID
			// NOTE: Check to see if we're doing an update or create
			if c.HasKey(key) == true {
				err = c.UpdateFrom(key, rec)
			} else {
				err = c.CreateFrom(key, rec)
			}
			if err == nil {
				if len(saveKeys) > 0 {
					exportedKeys = append(exportedKeys, key)
				}
				// We've exported a record successfully, now update select lists
				j++
			} else {
				if verbose == true {
					log.Printf("Failed to save eprint %s (%s) to %s, %s\n", key, uri, api.Dataset, err)
				}
				k++
			}
			c.AttachFile(key, key+".xml", bytes.NewReader(xmlSrc))
		}
		if verbose == true && (i%EPrintsExportBatchSize) == 0 {
			log.Printf("%d/%d uri processed, %d exported, %d unexported", i+1, count, j, k)
		}
	}
	if verbose == true {
		log.Printf("%d/%d uri processed, %d exported, %d unexported", len(uris), count, j, k)
	}
	if len(saveKeys) > 0 {
		if err := ioutil.WriteFile(saveKeys, []byte(strings.Join(exportedKeys, "\n")), 0664); err != nil {
			return fmt.Errorf("Failed to export %s, %s", saveKeys, err)
		}
	}
	return nil
}

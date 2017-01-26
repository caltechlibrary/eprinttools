//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
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
package epgo

import (
	"fmt"
	"log"
	"sort"
	"strings"

	// CaltechLibrary packages
	"github.com/caltechlibrary/dataset"
)

// ExportEPrints from highest ID to lowest for cnt. Saves each record in a DB and indexes published ones
func (api *EPrintsAPI) ExportEPrints(count int) error {
	var errs []error
	c, err := dataset.Create(api.Dataset, dataset.GenerateBucketNames(dataset.DefaultAlphabet, 2))
	failCheck(err, fmt.Sprintf("ExportEPrints() %s, %s", api.Dataset, err))
	defer c.Close()

	sLists := map[string]*dataset.SelectList{}
	for _, name := range slNames {
		c.Clear(name)
		if sLists[name], err = c.Select(name); err != nil {
			return fmt.Errorf("Can't create empty select list for %s, %s", name, err)
		}
	}

	uris, err := api.ListEPrintsURI()
	failCheck(err, fmt.Sprintf("Export %s failed, %s", api.URL.String(), err))

	//NOTE: I am sorting the URI by decscending ID number so that the newest articles
	// are exported first
	sort.Sort(byURI(uris))

	uriCount := len(uris)
	if count < 0 {
		count = uriCount
	}
	j := 0 // success count
	k := 0 // error count
	log.Printf("Exporting %d of %d uris", count, uriCount)
	for i := 0; i < uriCount && i < count; i++ {
		uri := uris[i]
		rec, err := api.GetEPrint(uri)
		if err != nil {
			log.Printf("Failed to get %s, %s\n", uri, err)
			k++
		} else {
			key := fmt.Sprintf("%d", rec.ID)
			err := c.Create(key, rec)
			if err != nil {
				errs = append(errs, err)
			} else {
				// We've exported a record successfully, now update select lists
				j++

				// Update pubDates select list
				dt := normalizeDate(rec.Date)
				if rec.DateType == "published" && rec.Date != "" {
					sLists["pubDate"].Push(fmt.Sprintf("%s%s%d", dt, indexDelimiter, rec.ID))
				}
				// Update localGroups select list
				if len(rec.LocalGroup) > 0 {
					for _, grp := range rec.LocalGroup {
						grp = strings.TrimSpace(grp)
						if len(grp) > 0 {
							sLists["localGroup"].Push(fmt.Sprintf("%s%s%s%s%d", grp, indexDelimiter, dt, indexDelimiter, rec.ID))
						}
					}
				}
				// Update orcids, isnis and authors select list
				if len(rec.Creators) > 0 {
					for _, person := range rec.Creators {
						orcid := strings.TrimSpace(person.ORCID)
						isni := strings.TrimSpace(person.ISNI)
						author := fmt.Sprintf("%s, %s", strings.TrimSpace(person.Family), strings.TrimSpace(person.Given))
						if len(orcid) > 0 {
							sLists["orcid"].Push(fmt.Sprintf("%s%s%s%s%d", orcid, indexDelimiter, dt, indexDelimiter, rec.ID))
						}
						if len(isni) > 0 {
							sLists["isni"].Push(fmt.Sprintf("%s%s%s%s%d", isni, indexDelimiter, dt, indexDelimiter, rec.ID))
						}
						if len(author) > 0 {
							sLists["author"].Push(fmt.Sprintf("%s%s%s%s%d", author, indexDelimiter, dt, indexDelimiter, rec.ID))
						}
					}
				}

				// Add funders and grantNumbers to select lists
				if len(rec.Funders) > 0 {
					for _, funder := range rec.Funders {
						funderName := strings.TrimSpace(funder.Agency)
						grantNumber := strings.TrimSpace(funder.GrantNumber)
						if len(funderName) > 0 {
							sLists["funder"].Push(fmt.Sprintf("%s%s%s%s%d", funderName, indexDelimiter, dt, indexDelimiter, rec.ID))
						}
						if len(grantNumber) > 0 {
							sLists["grantNumber"].Push(fmt.Sprintf("%s%s%s%s%d", grantNumber, indexDelimiter, dt, indexDelimiter, rec.ID))
						}
					}
				}
			}

			if err != nil {
				log.Printf("Failed to save eprint %s, %s\n", uri, err)
				k++
			}
		}
		if (i % EPrintsExportBatchSize) == 0 {
			log.Printf("%d uri processed, %d exported, %d unexported", i+1, j, k)
		}
	}
	log.Printf("%d uri processed, %d exported, %d unexported", len(uris), j, k)

	return nil
}
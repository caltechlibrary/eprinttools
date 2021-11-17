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
package clsrules

//
// clsrules.go (Caltech Library Specific Rules) is a package for
// implementing Caltech Library Specific features to processing
// and creating EPrint XML. Currently these include things like
// trimming prefixed "The " from titles, dropping series information,
// changing how the date is derived and very idiosencratic handling
// of Author and DOI references.
//

import (
	"fmt"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

// doiAsRelatedURL - Caltech Library stores the EPrint's DOI in the related URL for historical reasons
// If a DOI is in EPrint.DOI then it needs to migrate to EPrint.RelatedURLItemList as the initial item.
// Returns a revised URL list and boolean true if list was modified to include doi
func doiAsRelatedURL(doi string, description string, relatedURLs *eprinttools.RelatedURLItemList) (*eprinttools.RelatedURLItemList, bool) {
	if doi != "" {
		if relatedURLs == nil {
			relatedURLs = new(eprinttools.RelatedURLItemList)
		}
		entry := new(eprinttools.Item)
		entry.Type = "doi"
		if strings.Contains(doi, "://") {
			entry.URL = doi
		} else {
			entry.URL = fmt.Sprintf("https://doi.org/%s", doi)
		}
		if description != "" {
			entry.Description = description
		}

		//NOTE: doi needs to be inserted in the initial position
		newRelatedURLs := new(eprinttools.RelatedURLItemList)
		newRelatedURLs.Append(entry)
		if len(relatedURLs.Items) > 0 {
			newRelatedURLs.Items = append(newRelatedURLs.Items, relatedURLs.Items...)
		}
		return newRelatedURLs, true
	}
	return relatedURLs, false
}

// normalizeRelatedURLDescriptions
func normalizeRelatedURLDescriptions(relatedURLs *eprinttools.RelatedURLItemList) (*eprinttools.RelatedURLItemList, bool) {
	changed := false
	for i, item := range relatedURLs.Items {
		switch item.Description {
		case "book_section":
			item.Description = "Book Section"
			relatedURLs.Items[i] = item
			changed = true
		case "article":
			item.Description = "Article"
			relatedURLs.Items[i] = item
			changed = true
		}
	}
	return relatedURLs, changed
}

// trimTitle removes any leading "The " and trims spaces from title.
func trimTitle(s string) string {
	if strings.HasPrefix(strings.ToLower(s), "the ") {
		s = s[4:]
	}
	return strings.TrimSpace(s)
}

// trimNumberString, DR-46, George would like leading zeros in issue
// numbers trimmed.
func trimNumberString(s string) string {
	if strings.HasPrefix(s, "0") {
		s = strings.TrimLeft(s, "0")
	}
	return s
}

// normalizeCreators clears the creator list when there are more than
// 30 authors otherwise normalizes the content. If the list changes
// through normalization a new list and bool value of true is returned,
// otherwise the original list and false is returned.
// Limiting Creators to 30 is per George's email from the
// 1Science Load project.
func normalizeCreators(creators *eprinttools.CreatorItemList) (*eprinttools.CreatorItemList, bool) {
	// If more than 30 creators just dump the list and return an empty one
	if len(creators.Items) > 30 {
		return nil, true
	}
	return creators, false
}

func ClearRuleSet() map[string]bool {
	return map[string]bool{
		"trim_title":            false,
		"trim_volume":           false,
		"trim_number":           false,
		"prune_creators":        false,
		"prune_series":          false,
		"doi_as_related_url":    false,
		"normalize_related_url": false,
		"normalize_publisher":   false,
		"normalize_publication": false,
	}
}

func UseCLSRules() map[string]bool {
	return map[string]bool{
		// Conform titles to Caltech's practices, e.g. Trim "The" from titles
		"trim_title": true,
		// Conform Volume value per George and DR-46
		"trim_volume": true,
		// Conform Number value per George and DR-46
		"trim_number": true,
		// NOTE: Per Tools Incubator meeting discussion 2020-02-18
		// between George and Joy we're dropping the limitting of
		// the number of authors into EPrints/CaltechAUTHORS.
		// Normalize Creators and apply George's rules
		"prune_creators": false,
		// Caltech Library doesn't import series information
		"prune_series": true,
		// Apply Caltech Library's pecular DOI assignment behavior
		// As of July 2021 DOI should go in the DOI field not related URL.
		// RSD - 2021-07-30
		"doi_as_related_url": false,
		// Normalize related URL descriptions
		"normalize_related_url": true,
		// Normalize Publisher name from ISSN
		"normalize_publisher": true,
		// Normalize Publication from ISSN
		"normalize_publication": true,
	}
}

// Apply applies the current set of Caltech Library customizations
// to cross walked records to EPrints XML.
func Apply(eprintsList *eprinttools.EPrints, ruleSet map[string]bool) (*eprinttools.EPrints, error) {
	for i, eprint := range eprintsList.EPrint {
		changed := false
		for name, enabled := range ruleSet {
			if enabled {
				switch name {
				case "trim_title":
					if title := trimTitle(eprint.Title); title != eprint.Title {
						eprint.Title = title
						changed = true
					}
				case "trim_volume":
					if volNo := trimNumberString(eprint.Volume); volNo != eprint.Volume {
						eprint.Volume = volNo
						changed = true
					}
				case "trim_number":
					if no := trimNumberString(eprint.Number); no != eprint.Number {
						eprint.Number = no
						changed = true
					}
				case "prune_creators":
					if eprint.Creators != nil && len(eprint.Creators.Items) > 0 {
						if creators, hasChanged := normalizeCreators(eprint.Creators); hasChanged {
							eprint.Creators = creators
							changed = true
						}
					}
				case "prune_series":
					if eprint.Series != "" {
						eprint.Series = ""
						changed = true
					}
				case "doi_as_related_url":
					if eprint.DOI != "" {
						if relatedURLs, hasChanged := doiAsRelatedURL(eprint.DOI, eprint.Type, eprint.RelatedURL); hasChanged {
							eprint.RelatedURL = relatedURLs
							eprint.DOI = ""
							changed = true
						}
					}
				case "normalize_related_url":
					if eprint.RelatedURL != nil {
						if relatedURLs, hasChanged := normalizeRelatedURLDescriptions(eprint.RelatedURL); hasChanged {
							eprint.RelatedURL = relatedURLs
							changed = true
						}
					}
				case "normalize_publisher":
					if eprint.ISSN != "" {
						if publisher, ok := issnPublisher[eprint.ISSN]; ok == true {
							eprint.Publisher = publisher
							changed = true
						}
					}
				case "normalize_publication":
					// Normalize Publisher name and Publication from ISSN
					if eprint.ISSN != "" {
						if publication, ok := issnPublication[eprint.ISSN]; ok == true {
							eprint.Publication = publication
							changed = true
						}
					}
				}
			}
		}

		// If we've changed the eprint record update it.
		if changed {
			eprintsList.EPrint[i] = eprint
		}
	}
	return eprintsList, nil
}

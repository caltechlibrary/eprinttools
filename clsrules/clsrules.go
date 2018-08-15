//
// clsrules.go (Caltech Library Specific Rules) is a package for
// implementing Caltech Library Specific features to processing
// and creating EPrint XML. Currently these include things like
// trimming prefixed "The " from titles, dropping series information,
// changing how the date is derived and very idiosencratic handling
// of Author and DOI references.
//
package clsrules

import (
	"fmt"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

// migrateDOI - Caltech Library stores the EPrint's DOI in the related URL for historical reasons
// If a DOI is in EPrint.DOI then it needs to migrate to EPrint.RelatedURLItemList as the initial item.
// Returns a revised URL list and boolean true if list was modified to include doi
func migrateDOI(doi string, description string, relatedURLs *eprinttools.RelatedURLItemList) (*eprinttools.RelatedURLItemList, bool) {
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
		newRelatedURLs.AddItem(entry)
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

// normalizeCreators clears the creator list when there are more than
// 30 authors otherwise normalizes the content. If the list changes
// through normalization a new list and bool value of true is returned,
// otherwise the original list and false is returned.
// Limiting Creators to 30 is per George's email from the
// 1Science Load project.
func normalizeCreators(creators *eprinttools.CreatorItemList) (*eprinttools.CreatorItemList, bool) {
	// If more than 30 creators just dump the list and return an empty one
	if len(creators.Items) > 30 {
		return new(eprinttools.CreatorItemList), true
	}
	return creators, false
}

func Apply(eprintsList *eprinttools.EPrints) (*eprinttools.EPrints, error) {
	// Trim "The" from titles
	for i, eprint := range eprintsList.EPrint {
		changed := false
		// Conform titles to Caltech's practices
		if title := trimTitle(eprint.Title); title != eprint.Title {
			eprint.Title = title
			changed = true
		}

		// Normalize Creators and apply George's rules
		if len(eprint.Creators.Items) > 0 {
			if creators, hasChanged := normalizeCreators(eprint.Creators); hasChanged {
				eprint.Creators = creators
				changed = true
			}
		}

		// Caltech Library doesn't import series information
		if eprint.Series != "" {
			eprint.Series = ""
			changed = true
		}

		// Handle Caltech Library's pecular DOI assignment behavior
		if eprint.DOI != "" {
			if relatedURLs, hasChanged := migrateDOI(eprint.DOI, eprint.Type, eprint.RelatedURL); hasChanged {
				eprint.RelatedURL = relatedURLs
				eprint.DOI = ""
				changed = true
			}
		}

		// Normalize related URL descriptions
		if eprint.RelatedURL != nil {
			if relatedURLs, hasChanged := normalizeRelatedURLDescriptions(eprint.RelatedURL); hasChanged {
				eprint.RelatedURL = relatedURLs
				changed = true
			}
		}

		// Normalize Publisher name and Publication from ISSN
		if eprint.ISSN != "" {
			if publisher, ok := issnPublisher[eprint.ISSN]; ok == true {
				eprint.Publisher = publisher
				changed = true
			}
			if publication, ok := issnPublication[eprint.ISSN]; ok == true {
				eprint.Publication = publication
				changed = true
			}
		}

		// If we've changed the eprint record update it.
		if changed {
			eprintsList.EPrint[i] = eprint
		}

	}
	return eprintsList, nil
}

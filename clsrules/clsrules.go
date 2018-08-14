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
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/eprinttools"
)

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

func Apply(eprintsList *eprinttools.EPrints) (*eprinttools.EPrints, error) {
	// Trim "The" from titles
	for i, eprint := range eprintsList.EPrint {
		if title := trimTitle(eprint.Title); title != eprint.Title {
			eprintsList.EPrint[i].Title = title
		}
		if eprint.Series != "" {
			eprint.Series = ""
		}
		if eprint.RelatedURL != nil {
			if relatedURLs, hasChanged := normalizeRelatedURLDescriptions(eprint.RelatedURL); hasChanged {
				eprint.RelatedURL = relatedURLs
			}
		}
	}
	return eprintsList, nil
}

package eprinttools

import (
	"fmt"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/crossrefapi"
)

// normalizeCrossRefToLocalGroup takes an affiliation from CrossRef author
// and attempts to determine if it is a Caltech Group and returns
// the normalized name or an empty string.
func normalizeCrossRefToLocalGroup(s string) string {
	s = strings.ToLower(s)
	switch {
	case strings.Contains(s, "california institute of technology"):
		i := strings.Index(s, "california institute of technology")
		if i > -1 {
			s = strings.TrimSuffix(strings.TrimSpace(s[0:i-1]), ";")
		}
		return strings.TrimSpace(s)
	case strings.Contains(s, "california institution of technology"):
		i := strings.Index(s, "california institution of technology")
		if i > -1 {
			s = strings.TrimSuffix(strings.TrimSpace(s[0:i-1]), ";")
		}
		return strings.TrimSpace(s)
	case strings.Contains(s, "caltech"):
	case strings.Contains(s, "caltech"):
		i := strings.Index(s, "caltech")
		if i > -1 {
			s = strings.TrimSuffix(strings.TrimSpace(s[0:i-1]), ";")
		}
		return strings.TrimSpace(s)
	case strings.Contains(s, "ligo"):
		return "LIGO"
	}
	return ""
}

// normalizeCrossRefType converts content type from CrossRef
// to Authors (e.g. "journal-article" to "article")
func normalizeCrossRefType(s string) string {
	switch strings.ToLower(s) {
	case "proceedings-article":
		//NOTE: This seems vary idiosyncratic to CaltechAUTHORS
		return "book_section"
	case "journal-article":
		return "article"
	case "book-chapter":
		return "book_section"
	default:
		return s
	}
}

// CrossRefWorksToEPrint takes a works object from the CrossRef API
// and maps the fields into an EPrint struct return a new struct or
// error.
func CrossRefWorksToEPrint(obj crossrefapi.Object) (*EPrint, error) {
	eprint := new(EPrint)
	// Type
	if s, ok := indexInto(obj, "message", "type"); ok == true {
		eprint.Type = normalizeCrossRefType(fmt.Sprintf("%s", s))
	} else {
		return nil, fmt.Errorf("Can't find type in object")
	}
	// Title
	if a, ok := indexInto(obj, "message", "title"); ok == true {
		if len(a.([]interface{})) > 0 {
			eprint.Title = fmt.Sprintf("%s", a.([]interface{})[0].(string))
		}
	}

	// NOTE: Assuming IsPublished is true given that we're talking to
	// CrossRef API which holds published content.
	// IsPublished
	eprint.IsPublished = "pub"

	// Publisher
	if s, ok := indexInto(obj, "message", "publisher"); ok == true {
		eprint.Publisher = fmt.Sprintf("%s", s)
	}
	// Publication
	if eprint.Type == "article" {
		if l, ok := indexInto(obj, "message", "container-title"); ok == true {
			if len(l.([]interface{})) > 0 {
				eprint.Publication = l.([]interface{})[0].(string)
			}
		}
	}
	// Series
	if eprint.Type == "book" {
		if l, ok := indexInto(obj, "message", "container-title"); ok == true {
			if len(l.([]interface{})) > 0 {
				eprint.Series = l.([]interface{})[0].(string)
			}
		}
	}
	if l, ok := indexInto(obj, "message", "short-container-title"); ok == true {
		if len(l.([]interface{})) > 0 {
			eprint.Series = l.([]interface{})[0].(string)
		}
	}

	// Volume
	if eprint.Type == "article" {
		if s, ok := indexInto(obj, "message", "volume"); ok == true {
			eprint.Volume = fmt.Sprintf("%s", s)
		}
		// Number
		if s, ok := indexInto(obj, "message", "journal-issue", "issue"); ok == true {
			eprint.Number = fmt.Sprintf("%s", s)
		}

	}

	// PlaceOfPub taken from publisher-location in CrossRef
	if s, ok := indexInto(obj, "message", "publisher-location"); ok == true {
		eprint.PlaceOfPub = fmt.Sprintf("%s", s)
	}

	// PageRange
	// Pages
	if s, ok := indexInto(obj, "message", "page"); ok == true {
		if strings.Contains(s.(string), "-") {
			eprint.PageRange = fmt.Sprintf("%s", s)
		} else {
			eprint.Pages = fmt.Sprintf("%s", s)
		}
	}

	// ISBN
	if a, ok := indexInto(obj, "message", "ISBN"); ok == true {
		if len(a.([]interface{})) > 0 {
			s := a.([]interface{})[0]
			eprint.ISBN = fmt.Sprintf("%s", s)
		}
	}

	// ISSN
	if a, ok := indexInto(obj, "message", "ISSN"); ok == true {
		if len(a.([]interface{})) > 0 {
			eprint.ISSN = fmt.Sprintf("%s", a.([]interface{})[0])
		}
	}

	// NOTE: This doesn't appear to be used by CaltechAUTHORS for full book
	// BookTitle
	if eprint.Title != "" && eprint.Type == "book" {
		eprint.BookTitle = eprint.Title
	}

	// Funders
	if a, ok := indexInto(obj, "message", "funder"); ok == true {
		eprint.Funders = new(FunderItemList)
		for _, item := range a.([]interface{}) {
			entry := new(Item)
			m := item.(map[string]interface{})
			if name, ok := indexInto(m, "name"); ok == true && name != "N/A" {
				entry.Agency = fmt.Sprintf("%s", name)
			}
			if a2, ok := indexInto(m, "award"); ok == true && a2 != "N/A" {
				if len(a2.([]interface{})) > 0 {
					entry.GrantNumber = fmt.Sprintf("%s", a2.([]interface{})[0])
				}
			}
			if entry.Agency != "" || entry.GrantNumber != "" {
				eprint.Funders.AddItem(entry)
			}
		}
	}

	// NOTE: Caltech Library puts the DOI in the related URL field rather than
	// in EPrint's default location. This code puts the DOI in the default
	// location. If you need Caltech Library's bahavior use clsrules.Apply()
	// to conform to that regime.
	if doi, ok := indexInto(obj, "message", "DOI"); ok == true {
		eprint.DOI = doi.(string)
	}
	if l, ok := indexInto(obj, "message", "update-to"); ok == true {
		for _, o := range l.([]interface{}) {
			m := o.(map[string]interface{})
			if newDoi, ok := indexInto(m, "DOI"); ok == true && newDoi != "" {
				dt, _ := indexInto(m, "updated", "date-time")
				when := dt.(string)
				l, _ := indexInto(m, "label")
				label := l.(string)
				if len(when) > 10 {
					when = when[0:10]
				}
				entry := new(Item)
				entry.Type = "doi"
				entry.URL = fmt.Sprintf("https://doi.org/%s", newDoi)
				entry.Description = fmt.Sprintf("%s, %s", label, when)
				eprint.RelatedURL.AddItem(entry)
			}
		}
	}

	// RelatedURLs (links in message of CrossRef works object)
	if l, ok := indexInto(obj, "message", "link"); ok == true {
		if eprint.RelatedURL == nil {
			eprint.RelatedURL = new(RelatedURLItemList)
		}
		for _, o := range l.([]interface{}) {
			entry := new(Item)
			if s, ok := indexInto(o.(map[string]interface{}), "URL"); ok == true {
				entry.URL = s.(string)
			}
			// NOTE: Related URL Type is not related to mime-type,
			// import related URLs without type information.
			if s, ok := indexInto(o.(map[string]interface{}), "type"); ok == true {
				entry.Type = s.(string)
			}
			if len(entry.URL) > 0 { //&& len(entry.Type) > 0 {
				eprint.RelatedURL.AddItem(entry)
			}
		}
	}

	// NOTE: Assuming date is published given we're talking to CrossRef
	// Date. We prefer issued date but fallback to created date.
	eprint.DateType = "published"
	if issued, ok := indexInto(obj, "message", "issued", "date-time"); ok == true {
		// DateType
		eprint.Date = fmt.Sprintf("%s", issued)
	} else if created, ok := indexInto(obj, "message", "created", "date-time"); ok == true {
		// DateType
		eprint.Date = fmt.Sprintf("%s", created)
	}
	if len(eprint.Date) > 10 {
		eprint.Date = eprint.Date[0:10]
	}

	// Authors list
	if l, ok := indexInto(obj, "message", "author"); ok == true {
		creators := new(CreatorItemList)
		corpCreators := new(CorpCreatorItemList)
		localGroups := new(LocalGroupItemList)
		for _, entry := range l.([]interface{}) {
			author := entry.(map[string]interface{})
			item := new(Item)
			item.Name = new(Name)
			if affiliations, ok := author["affiliation"]; ok == true {
				for _, affiliation := range affiliations.([]interface{}) {
					if s, ok := indexInto(affiliation.(map[string]interface{}), "name"); ok == true {
						group := new(Item)
						group.Value = normalizeCrossRefToLocalGroup(s.(string))
						if group.Value != "" {
							localGroups.AddItem(group)
						}
					}
				}
			}
			if orcid, ok := author["ORCID"]; ok == true {
				item.ORCID = orcid.(string)
				if strings.HasPrefix(orcid.(string), "http://orcid.org/") {
					item.ORCID = strings.TrimPrefix(orcid.(string), "http://orcid.org/")
				}
				if strings.HasPrefix(orcid.(string), "https://orcid.org/") {
					item.ORCID = strings.TrimPrefix(orcid.(string), "http://orcid.org/")
				}
			}
			if family, ok := author["family"]; ok == true {
				item.Name.Family = family.(string)
			}
			if given, ok := author["given"]; ok == true {
				item.Name.Given = given.(string)
			}
			//NOTE: if as have a 'name' then we'll add it to
			// as a corp_creators
			if name, ok := author["name"]; ok == true {
				item.Name.Value = name.(string)
				if strings.HasPrefix(item.Name.Value, "(") && strings.HasSuffix(item.Name.Value, ")") {
					item.Name.Value = strings.TrimSuffix(strings.TrimPrefix(item.Name.Value, "("), ")")
				}
			}
			if item.Name.Given != "" && item.Name.Family != "" {
				creators.AddItem(item)
			}
			if item.Name.Value != "" {
				corpCreators.AddItem(item)
			}
		}
		if len(creators.Items) > 0 {
			eprint.Creators = creators
		}
		if len(corpCreators.Items) > 0 {
			eprint.CorpCreators = corpCreators
		}
	}

	// Edition
	//FIXME: Need to find value in CrossRef works metadata for this

	// FullTextStatus
	//FIXME: Need to find value in CrossRef works metadata for this

	// Keywords
	//FIXME: Need to find value in CrossRef works metadata for this

	// Note
	//FIXME: Need to find value in CrossRef works metadata for this

	// Abstract
	//FIXME: Need to find value in CrossRef works metadata for this

	// Refereed
	//FIXME: Need to find value in CrossRef works metadata for this

	// Editors (*EditorItemList)
	//FIXME: Need to find value in CrossRef works metadata for this

	// Projects
	//FIXME: Need to find value in CrossRef works metadata for this

	// Contributors (*ContriborItemList)
	//FIXME: Need to find value in CrossRef works metadata for this

	// MonographType
	//FIXME: Need to find value in CrossRef works metadata for this

	// Subjects
	//FIXME: Need to find value in CrossRef works metadata for this

	// PresType (presentation type)
	//FIXME: Need to find value in CrossRef works metadata for this
	return eprint, nil
}

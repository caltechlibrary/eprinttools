package eprinttools

import (
	"fmt"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/crossrefapi"
)

// indexInto takes a map and walks the path and returns a string
// value and succeess bool, if path not found string value is
// an empty string and bool is false.
func indexInto(m map[string]interface{}, parts ...string) (interface{}, bool) {
	switch len(parts) {
	case 0:
		return "", false
	case 1:
		if val, ok := m[parts[0]]; ok == true {
			return val, true
		}
		return "", false
	default:
		if val, ok := m[parts[0]]; ok == true {
			switch val.(type) {
			case map[string]interface{}:
				return indexInto(val.(map[string]interface{}), parts[1:]...)
			}
		}
		return "", false
	}
}

// normalizeType converts content type from CrossRef to Authors (e.g.
// "journal-article" to "article"
func normalizeType(s string) string {
	switch strings.ToLower(s) {
	case "journal-article":
		return "article"
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
		eprint.Type = normalizeType(fmt.Sprintf("%s", s))
	} else {
		return nil, fmt.Errorf("Can't find type in object")
	}
	// Title
	if a, ok := indexInto(obj, "message", "title"); ok == true {
		if len(a.([]interface{})) > 0 {
			eprint.Title = fmt.Sprintf("%s", a.([]interface{})[0].(string))
		}
	}
	// ShortTitle
	// SubTitle

	// NOTE: Assuming IsPublished is true given that we're talking to
	// CrossRef API which holds published content.
	// IsPublished
	eprint.IsPublished = "pub"

	// Publisher
	if s, ok := indexInto(obj, "message", "publisher"); ok == true {
		eprint.Publisher = fmt.Sprintf("%s", s)
	}
	// Publication
	if l, ok := indexInto(obj, "message", "container-title"); ok == true {
		if len(l.([]interface{})) > 0 {
			eprint.Publication = l.([]interface{})[0].(string)
		}
	}
	// Series
	//FIXME: Need to find value in CrossRef works metadata for this

	// Volume
	if s, ok := indexInto(obj, "message", "volume"); ok == true {
		eprint.Volume = fmt.Sprintf("%s", s)
	}
	// Number
	//FIXME: Need to find value in CrossRef works metadata for this

	// Edition
	//FIXME: Need to find value in CrossRef works metadata for this

	// PlaceOfPub
	//FIXME: Need to find value in CrossRef works metadata for this

	// PageRange
	//FIXME: Need to find value in CrossRef works metadata for this

	// Pages
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

	// ISBN
	if s, ok := indexInto(obj, "message", "ISBN"); ok == true {
		eprint.ISSN = fmt.Sprintf("%s", s)
	}

	// ISSN
	if a, ok := indexInto(obj, "message", "ISSN"); ok == true {
		if len(a.([]interface{})) > 0 {
			eprint.ISSN = fmt.Sprintf("%s", a.([]interface{})[0])
		}
	}

	// BookTitle
	//FIXME: Need to find value in CrossRef works metadata for this

	// Editors (*EditorItemList)
	//FIXME: Need to find value in CrossRef works metadata for this

	// Projects
	//FIXME: Need to find value in CrossRef works metadata for this

	// Funders
	//FIXME: Need to find value in CrossRef works metadata for this

	// Contributors (*ContriborItemList)
	//FIXME: Need to find value in CrossRef works metadata for this

	// MonographType???
	//FIXME: Need to find value in CrossRef works metadata for this

	// Subjects???
	//FIXME: Need to find value in CrossRef works metadata for this

	// PresType (presentation type)???
	//FIXME: Need to find value in CrossRef works metadata for this

	// NOTE: Caltech Library puts the DOI in a different field than
	// EPrints' standard DOI location
	// DOI
	if doi, ok := indexInto(obj, "message", "DOI"); ok == true {
		eprint.RelatedURL = new(RelatedURLItemList)
		entry := new(Item)
		entry.Type = "DOI"
		entry.URL = fmt.Sprintf("https://doi.org/%s", doi)
		entry.Description = eprint.Type
		eprint.RelatedURL.AddItem(entry)
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
				entry.Type = "DOI"
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
				entry.URL = fmt.Sprintf("%s", s)
			}
			if s, ok := indexInto(o.(map[string]interface{}), "content-type"); ok == true {
				entry.Type = fmt.Sprintf("%s", s)
			}
			if len(entry.URL) > 0 && len(entry.Type) > 0 {
				eprint.RelatedURL.AddItem(entry)
			}
		}
	}

	// NOTE: Assuming date is published given we're talking to CrossRef
	// Date
	if created, ok := indexInto(obj, "message", "created", "date-time"); ok == true {
		// DateType
		eprint.DateType = "published"
		eprint.Date = fmt.Sprintf("%s", created)
		if len(eprint.Date) > 10 {
			eprint.Date = eprint.Date[0:10]
		}
	}
	// Authors list
	if l, ok := indexInto(obj, "message", "author"); ok == true {
		creators := new(CreatorItemList)
		corpCreators := new(CorpCreatorItemList)
		for _, entry := range l.([]interface{}) {
			author := entry.(map[string]interface{})
			item := new(Item)
			item.Name = new(Name)
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
	return eprint, nil
}

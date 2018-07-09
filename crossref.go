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

// CrossRefWorksToEPrint takes a works object from the CrossRef API
// and maps the fields into an EPrint struct return a new struct or
// error.
func CrossRefWorksToEPrint(obj crossrefapi.Object) (*EPrint, error) {
	eprint := new(EPrint)
	// Type
	if s, ok := indexInto(obj, "message", "type"); ok == true {
		eprint.Type = fmt.Sprintf("%s", s)
	} else {
		return nil, fmt.Errorf("Can't find type in object")
	}
	// Title
	if a, ok := indexInto(obj, "message", "title"); ok == true {
		if len(a.([]interface{})) > 0 {
			eprint.Title = fmt.Sprintf("%s", a.([]interface{})[0].(string))
		}
	}

	// IsPublished
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
	// Volume
	// Number
	// Edition
	// PlaceOfPub
	// PageRange
	// Pages
	// FullTextStatus
	// Keywords
	// Note
	// Abstract
	// Refereed
	// ISBN
	// ISSN
	// BookTitle
	// Editors (*EditorItemList)
	// Projects
	// Funders
	// Contributors (*ContriborItemList)
	// MonographType???
	// Subjects???
	// PresType (presentation type)???

	// DOI (added to the related URL field, not DOI field per normal EPrints)
	if doi, ok := indexInto(obj, "message", "DOI"); ok == true {
		eprint.RelatedURL = new(RelatedURLItemList)
		entry := new(Item)
		entry.Type = "DOI"
		entry.URL = fmt.Sprintf("https://doi.org/%s", doi)
		entry.Description = eprint.Type
		eprint.RelatedURL.AddItem(entry)
	}

	// DateType
	// Date
	if created, ok := indexInto(obj, "message", "created", "date-time"); ok == true {
		eprint.Date = fmt.Sprintf("%s", created)
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

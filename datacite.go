package eprinttools

import (
	"fmt"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataciteapi"
)

// normalizeDataCiteAuthorName takes a name literal and normalizes it to
// an EPrints creator or corp creator item's name attribute.
func normalizeDataCiteAuthorName(obj map[string]interface{}) *Name {
	name := new(Name)
	if s, ok := obj["literal"]; ok == true {
		literal := s.(string)
		if strings.Contains(literal, ",") {
			parts := strings.SplitN(literal, ",", 2)
			name.Family = strings.TrimSpace(parts[0])
			name.Given = strings.TrimSpace(parts[1])
			return name
		}
		if strings.Contains(literal, " ") {
			parts := strings.Split(literal, " ")
			switch {
			case len(parts) == 1:
				name.Value = strings.TrimSpace(literal)
				return name
			case len(parts) == 2:
				name.Family = strings.TrimSpace(parts[1])
				name.Given = strings.TrimSpace(parts[1])
				return name
			case len(parts) > 2:
				last, next_to_last := (len(parts) - 1), (len(parts) - 2)
				switch strings.ToLower(parts[next_to_last]) {
				case "de":
					name.Family = strings.Join(parts[next_to_last:last], " ")
					name.Given = strings.Join(parts[0:next_to_last-1], " ")
				case "van":
					name.Family = strings.Join(parts[next_to_last:last], " ")
					name.Given = strings.Join(parts[0:next_to_last-1], " ")
				default:
					name.Family = strings.TrimSpace(parts[last])
					name.Given = strings.Join(parts[0:next_to_last], " ")
				}
				return name
			}
		}
		// If no spaces or commas then assume a corp style name
		name.Value = literal
		return name
	}
	if given, ok := obj["given"]; ok == true {
		name.Given = given.(string)
	}
	if family, ok := obj["family"]; ok == true {
		name.Family = family.(string)
	}
	return name
}

// normalizeDataCiteToLocalGroup takes an affiliation from DataCite author
// and attempts to determine if it is a Caltech Group and returns
// the normalized name or an empty string.
func normalizeDataCiteToLocalGroup(s string) string {
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

// normalizeDataCiteType converts content type from DataCite to Authors (e.g.
// "journal-article" to "article"
func normalizeDataCiteType(s string) string {
	switch strings.ToLower(s) {
	case "proceedings-article":
		//NOTE: This seems vary idiosyncratic to CaltechAUTHORS
		return "book_section"
	case "journal-article":
		return "article"
	default:
		return s
	}
}

// DataCiteWorksToEPrint takes a works object from the DataCite API
// and maps the fields into an EPrint struct return a new struct or
// error.
func DataCiteWorksToEPrint(obj dataciteapi.Object) (*EPrint, error) {
	if obj == nil {
		return nil, fmt.Errorf("Nothing to convert")
	}
	eprint := new(EPrint)
	// Type
	if s, ok := indexInto(obj, "data", "attributes", "resource-type-id"); ok == true {
		eprint.Type = s.(string)
		switch eprint.Type {
		case "text":
			if s2, ok := indexInto(obj, "data", "attributes", "resource-type-subtype"); ok == true {
				eprint.Type = s2.(string)
			}
		}
	}

	// Title
	if s, ok := indexInto(obj, "data", "attributes", "title"); ok == true {
		eprint.Title = fmt.Sprintf("%s", s)
	}

	// Publisher
	if s, ok := indexInto(obj, "data", "attributes", "publisher"); ok == true {
		eprint.Publisher = fmt.Sprintf("%s", s)
	}

	// Publication
	if s, ok := indexInto(obj, "data", "attributes", "container-title"); ok == true {
		eprint.Publication = s.(string)
	}

	// FIXME: Series

	// FIXME: Volume

	// FIXME: Number

	// FIXME: PlaceOfPub

	// FIXME: PageRange

	// FIXME: Pages

	// FIXME: ISBN

	// FIXME: ISSN

	// NOTE: This doesn't appear to be used by CaltechAUTHORS for full book
	// FIXME: BookTitle

	// FIXME: Funders

	// NOTE: Caltech Library puts the DOI in a different field than
	// EPrints' standard DOI location (i.e. not in eprint.DOI but in
	// the related url item list)
	// DOI
	if doi, ok := indexInto(obj, "data", "attributes", "doi"); ok == true {
		eprint.RelatedURL = new(RelatedURLItemList)
		entry := new(Item)
		entry.Type = "doi"
		entry.URL = fmt.Sprintf("https://doi.org/%s", doi)
		entry.Description = eprint.Type
		eprint.RelatedURL.AddItem(entry)
	}

	// FIXME: RelatedURLs (links in message of DataCite works object)
	// NOTE: related URL type is NOT Mime-Type in CaltechAUTHORS, import URL without type being set.

	// NOTE: Assuming date is published given we're talking to DataCite
	// Date, DateType, IsPublished
	if s, ok := indexInto(obj, "data", "attributes", "published"); ok == true {
		eprint.Date = s.(string)
		eprint.DateType = "published"
		eprint.IsPublished = "pub"
	}

	// Creators/CorpCreators list
	if a, ok := indexInto(obj, "data", "attributes", "author"); ok == true {
		eprint.Creators = new(CreatorItemList)
		eprint.CorpCreators = new(CorpCreatorItemList)
		for _, o := range a.([]interface{}) {
			m := o.(map[string]interface{})
			name := normalizeDataCiteAuthorName(m)
			entry := new(Item)
			entry.Name = name
			if name.Value == "" {
				//NOTE: Assume a person name
				eprint.Creators.AddItem(entry)
			} else {
				//NOTE: Assume a corporate name if we have only a single name
				eprint.CorpCreators.AddItem(entry)
			}
		}
	}

	// Edition
	//FIXME: Need to find value in DataCite works metadata for this

	// FullTextStatus
	//FIXME: Need to find value in DataCite works metadata for this

	// Keywords
	//FIXME: Need to find value in DataCite works metadata for this

	// Note
	//FIXME: Need to find value in DataCite works metadata for this

	// Abstract
	//FIXME: Need to find value in DataCite works metadata for this
	if s, ok := indexInto(obj, "data", "attributes", "description"); ok == true {
		eprint.Abstract = fmt.Sprintf("%s", s)
	}

	// Refereed
	//FIXME: Need to find value in DataCite works metadata for this

	// Editors (*EditorItemList)
	//FIXME: Need to find value in DataCite works metadata for this

	// Projects
	//FIXME: Need to find value in DataCite works metadata for this

	// Contributors (*ContriborItemList)
	//FIXME: Need to find value in DataCite works metadata for this

	// MonographType
	//FIXME: Need to find value in DataCite works metadata for this

	// Subjects
	//FIXME: Need to find value in DataCite works metadata for this

	// PresType (presentation type)
	//FIXME: Need to find value in DataCite works metadata for this

	return eprint, nil
}

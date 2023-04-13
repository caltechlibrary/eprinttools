package eprinttools

import (
	"encoding/json"
	"testing"

	// Caltech Library Packages
	"github.com/caltechlibrary/simplified"
)

// TestSimplifiedRecord checks if the a record rendered in a simplified form
// can be round trip to/from an EPrint type struct
func TestSimplifiedRecord(t *testing.T) {
	// Test data is taken from   https://github.com/inveniosoftware/invenio-rdm-records/blob/396b1f2ff802e8483e30fa2e42cfe03f597d1e87/tests/conftest.py#L361
	example1Text := []byte(`{
        "pids": {
            "doi": {
                "identifier": "10.5281/inveniordm.1234",
                "provider": "datacite",
                "client": "inveniordm"
            }
        },
        "metadata": {
            "resource_type": {"id": "image-photo"},
            "creators": [
                {
                    "person_or_org": {
                        "name": "Nielsen, Lars Holm",
                        "type": "personal",
                        "given_name": "Lars Holm",
                        "family_name": "Nielsen",
                        "identifiers": [
                            {"scheme": "orcid", "identifier": "0000-0001-8135-3489"}
                        ]
                    },
                    "affiliations": [{"id": "cern"}, {"name": "free-text"}]
                }
            ],
            "title": "InvenioRDM",
            "additional_titles": [
                {
                    "title": "a research data management platform",
                    "type": {"id": "subtitle"},
                    "lang": {"id": "eng"}
                }
            ],
            "publisher": "InvenioRDM",
            "publication_date": "2018/2020-09",
            "subjects": [
                {"id": "http://id.nlm.nih.gov/mesh/A-D000007"},
                {"subject": "custom"}
            ],
            "contributors": [
                {
                    "person_or_org": {
                        "name": "Nielsen, Lars Holm",
                        "type": "personal",
                        "given_name": "Lars Holm",
                        "family_name": "Nielsen",
                        "identifiers": [
                            {"scheme": "orcid", "identifier": "0000-0001-8135-3489"}
                        ]
                    },
                    "role": { 
                        "id": "other"
                    },
                    "affiliations": [{"id": "cern"}]
                }
            ],
            "dates": [
                {"date": "1939/1945", "type": {"id": "other"}, "description": "A date"}
            ],
            "languages": [{"id": "dan"}, {"id": "eng"}],
            "identifiers": [{"identifier": "1924MNRAS..84..308E", "scheme": "bibcode"}],
            "related_identifiers": [
                {
                    "identifier": "10.1234/foo.bar",
                    "scheme": "doi",
                    "relation_type": {"id": "iscitedby"},
                    "resource_type": {"id": "dataset"}
                }
            ],
            "sizes": ["11 pages"],
            "formats": ["application/pdf"],
            "version": "v1.0",
            "rights": [
                {
                    "title": {"en": "A custom license"},
                    "description": {"en": "A description"},
                    "link": "https://customlicense.org/licenses/by/4.0/"
                },
                {"id": "cc-by-4.0"}
            ],
            "description": "<h1>A description</h1> <p>with HTML tags</p>",
            "additional_descriptions": [
                {
                    "description": "Bla bla bla",
                    "type": {"id": "methods"},
                    "lang": {"id": "eng"}
                }
            ],
            "locations": {
                "features": [
                    {
                        "geometry": {
                            "type": "Point",
                            "coordinates": [-32.94682, -60.63932]
                        },
                        "place": "test location place",
                        "description": "test location description",
                        "identifiers": [
                            {"identifier": "12345abcde", "scheme": "wikidata"},
                            {"identifier": "12345abcde", "scheme": "geonames"}
                        ]
                    }
                ]
            },
            "funding": [
                {
                    "funder": {
                        "id": "00k4n6c32"
                    },
                    "award": {"id": "00k4n6c32::755021"}
                }
            ],
            "references": [
                {
                    "reference": "Nielsen et al,..",
                    "identifier": "0000 0001 1456 7559",
                    "scheme": "isni"
                }
            ]
        },
        "ext": {
            "dwc": {
                "collectionCode": "abc",
                "collectionCode2": 1.1,
                "collectionCode3": true,
                "test": [ "abc", 1, true ]
            }
        },
        "provenance": {
            "created_by": { "user": "<USER_ID_1>" },
            "on_behalf_of": { "user": "<USER_ID_2>" }
        },
        "access": {
            "record": "public",
            "files": "restricted",
            "embargo": {
                "active": true,
                "until": "2131-01-01",
                "reason": "Only for medical doctors."
            }
        },
        "files": {
            "enabled": true,
            "total_size": 1114324524355,
            "count": 1,
            "bucket": "81983514-22e5-473a-b521-24254bd5e049",
            "default_preview": "big-dataset.zip",
            "order": ["big-dataset.zip"],
            "entries": {
                "big-dataset.zip": {
                    "checksum": "md5:234245234213421342",
                    "mimetype": "application/zip",
                    "size": 1114324524355,
                    "key": "big-dataset.zip",
                    "file_id": "445aaacd-9de1-41ab-af52-25ab6cb93df7"
                }
            },
            "meta": {"big-dataset.zip": {"description": "File containing the data."}}
        },
        "notes": ["Under investigation for copyright infringement."]
    }`)
	simpleRecord := new(simplified.Record)
	if err := json.Unmarshal(example1Text, &simpleRecord); err != nil {
		t.Errorf("Unmarshal failed, %s", err)
		t.FailNow()
	}
	t.Errorf("FIXME: need to test to see if all the elements map across from example.")
}

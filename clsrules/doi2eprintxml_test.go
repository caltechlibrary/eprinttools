package clsrules

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/caltechlibrary/eprinttools"
)

// Custom JSON decoder so we can treat numbers easier
func jsonDecode(src []byte, obj interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	dec.UseNumber()
	err := dec.Decode(&obj)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func TestIssue63(t *testing.T) {
	// DataCite record for DOI 10.1364/prj.437518
	src := []byte(`{
   "data": {
      "attributes": {
         "author": [
            {
               "family": "Khachaturian",
               "given": "Aroutin"
            },
            {
               "family": "Fatemi",
               "given": "Reza"
            },
            {
               "family": "Hajimiri",
               "given": "Ali"
            }
         ],
         "checked": null,
         "citation-count": 0,
         "citations-over-time": [],
         "container-title": "Optica Publishing Group",
         "data-center-id": null,
         "description": null,
         "doi": "10.1364/prj.437518",
         "download-count": 0,
         "downloads-over-time": [],
         "identifier": "https://doi.org/10.1364/prj.437518",
         "license": null,
         "media": [],
         "member-id": null,
         "published": "2022",
         "registered": "2022-05-01T07:26:01.000Z",
         "related-identifiers": [],
         "related-items": [],
         "resource-type-id": "journalarticle",
         "resource-type-subtype": "JournalArticle",
         "results": [],
         "schema-version": "4",
         "title": "Achieving full grating-lobe-free field of view with low-complexity co-prime photonic beamforming transceivers",
         "updated": "2022-05-03T01:33:56.000Z",
         "url": "https://opg.optica.org/abstract.cfm?URI=prj-10-5-A66",
         "version": null,
         "view-count": 0,
         "views-over-time": [],
         "xml": "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPHJlc291cmNlIHhtbG5zOnhzaT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS9YTUxTY2hlbWEtaW5zdGFuY2UiIHhtbG5zPSJodHRwOi8vZGF0YWNpdGUub3JnL3NjaGVtYS9rZXJuZWwtNCIgeHNpOnNjaGVtYUxvY2F0aW9uPSJodHRwOi8vZGF0YWNpdGUub3JnL3NjaGVtYS9rZXJuZWwtNCBodHRwOi8vc2NoZW1hLmRhdGFjaXRlLm9yZy9tZXRhL2tlcm5lbC00L21ldGFkYXRhLnhzZCI+CiAgPGlkZW50aWZpZXIgaWRlbnRpZmllclR5cGU9IkRPSSI+MTAuMTM2NC9wcmouNDM3NTE4PC9pZGVudGlmaWVyPgogIDxjcmVhdG9ycz4KICAgIDxjcmVhdG9yPgogICAgICA8Y3JlYXRvck5hbWUgbmFtZVR5cGU9IlBlcnNvbmFsIj5LaGFjaGF0dXJpYW4sIEFyb3V0aW48L2NyZWF0b3JOYW1lPgogICAgICA8Z2l2ZW5OYW1lPkFyb3V0aW48L2dpdmVuTmFtZT4KICAgICAgPGZhbWlseU5hbWU+S2hhY2hhdHVyaWFuPC9mYW1pbHlOYW1lPgogICAgICA8bmFtZUlkZW50aWZpZXIgbmFtZUlkZW50aWZpZXJTY2hlbWU9Ik9SQ0lEIiBzY2hlbWVVUkk9Imh0dHBzOi8vb3JjaWQub3JnIj5odHRwczovL29yY2lkLm9yZy8wMDAwLTAwMDEtODMwNC0zMzAyPC9uYW1lSWRlbnRpZmllcj4KICAgIDwvY3JlYXRvcj4KICAgIDxjcmVhdG9yPgogICAgICA8Y3JlYXRvck5hbWUgbmFtZVR5cGU9IlBlcnNvbmFsIj5GYXRlbWksIFJlemE8L2NyZWF0b3JOYW1lPgogICAgICA8Z2l2ZW5OYW1lPlJlemE8L2dpdmVuTmFtZT4KICAgICAgPGZhbWlseU5hbWU+RmF0ZW1pPC9mYW1pbHlOYW1lPgogICAgICA8bmFtZUlkZW50aWZpZXIgbmFtZUlkZW50aWZpZXJTY2hlbWU9Ik9SQ0lEIiBzY2hlbWVVUkk9Imh0dHBzOi8vb3JjaWQub3JnIj5odHRwczovL29yY2lkLm9yZy8wMDAwLTAwMDEtOTA4MS0yNjA4PC9uYW1lSWRlbnRpZmllcj4KICAgIDwvY3JlYXRvcj4KICAgIDxjcmVhdG9yPgogICAgICA8Y3JlYXRvck5hbWUgbmFtZVR5cGU9IlBlcnNvbmFsIj5IYWppbWlyaSwgQWxpPC9jcmVhdG9yTmFtZT4KICAgICAgPGdpdmVuTmFtZT5BbGk8L2dpdmVuTmFtZT4KICAgICAgPGZhbWlseU5hbWU+SGFqaW1pcmk8L2ZhbWlseU5hbWU+CiAgICAgIDxuYW1lSWRlbnRpZmllciBuYW1lSWRlbnRpZmllclNjaGVtZT0iT1JDSUQiIHNjaGVtZVVSST0iaHR0cHM6Ly9vcmNpZC5vcmciPmh0dHBzOi8vb3JjaWQub3JnLzAwMDAtMDAwMS02NzM2LTgwMTk8L25hbWVJZGVudGlmaWVyPgogICAgPC9jcmVhdG9yPgogIDwvY3JlYXRvcnM+CiAgPHRpdGxlcz4KICAgIDx0aXRsZT5BY2hpZXZpbmcgZnVsbCBncmF0aW5nLWxvYmUtZnJlZSBmaWVsZCBvZiB2aWV3IHdpdGggbG93LWNvbXBsZXhpdHkgY28tcHJpbWUgcGhvdG9uaWMgYmVhbWZvcm1pbmcgdHJhbnNjZWl2ZXJzPC90aXRsZT4KICA8L3RpdGxlcz4KICA8cHVibGlzaGVyPk9wdGljYSBQdWJsaXNoaW5nIEdyb3VwPC9wdWJsaXNoZXI+CiAgPHB1YmxpY2F0aW9uWWVhcj4yMDIyPC9wdWJsaWNhdGlvblllYXI+CiAgPHJlc291cmNlVHlwZSByZXNvdXJjZVR5cGVHZW5lcmFsPSJKb3VybmFsQXJ0aWNsZSI+Sm91cm5hbEFydGljbGU8L3Jlc291cmNlVHlwZT4KICA8ZGF0ZXM+CiAgICA8ZGF0ZSBkYXRlVHlwZT0iSXNzdWVkIj4yMDIyLTA0LTI5PC9kYXRlPgogICAgPGRhdGUgZGF0ZVR5cGU9IlVwZGF0ZWQiPjIwMjItMDUtMDFUMjM6MjY6MjRaPC9kYXRlPgogIDwvZGF0ZXM+CiAgPHJlbGF0ZWRJZGVudGlmaWVycz4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IklTU04iIHJlbGF0aW9uVHlwZT0iSXNQYXJ0T2YiIHJlc291cmNlVHlwZUdlbmVyYWw9IkNvbGxlY3Rpb24iPjIzMjctOTEyNTwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTM2NC9jbGVvX2F0LjIwMTguanRoNWMuODwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTAzOC9zNDE1NjYtMDE4LTAyNjYtNTwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTM2NC9vZS4yMy4wMDUxMTc8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjEzNjQvY2xlb19hdC4yMDE3Lmp3MmEuOTwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTAzOC9zNDE1OTgtMDIwLTU4MDI3LTE8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjExMDkvanN0cWUuMjAxOS4yOTA4NTU1PC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMzY0L2NsZW9fYXQuMjAxOC5qdGg1Yy4yPC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMDM4L25hdHVyZTExNzI3PC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMDE2L3MwOTI0LTQyNDcoMDEpMDA2MDktNDwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTM2NC9vcHRpY2EuNi4wMDA1NTc8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjExMDkvaXNzY2MuMjAxNy43ODcwMzYxPC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMzY0L29lLjIzLjAyMTAxMjwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTM2NC9vZS4xOS4wMjE1OTU8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjEzNjQvY2xlb19hdC4yMDIwLmp0aDRhLjM8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjExMDkvanNzYy4yMDE5LjI5MzQ2MDE8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjEzNjQvb2UuMjcuMDI3MTgzPC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMzY0L29wdGljYS4zODkwMDY8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjEzNjQvb2wuNDIuMDAwMDIxPC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMzY0L2FvLjQwMzMxNDwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTAwMi85NzgwNDcwNTI5MTg4PC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMTA5L2pzc2MuMjAxOS4yODk2NzY3PC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMzY0L2NsZW9fYXQuMjAxOS5hdzNrLjI8L3JlbGF0ZWRJZGVudGlmaWVyPgogICAgPHJlbGF0ZWRJZGVudGlmaWVyIHJlbGF0ZWRJZGVudGlmaWVyVHlwZT0iRE9JIiByZWxhdGlvblR5cGU9IlJlZmVyZW5jZXMiPjEwLjExMDkvdHNwLjIwMTAuMjA4OTY4MjwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTExNy8xMi4yMDQ0ODIwPC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMzY0L29sLjQyLjAwNDA5MTwvcmVsYXRlZElkZW50aWZpZXI+CiAgICA8cmVsYXRlZElkZW50aWZpZXIgcmVsYXRlZElkZW50aWZpZXJUeXBlPSJET0kiIHJlbGF0aW9uVHlwZT0iUmVmZXJlbmNlcyI+MTAuMTM2NC9wcmouNDM3ODQ2PC9yZWxhdGVkSWRlbnRpZmllcj4KICAgIDxyZWxhdGVkSWRlbnRpZmllciByZWxhdGVkSWRlbnRpZmllclR5cGU9IkRPSSIgcmVsYXRpb25UeXBlPSJSZWZlcmVuY2VzIj4xMC4xMDM4L3M0MTU4Ni0wMjEtMDMyNTkteTwvcmVsYXRlZElkZW50aWZpZXI+CiAgPC9yZWxhdGVkSWRlbnRpZmllcnM+CiAgPHNpemVzLz4KICA8Zm9ybWF0cy8+CiAgPHZlcnNpb24vPgogIDxkZXNjcmlwdGlvbnM+CiAgICA8ZGVzY3JpcHRpb24gZGVzY3JpcHRpb25UeXBlPSJTZXJpZXNJbmZvcm1hdGlvbiI+UGhvdG9uaWNzIFJlc2VhcmNoLCAxMCg1KSwgQTY2PC9kZXNjcmlwdGlvbj4KICA8L2Rlc2NyaXB0aW9ucz4KPC9yZXNvdXJjZT4K"
      },
      "id": "https://doi.org/10.1364/prj.437518",
      "relationships": {
         "data-center": {
            "data": null
         },
         "member": {
            "data": null
         },
         "resource-type": {
            "data": {
               "id": "journal-article",
               "type": "resource-types"
            }
         }
      },
      "type": "works"
   }
}`)
	obj := make(map[string]interface{})
	if err := jsonDecode(src, &obj); err != nil {
		t.Errorf("CrossRef Object unmarshal, %s", err)
		t.FailNow()
	}

	doi := `10.1364/prj.437518`
	eprint, err := eprinttools.DataCiteWorksToEPrint(obj)
	if err != nil {
		t.Errorf(`Expected to crosswalk DataCite record for %q, %s`, doi, err)
		t.FailNow()
	}
	// Now need to apply Caltech Library rules to populate eprint.IDNumber
	eprintsList := new(eprinttools.EPrints)
	eprintsList.Append(eprint)
	ruleSet := UseCLSRules()

	eprintsList, err = Apply(eprintsList, ruleSet)
	if err != nil {
		t.Errorf("expected to apply CLS Rules, %s", err)
	}
	eprint = eprintsList.EPrint[0]
	if eprint.IDNumber == "" {
		t.Errorf("expected IDNUmber, got %q", eprint.IDNumber)
	}
	if eprint.OfficialURL == "" {
		t.Errorf("expected OfficialURL, got %q", eprint.OfficialURL)
	}
	if !strings.HasSuffix(eprint.OfficialURL, eprint.IDNumber) {
		t.Errorf("Expected ID Number path in Official URL to match\n%q\n%q\n", eprint.IDNumber, eprint.OfficialURL)
	}
}

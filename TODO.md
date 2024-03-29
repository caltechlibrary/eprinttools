
Action items
============

Next
----

Transitional support for feeds

- [ ] Simplfied record (invenio-RDM record) to EPrint XML/EPrint JSON
    - [ ] implement libdataset for v2 dataset collection
    - [ ] convert feeds' datasets to v2 ds collections
    - [ ] CaltechDATA records are a good test case for mapping Invenio-RDM to JSON that will drive EPrints templates for feeds (v1 of feeds is EPrints centric)
    - [ ] I need to generate EPrint XML as Invenio-RDM records (simplified records) first and input into the test instance of CaltechAUTHORS as RDM based repository
    - [ ] I can then take the Invenio-RDM records returned and turn those into EPrint XMl/JSON for integration into existing feeds


This is for the simplified eprinttools codebase.

- [x] Make a list of the missing page types I need to render to match feeds
- [ ] Render json types  
	- Misc. JSON files
    	- [ ] caltechauthors-grid.json
    	- [ ] caltechdata-grid.json
    	- [ ] caltechthesis-grid.json
    	- [ ] directory_info.json    (only found for people)
    	- [x] index.json
    	- [ ] pagefind-entry.json
	    - [x] group.json
	    - [x] group_list.json
  		- [ ] people.json
    	- [x] people_list.json
	- CaltechTHESIS derived
    	- [ ] advisor-bachelors.json (needed in groups only)
    	- [ ] advisor-combined.json  (needed in groups only)
    	- [ ] advisor-engd.json      (needed in groups only)
    	- [ ] advisor-masters.json   (needed in groups only)
    	- [ ] advisor-other.json     (needed in groups only)
    	- [ ] advisor-phd.json       (needed in groups only)
    	- [ ] advisor-senior_major.json (needed in groups only)
    	- [ ] advisor-senior_minor.json (needed in groups only)
    	- [x] advisor.json           (person/group)
    	- [ ] bachelors.json         (person/group)
    	- [ ] engd.json              (person/group)
    	- [ ] masters.json           (person/group)
    	- [ ] phd.json               (person/group)
    	- [ ] senior_major.json      (person/group)
    	- [ ] senior_minor.json      (person/group)
	- CaltechAUTHORS derived
   		- [ ] combined.json (from CaltechAUTHORS), should be renamed combined_authors.json
    	- [ ] pub_types.json
   		- [ ] article.json
   		- [ ] audiovisual.json
   		- [ ] book.json
   		- [ ] book_section.json
   		- [ ] collection.json
    	- [ ] conference_item.json
   		- [ ] dataset.json
   		- [ ] image.json
   		- [ ] interactiveresource.json
    	- [ ] model.json
   		- [ ] monograph.json
   		- [ ] object_types.json
   		- [ ] patent.json
   		- [ ] software.json
   		- [ ] teaching_resource.json
   		- [ ] text.json
   		- [ ] thesis.json
   		- [ ] video.json
   		- [ ] workflow.json
	- CaltechDATA
   		- [ ] combined_data.json
    	- [ ] software.json
    	- [ ] data.json
    	- [ ] teaching_resource.json
    	- [ ] data_object_types.json
    	- [ ] data_pub_types.json
    	- [ ] data_types.json
    	- [ ] image.json
- [ ] Render keys types
- [ ] Render Markdown types
- [ ] Render include types
- [ ] Render BibTeX types
- [ ] Render RSS types
- [ ] Make a list of the new pages types reflecting use of simplified records

Bugs
----

- [ ] ep3datasets renders EPrints JSON objects without primary_object being set
- [x] Cleanup eprint content for public views
	- [x] Does sanitization happen at rendering of JSON/Markdown documents or when harvesting the content?
        - Santization happens when we render content, this lets us do one harvest for both dark and public archives
	- [x] Cleanup email content fields
	- [x] Remove "notes"
- [ ] Still debugging mapping the advisor_id, thesis_id and authors_id to person_id for aggregation tables and people feed generation
	- caltechthesis record 15078 is showing up with a local group of "Scott Cushing" who is actaully a committee member not a local gorup.
- [ ] Are messy people identifiers in EPrints are preventing a simple mapping to a single person id, when the EPRint record is read in it needs the ID should be corsswalked to the cl_people_id value.
- [ ] If feeds are "public only" then I need to strip email addresses from the JSON objects.
- [ ] For feeds generated as REPO_NAME-RECORD_TYPE.json to name the feed by record type only, but before I add this I need to see if there is any case where thesis in CaltechAUTHORS need to be itemized along with thesis in CaltechTHESIS
- [x] updated value retrieved from database isn't converting correctly into a time.Time object in Go. Need to figure the best way to make this correct
- [x] Aggregation group_list.json has empty "combined" mapped when there are no eprintid for the specific group in the respository
- [ ] each index.html under people and group should have a corresponding index.json that is used by Pandoc to render index.md that then renders index.html, include.include
- [x] Issue 40, SQL reference document_relation_type table issues
- [ ] Issue 41, Add related URL as DOI value (really make eprints show this as a linked field in the display, don't do that in the data structure)
- [x] Issue 44, Funders are coming up as "UNSPECIFIED"
- [ ] Issue 45, Related URLs are coming in as "UNSPECIFIED"
- [ ] Issue 47, Need to strip HTML from Abstract field
- [ ] Issue 48, Imported EPrint doesn't show up in review buffer
    - [ ] in release 1.1.1-next datestamp isn't set, example eprintid 111912
    - [x] I might be setting the wrong event_status (e.g. buffer or inbox)
    - [x] I need to confirm all timestamp fields and datestamp field is being set correctly
- [x] Issue 49, Field defaults on import including resolver URL and collection
- [ ] Issue 50, Verify why imported and published EPrints don't show in recent additions (is the an issue with generated views or with a datestamp not getting set correctly?).


Next
----

- [ ] Add deposit info to EPrintXML output
- [ ] ioutil is depreciated, need to update the code that uses it
- [x] Need a means of filtering for public EPrint records only
    - `is-public` end point added to ep3apid
    - `?eprint_status=...` added for keys and keys by timestamp ranges
- [x] Add Extended API support to eputil command
- [ ] Implement Solr index record view for Solr 8.9 ingest
- [ ] Add update end point to support update EPrints Metadata
    - [ ] Figure out how historical diffs of EPrints XML are generated in EPrints' History tab
    - [ ] Implement updates versioning the EPrint Metadata record
    - [ ] Implement file upload and manage document versioning

Completed
---------

- [x] Implement an example ep3apid Python API
- [x] Implement a /version end point displaying ep3apid version number
- [x] Create an example service file for running ep3apid as a service under SystemD (Linux)
- [x] Create an example service file for running ep3apid as a service under LaunchD (macOS)
- [x] Need a Users end point to get a list of users in the system and retrieve their numeric user id
- [x] the various related tables that represent item lists don't have the same row count so I need to explicitly query for eprintid, pos or do JOIN and handle the NULL column cases.
- [x] Fix lemurprints-import-api-16 through 21 examples, re-import with ./bin/doi2eprintxml tool
- [x] Add script to generate "lemurprints" database with support for all fields present across our repositories so I can do robust testing and generate appropriate testdata
    - [x] Include all fields and tables in caltechauthors
    - [x] Include all fields and tables in caltechthesis
    - [x] Include all fields and tables in caltechconf
    - [x] Include all fields and tables in caltechcampuspubs
    - [x] Include all fields and tables in calteches
    - [x] Include all fields and tabels in caltechoh
    - [x] Include all fields and tabels in caltechln
    - [x] Exported selected records from production, sanitize them and write import test against lemurprints test database
    - [x] Fetch DOI of records found in EPrints use them to test in lemurprints
- [x] Add create end points to support importing EPrint XML metadata into eprints
    - [x] Implement SQLReadEPrint
    - [x] Implement SQLCreateEPrint
    - [x] Implement ImportEPrint for importing EPrint XML metadata
    - [x] Implement a method that takes a table/column map and EPrint structure then renders a INSERT or REPLACE sequence to create or update an EPrint record
    - [x] Implement a method that takes a table/column map and EPrint structure and update the EPrint structure from a sequnce of SELECT statements
- [x] Split clsrules into separate options to allow for more specific control
- [x] Add end point for `/{REPO_ID}/year` (list years that have eprint records with a "published" date type)
- [x] Add end point for `/{REPO_ID}/year/{YEAR}` lists eprint records published in that year
- [x] Implement a method to show which tables a repository instance has and the column names in each table
    - [x] Implement a startup data structure that captures the `/repository/` end point data so that table/column map can be used to build the SQL queries need to read, create, and update an EPrint record
    - [x] Implement `/repository/<REPO_ID>` end point with `map[string][]string{}` output
- [x] doi2eprintxml list of DOI should allow for pipe separator and URL to object and handle it like Acacia does
- [x] doi2eprintxml needs to fetch the object URL and save results along side the generated EPrints XML
    - added with a -D,-download option in doi2eprintxml.
- [x] Added created (datestamp) end point for feeds
- [x] Implement Simplified JSON record based on
    - https://inveniordm.docs.cern.ch/reference/metadata/
    - https://github.com/caltechlibrary/caltechdata_api/blob/ce16c6856eb7f6424db65c1b06de741bbcaee2c8/tests/conftest.py#L147
- [x] Add simplified JSON output option to
    - [x] eputil
    - [x] epfmt
    - [x] doi2eprintxml

Someday, Maybe
--------------

- [ ] Add end point to recreate Person A-Z list
- [ ] Add end point for subjects
- [ ] Add end point for events (Conferences)
- [ ] Add end point for collection
- [ ] Add end point for publication
- [ ] Add end point for place_of_pub
- [ ] Add end point for issn
- [ ] Add end point for Person (Person IDs)
- [ ] Add end point for Authors (creators)
- [ ] Add end point for Editors
- [ ] Add end point for contributors
- [ ] Add end point for types
- [ ] Add end point for corp_creators
- [ ] Add end point ofr issuing_body



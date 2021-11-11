
Action items
============

This is for the simplified eprinttools codebase.

Bugs
----

+ (branch next) the various related tables that represent item lists don't have the same row count so I need to explicitly query for eprintid, pos or do JOIN and handle the NULL column cases.

Next
----


- [ ] Add script to generate "lemurprints" database with support for all fields present across our repositories so I can do robust testing and generate appropriate testdata
    - [x] Include all fields and tables in caltechauthors
    - [ ] Include all fields and tables in caltechthesis
    - [ ] Include all fields and tables in caltechconf
    - [ ] Include all fields and tables in caltechcampuspubs
    - [ ] Include all fields and tables in calteches
    - [ ] Include all fields and tabels in caltechoh
    - [ ] Include all fields and tabels in caltechln
- [ ] Add create end points to support importing EPrint XML metadata into eprints
    - [x] Creating a new row in eprint table and getting the eprintid assigne needs to be atomic. Wrap in a transaction.
    - [x] Implement a method that takes a table/column map and EPrint structure then renders a INSERT or REPLACE sequence to create or update an EPrint record
    - [x] Implement a method that takes a table/column map and EPrint structure and update the EPrint structure from a sequnce of SELECT statements
- [ ] Add update end point to support update EPrints Metadata
    - [ ] Figure out host historical diffs of EPrints XML are generated in EPrints' History tab
    - [ ] Impement updates versioning the EPrint Metadata record
    - [ ] Implement file upload and manage document versioning
- [ ] Implement Solr index record view for Solr 8.9 ingest
- [ ] Add Person A-Z list support
- [ ] Split clsrules into separate options to allow for more specific control

Completed
---------

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


Completed
---------


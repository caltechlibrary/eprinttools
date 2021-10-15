
Action items
============

This is for the simplified eprinttools codebase.

Bugs
----

Next
----

- [ ] Add create/update/delete of eprints record support to ep3apid, needed to push records generated through Acacia into EPrints
    - [ ] Implement a method that takes a table/column map and EPrint structure then renders a INSERT or REPLACE sequence to create or update an EPrint record
    - [ ] Implement a method that takes a table/column map and EPrint structure and update the EPrint structure from a sequnce of SELECT statements
- [ ] Implement a method to show which tables a repository instance has and the column names in each table
    - [x] Implement `/repository/<REPO_ID>` end point with map[string][]string{} output
    - [ ] Implement a start data structure that captures the `/repository/` end point data so that table/column map can be used to build the SQL queries need to read, create, and update an EPrint record
- [ ] Implement Solr index record view for Solr 8.9 ingest
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

Completed
---------



# Overview

E-Prints produces a list of recently added items to the repository but doesn't provide a feed of items by publication dates.
The ep command line utility provides a way to generate that list by leveraging E-Prints REST API.


## ep functions

+ New() - returns an new API structure
+ ListEPrintsURI() - returns a list of eprint URI from the REST API
+ GetEPrint() - returns an eprint record by URI from the REST API
+ ListURI() - return a list of eprint URI from the database
+ Get() - returns an eprint record by URI from the database
+ ExportEPrints() - exports the contents from the EPrint REST API into a database with a publication dates index
+ GetAllRecordIDs() - return an array of all record ids
+ GetAllRecords() - return an array of all eprint records
+ GetPublishedRecords() - returns an array of records with a date type of "published" from the database
    + IsPublished is "pub"
    + Date Type is "publish"
+ GetPublishedArticles() - recently published articles 
    + Type is "article"
    + IsPublished is "pub"
    + Date Type is "publish"
+ GetLocalGroups() - get a list of unique local groups used in EPrints
+ GetLocalGroupRecords() - get a list of EPrint records for given Local Group
+ GetORCIDS() - get a list of unique orcid ids used in EPrints
+ GetORCIDRecords() - get a list of EPrint records for a given ORCID id
+ BuildPages() - generate a page from an EPrints (BuildSite() makes repeated calls to BuildPage())
+ BuildSite() - generates a JSON file of 25 most recently published articles, a version in RSS and one as an HTML include
+ BuildEPrintMirror() - mirror all the EPrint records in JSON form from BoltDB copy
+ RenderDocuments() - takes a basepath and records array and renders out a directory with *.rss, *.html, *.include and *.json
+ RenderEPrint() - render a mirrored eprint record

Additional functions working off lists

+ ToNames() - from a PersonList type, returns a list of names in Family, Given format
+ ToORCIDs() - from a PersonList type, return a list of ORCID ids found
+ ToISNIs() - from a PersonList type, return a list of ISNI ids found
+ ToAgencies() - from a FunderList type, return a list of agency names
+ ToGrantNumbers() - from a FunderList type, return a list of grant numbers
+ PubDate() - from a Record structure, return a publication date or empty string


## ep Search

A feature that needs to be added to accessing simple and advanced search via the *ep* library.

Two examples of "search" using *curl* for titles starting with "flood characteristic of alluvial"

```shell
    # Basic search URL from using form
    curl "$EPRINTS_URL/cgi/search/simple?screen=Search&order=&q_merge=ALL&q=flood+characteristics+of+alluvial&_action_search=Search"\
    | jq 
    # Advanced Search returned as JSON (title = ??)
    curl "$EPRINTS_URL/cgi/search/archive/advanced/export_caltechauthors_JSON.js?screen=Search&dataset=archive&_action_export=1&output=JSON&exp=0%7C1%7C-date%2Fcreators_name%2Ftitle%7Carchive%7C-%7Ctitle%3Atitle%3AALL%3AIN%3Aflood+characteristics+of+alluvial%7C-%7Ceprint_status%3Aeprint_status%3AANY%3AEQ%3Aarchive%7Cmetadata_visibility%3Ametadata_visibility%3AANY%3AEQ%3Ashow&n="\
    | jq '{title: .[].title, uri: .[].uri}'
```


## environment vars

+ EP_API_URL - base URL for your eprints repository (e.g. http://lemurprints.local)
+ EP_SITE_URL - the website url (might be the same as your eprints repository)
+ EP_DBNAME - the database name used in exporting or building
+ EP_HTDOCS - the htdocs directory where files are written in building
+ EP_TEMPLATES - the template directory holding the templates for building


## htdocs layout

+ publications.*, articles.* come from EPrints
+ data.* comes from the data repository
+ FORMAT: .html, .include (html fragment), .rss, .json (as JSON document), .bib (BibTeX)
+ ORCID-OR-ISNI: if we have an ORCID use that otherwise use ISNI if available (e.g. Richard Feynman only has an ISNI but not an ORCID)
    + ORCID and ISNI do not collide, ORCID is a proper subset of ISNI's four four digit grouping

```
    /index.html <-- docs & description of structure
    /publications.FORMAT
    /articles.FORMAT
    /data.FORMAT
    /recent/index.html <-- docs describing recent 25 and data formats
    /recent/publications.FORMAT
    /recent/articles.FORMAT
    /recent/data.FORMAT
    /person/ORCID-OR-ISNI/publications.FORMAT
    /person/ORCID-OR-ISNI/articles.FORMAT
    /person/ORCID-OR-ISNI/data.FORMAT
    /person/ORCID-OR-ISNI/recent/publications.FORMAT
    /person/ORCID-OR-ISNI/recent/articles.FORMAT
    /person/ORCID-OR-ISNI/recent/data.FORMAT
    ...
    /affiliations/... <-- needs to be defined around a uniqu ID (e.g. what was talked about at pidpaloosa)
    /search/?q=... <-- search engine results (also available in various FORMAT)
    /tools/... <-- any web based tools we invent to work with the feeds
```


## Reference materials

+ [Web API Docs](http://wiki.eprints.org/w/API:EPrints/Apache/CRUD)
    + includes example curl interactions
+ [Golang Template Tutorial](https://elithrar.github.io/article/approximating-html-template-inheritance/)
    + page.html and page.include are the two templates needed to produce the publications feed pages




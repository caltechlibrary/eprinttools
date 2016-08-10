
# Overview

E-Prints produces a list of recently added items to the repository but doesn't provide a feed of items by publication dates.
The epgo command line utility provides a way to generate that list by leveraging E-Prints REST API.

## epgo functions

+ New() - returns an new API structure
+ ListEPrintsURI() - returns a list of eprint URI from the REST API
+ GetEPrint() - returns an eprint record by URI from the REST API
+ ListURI() - return a list of eprint URI from the database
+ Get() - returns an eprint record by URI from the database
+ Export() - exports the contents from the EPrint REST API into a database with a publication dates index
+ GetPublishedRecords() - returns an array of records with a date type of "published" from the database
    + IsPublished is "pub"
    + Date Type is "publish"
+ GetPublishedArticles() - recently published articles 
    + Type is "article"
    + IsPublished is "pub"
    + Date Type is "publish"
+ BuildSite() - generates a JSON file of 25 most recently published articles, a version in RSS and one as an HTML include
+ RenderDocuments() - takes a basepath and records array and renders out a directory with rss.xml, index.html, index.include and index.json

## epgo Search

A feature that needs to be added to accessing simple and advanced search via the *epgo* library.

Two examples of "search" using *curl* for titles starting with "flood characteristic of alluvial"

```path
    # Basic search URL from using form
    curl "$EPRINTS_URL/cgi/search/simple?screen=Search&order=&q_merge=ALL&q=flood+characteristics+of+alluvial&_action_search=Search" | \
        jq 
    # Advanced Search returned as JSON (title = ??)
    curl $"EPRINTS_URL/cgi/search/archive/advanced/export_caltechauthors_JSON.js?screen=Search&dataset=archive&_action_export=1&output=JSON&exp=0%7C1%7C-date%2Fcreators_name%2Ftitle%7Carchive%7C-%7Ctitle%3Atitle%3AALL%3AIN%3Aflood+characteristics+of+alluvial%7C-%7Ceprint_status%3Aeprint_status%3AANY%3AEQ%3Aarchive%7Cmetadata_visibility%3Ametadata_visibility%3AANY%3AEQ%3Ashow&n=" | jq '{title: .[].title, uri: .[].uri}'
```


## environment vars

+ EPGO_API_URL - base URL for your eprints repository (e.g. http://lemurprints.local)
+ EPGO_SITE_URL - the website url (might be the same as your eprints repository)
+ EPGO_DBNAME - the database name used in exporting or building
+ EPGO_HTDOCS - the htdocs directory where files are written in building
+ EPGO_TEMPLATES - the template directory holding the templates for building

## Reference materials

+ [Web API Docs](http://wiki.eprints.org/w/API:EPrints/Apache/CRUD)
    + includes example curl interactions
+ [Golang Template Tutorial](https://elithrar.github.io/article/approximating-html-template-inheritance/)
    + page.html and page.include are the two templates needed to produce the publications feed pages

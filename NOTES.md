
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
+ BuildSite() - generates a JSON file of 25 most recently published articles, a version in RSS and one as an HTML include

## environment vars

+ EPGO_BASE_URL - base URL for your eprints repository (e.g. http://lemurprints.local)
+ EPGO_DBNAME - the database name used in exporting or building
+ EPGO_HTDOCS - the htdocs directory where files are written in building
+ EPGO_TEMPLATES - the template directory holding the templates for building

## Reference materials

+ [Web API Docs](http://wiki.eprints.org/w/API:EPrints/Apache/CRUD)
    + includes example curl interactions
+ [Golang Template Tutorial](https://elithrar.github.io/article/approximating-html-template-inheritance/)
    + page.html and page.include are the two templates needed to produce the publications feed pages

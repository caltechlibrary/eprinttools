
# Overview

## epgo functions

+ New() - returns an new API structu
+ ListEPrintsURI() - returns a list of eprint URI
+ GetEPrint() - returns an eprint record from a URI

## environment vars

+ EPGO_BASE_URL - base URL for your eprints repository (e.g. http://lemurprints.local)
+ EPGO_DBNAME - the database name used in exporting or building
+ EPGO_HTDOCS - the htdocs directory where files are written in building
+ EPGO_TEMPLATES - the template directory holding the templates for building

## Reference

+ [Web API Docs](http://wiki.eprints.org/w/API:EPrints/Apache/CRUD)
    + includes example curl interactions
+ [Golang Template Tutorial](https://elithrar.github.io/article/approximating-html-template-inheritance/)
    + page.html and page.include are the two templates needed to produce the publications feed pages

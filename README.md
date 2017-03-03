[![Go Report Card](http://goreportcard.com/badge/caltechlibrary/epgo)](http://goreportcard.com/report/caltechlibrary/epgo)

# epgo

This project contains the _epgo_ go package for working with EPrints REST API. It also includes
a set of command line utilities for harvesting content,  building feeds and website pages
that run indepenantly of EPrints itself (e.g. You could have a website/feeds generated from 
an EPrints repository running on a different system).

## The command line utilities

+ _epgo_ is a command line utility utilizing EPrints' REST API to produce alternative
feeds and formats. Currently it supports generating a feed of repository items based
on publication dates.
+ _epgo-genpages_ is a command line utility that builds HTML and feed pages from content harvested with _epgo_
    + NOTE: this will change, epgo-genpages should probably be epgo-gendocs and render JSON, BibTeX and RSS2 only
+ _epgo-indexpages_ is a command line utlity that will build a [Bleve](https://blevesearch.com) index to support website search
    + NOTE: could this be depreciated in favor of a generalized indexer?
+ _epgo-servepages_ is a web server for serving the static content generated with _genpages_ as well as supporting search from the index created with _epgo-indexpages_
    + NOTE: this should be depreciated in favor of _ws_ from the _mkpage_ project
+ _epgo-sitemapper_ is a command line utility for generating a sitemap.xml file
    + NOTE: this should be depreciated in favor of _sitemapper_ from the _mkpage_ project

All the utilities can be configured from the environment.  The environment
can be overridden by command line options. For details run the individual command wiht the '-help'
option.  E.g. `./bin/epgo -help`





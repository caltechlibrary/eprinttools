
# USAGE

    epgo-servepages [OPTIONS]

## OVERVIEW

epgo-servepages a webserver for explosing EPrints as HTML pages,  HTML .include pages, JSON and BibTeX formats.

## CONFIGURATION

epgo-servepages can be configured through setting the following environment
variables-

+ EPGO_BLEVE     a colon delimited list of Bleve indexes

+ EPGO_HTDOCS    this is the directory where the HTML files are written.

+ EPGO_SITE_URL  this is the website URL that the public will use

+ EPGO_TEMPLATE_PATH this is the directory that contains the templates
                 used used to generate the static content of the website.

## OPTIONS

```
	-bleve	a colon delimited list of Bleve index db names
	-enable-search	turn on search support in webserver
	-h	display help
	-help	display help
	-htdocs	specify where to write the HTML files to
	-l	display license
	-license	display license
	-site-url	the website url
	-template-path	specify where to read the templates from
	-v	display version
	-version	display version
``` 



# USAGE

    epgo [OPTIONS] [EPRINT_URI]

## SYSNOPSIS

epgo wraps the REST API for E-Prints 3.3 or better. It can return a list of uri,
a JSON view of the XML presentation as well as generates feeds and web pages.

## CONFIG

epgo can be configured with following environment variables

+ EPGO_API_URL (required) the URL to your E-Prints installation
+ EPGO_DATASET   (required) the dataset and collection name for exporting, site building, and content retrieval
+ EPGO_BLEVE (optional) the name for the Bleve index/db
+ EPGO_SITE_URL (optional) the URL to your public website (might be the same as E-Prints)
+ EPGO_HTDOCS   (optional) the htdocs root for site building
+ EPGO_TEMPLATE_PATH (optional) the template directory to use for site building

If EPRINT_URI is provided then an individual EPrint is return as
a JSON structure (e.g. /rest/eprint/34.xml). Otherwise a list of EPrint paths are
returned.

## OPTIONS

```
	-api	url for EPrints API
	-articles-newest	list the N newest published articles
	-dataset	dataset/collection name
	-export	export N EPrints from highest ID to lowest
	-feed-size	number of items rendering in feeds
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-p	pretty print JSON output
	-published-newest	list the N newest published items
	-read-api	read the contents from the API without saving in the database
	-s	generate select lists in dataset
	-select	generate select lists in dataset
	-v	display version
	-version	display version
```


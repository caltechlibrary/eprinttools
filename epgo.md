
# USAGE 

    epgo [OPTIONS] [EPGO_EPRINTS_URL]

## SYNOPSIS

epgo wraps the REST API for E-Prints 3.3 or better. It can return a list 
of uri, a JSON view of the XML presentation as well as generates feeds 
and web pages.

## CONFIGURATION

epgo can be configured with following environment variables

+ EPGO_EPRINTS_URL the URL to your E-Prints installation
+ EPGO_DATASET the dataset and collection name for exporting, site building, and content retrieval

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
	-o	output filename (logging)
	-output	output filename (logging)
	-p	pretty print JSON output
	-published-newest	list the N newest published items
	-read-api	read the contents from the API without saving in the database
	-s	generate select lists in dataset
	-select	generate select lists in dataset
	-v	display version
	-version	display version
```

## EXAMPLE

```shell
    epgo -export all
```

Would export the entire EPrints repository public content defined by the
environment virables EPGO_API_URL, EPGO_DATASET.

```shell
    epgo -export 2000
```

Would export 2000 EPrints from the repository with the heighest ID values.

```shell
    epgo -select
```

Would (re)build the select lists based on contents of $EPGO_DATASET.

```shell
    epgo -select -export all
```

Would export all eprints and rebuild the select lists.


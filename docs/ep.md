
# ep

## USAGE 

    ep [OPTIONS] [EP_EPRINTS_URL]

## SYNOPSIS

ep wraps the REST API for E-Prints 3.3 or better. It can return a list 
of uri, a JSON view of the XML presentation as well as generates feeds 
and web pages.

## CONFIGURATION

ep can be configured with following environment variables

+ EP_EPRINTS_URL the URL to your E-Prints installation
+ EP_DATASET the dataset and collection name for exporting, site building, and content retrieval

## OPTIONS

```
	-api	url for EPrints API
	-articles-newest	list the N newest published articles
	-auth	set the authentication method (e.g. none, basic, oauth, shib)
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
	-pretty	pretty print JSON output
	-published-newest	list the N newest published items
	-pw	set the password
	-read-api	read the contents from the API without saving in the database
	-s	generate select lists in dataset
	-select	generate select lists in dataset
	-un	set the username
	-username	set the username
	-v	display version
	-version	display version
```


## EXAMPLE

```shell
    ep -export all
```

Would export the entire EPrints repository public content defined by the
environment virables EP_API_URL, EP_DATASET.

```shell
    ep -export 2000
```

Would export 2000 EPrints from the repository with the heighest ID values.

```shell
    ep -select
```

Would (re)build the select lists based on contents of $EP_DATASET.

```shell
    ep -select -export all
```

Would export all eprints and rebuild the select lists.

ep v0.0.10-beta1

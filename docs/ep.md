
# ep

## USAGE 

    ep [OPTIONS] [EP_EPRINTS_URL|ONE_OR_MORE_EPRINT_ID]

## SYNOPSIS

ep wraps the REST API for EPrints 3.3 or better. It can return a list 
of uri, a JSON view of the XML presentation as well as generates feeds 
and web pages.

## CONFIGURATION

ep can be configured with following environment variables

+ EP_EPRINTS_URL the URL to your EPrints installation
+ EP_DATASET the dataset and collection name for exporting, site building, and content retrieval

## OPTIONS

```
	-api	url for EPrints API
	-auth	set the authentication method (e.g. none, basic, oauth, shib)
	-dataset	dataset/collection name
	-export	export N EPrints from highest ID to lowest
	-export-since	export  EPrints from a given date to present (e.g. 2017-07-01)
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-o	output filename (logging)
	-output	output filename (logging)
	-p	pretty print JSON output
	-pretty	pretty print JSON output
	-pw	set the password
	-read-api	read the contents from the API without saving in the database
	-un	set the username
	-updated-since	list EPrint IDs updated since a given date (e.g 2017-07-01)
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


ep v0.0.10-beta2

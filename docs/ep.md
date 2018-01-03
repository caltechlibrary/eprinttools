DEBUG showHelp false

# USAGE

	ep [OPTIONS] [EP_EPRINTS_URL|ONE_OR_MORE_EPRINT_ID]

## SYNOPSIS


ep wraps the REST API for EPrints 3.3 or better. It can return a list 
of uri, a JSON view of the XML presentation as well as generates feeds 
and web pages.

CONFIGURATION

ep can be configured with following environment variables

EP_EPRINTS_URL the URL to your EPrints installation

EP_DATASET the dataset and collection name for exporting, site building, and content retrieval

EP_SUPPRESS_NOTE if set to true or 1 will suppress the note field in harvesting


## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    EP_DATASET         Sets the dataset collection for storing EPrint harvested records
    EP_EPRINT_URL      Sets the EPRints API URL
    EP_SUPPRESS_NOTE   Suppress the note field on harvesting


## OPTIONS

Options will override any corresponding environment settings.

```
    -api                      url for EPrints API
    -auth                     set the authentication method (e.g. none, basic, oauth, shib)
    -dataset                  dataset/collection name
    -examples                 display example(s)
    -export                   export N EPrints from highest ID to lowest
    -export-modified          export records by date or date range (e.g. 2017-07-01)
    -export-save-keys         save the keys exported in a file with provided filename
    -generate-markdown-docs   generation markdown documentation
    -h, -help                 display help
    -l, -license              display license
    -o, -output               output filename
    -p, -pretty               pretty print JSON output
    -pw, -password            set the password
    -quiet                    suppress error output
    -read-api                 read the contents from the API without saving in the database
    -suppress-note            suppress note
    -un, -username            set the username
    -updated-since            list EPrint IDs updated since a given date (e.g 2017-07-01)
    -v, -version              display version
    -verbose                  verbose logging
```


## EXAMPLES


Would export the entire EPrints repository public content defined by the
environment virables EP_API_URL, EP_DATASET.

    ep -export all

Would export 2000 EPrints from the repository with the heighest ID values.

    ep -export 2000

Would export the EPrint records modified since July 1, 2017.

    ep -export-modified 2017-07-01

Would export the EPrint records with modified times in July 2017 and
save the keys for the records exported with one key per line. 

    ep -export-modified 2017-07-01,2017-07-31 \
       -export-save-keys=july-keys.txt 


ep v0.0.10-beta7


# USAGE

	ep [OPTIONS] [EPRINT_URL] [ONE_OR_MORE_EPRINT_ID]

## DESCRIPTION


ep uses the REST API for EPrints 3.x to harvest EPrints content into
a dataset collection. If you don't need dataset integration use eputil 
instead. If you want to view  the harvested content then use the
dataset command.

CONFIGURATION

ep can be configured with following environment variables

EPRINT_URL the URL to your EPrints installation

DATASET the dataset collection name to use for storing your harvested EPrint content.


## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    DATASET      # Sets the dataset collection for storing EPrint harvested records
    EPRINT_URL   # Sets the EPRints API URL
```

## OPTIONS

Below are a set of options available. Options will override any corresponding environment settings.

```
    -api                    url for EPrints API
    -auth                   set the authentication method (e.g. none, basic, oauth, shib)
    -dataset                dataset collection name
    -examples               display example(s)
    -export                 export N EPrints from highest ID to lowest
    -export-keys            export using a delimited list of EPrint keys
    -export-modified        export records by date or date range (e.g. 2017-07-01)
    -export-save-keys       save the keys exported in a file with provided filename
    -export-with-docs       include EPrint documents with export
    -generate-manpage       generation man page
    -generate-markdown      generation markdown documentation
    -h, -help               display help
    -l, -license            display license
    -nl, -newline           set to false to exclude trailing newline
    -o, -output             output filename
    -p, -pretty             pretty print JSON output
    -pw, -password          set the password
    -quiet                  suppress error output
    -read-api               read the contents from the API without saving in the database
    -suppress-suggestions   suppress the suggestions field from output
    -un, -username          set the username
    -updated-since          list EPrint IDs updated since a given date (e.g 2017-07-01)
    -v, -version            display version
    -verbose                verbose logging
```


## EXAMPLES


Save a list the URI end points for eprint records found at EPRINT_URL.

	ep -o uris.txt

Export the entire EPrints repository public content defined by the
environment variables EPRINT_URL, DATASET.

    ep -export all

Export 2000 EPrints from the repository with the heighest ID values.

    ep -export 2000

Export the EPrint records modified since July 1, 2017.

    ep -export-modified 2017-07-01

Explore a specific list of keys (e.g. "101,102,1304" or
if list is '-' then read from standard input, one key per line)

	ep -export-keys "101,102,1304"

	ep -export-keys "-"

Export the EPrint records with modified times in July 2017 and
save the keys for the records exported with one key per line. 

    ep -export-modified 2017-07-01,2017-07-31 \
       -export-save-keys=july-keys.txt 


ep v0.0.58


# USAGE

	eputil [OPTIONS]

## SYNOPSIS


	eputil parses XML content retrieved from disc or the EPrints API. It will 
	render JSON if the XML is valid otherwise return errors.


## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    EP_USER   # Sets the eprint USERNAME:PASSWORD to use with URL basic auth
```

## OPTIONS

Options will override any corresponding environment settings.

```
    -auth                     set the authorization type, e.g. basic
    -document, -eprints       parse an eprints (e.g. rest response) document
    -e, -examples             display examples
    -generate-markdown-docs   output documentation in Markdown
    -get, -url                do an HTTP GET to fetch the XML from the URL then parse
    -h, -help                 display help
    -i, -input                input file name
    -ids                      get a list of doc paths (e.g. ids or sub-fields depending on the URL provided
    -json                     attempt to parse XML into generaic JSON structure
    -l, -license              display license
    -nl, -newline             if true add a trailing newline
    -o, -output               output file name
    -post                     do an HTTP POST to add/update a record (e.g. update a URL endpoint
    -put                      do an HTTP POST to add/update a record (e.g. update a URL endpoint
    -pw, -password            set the password for authenticated access
    -quiet                    suppress error messages
    -revision, -eprint        parse a eprint (revision) document
    -u, -user                 set the basic auth string to 'username:password' for authenticated access
    -un, -username            set the username for authenticated access
    -v, -version              display version
```


## EXAMPLES


Fetch an EPrints document as JSON from a URL for an EPrint with an id of 123

    eputil -get https://eprints.example.org/rest/eprint/123.xml -json

Fetch an EPrints document as XML from a URL for an EPrint with an id of 123

    eputil -get https://eprints.example.org/rest/eprint/123.xml

Fetch the creators.xml as JSON for an EPrint with the id of 123.

    eputil -get https://eprints.example.org/rest/eprint/123/creators.xml -json

Parse an EPrint reversion XML document

    eputil -i revision/2.xml -eprint

Get a JSON array of eprint ids from the REST API

    eputil -get https://eprints.example.org/rest/eprint/ -ids


eputil v0.0.13-dev

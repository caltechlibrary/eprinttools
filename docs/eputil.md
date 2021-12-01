

USAGE
=====

```
	eputil [OPTIONS] EPRINT_REST_URL
```
	
SYNOPSIS
--------

eputil is a command line program for exploring 
the EPrints Extended API (provided by ep3apid) or EPrint's
own REST API.  Records are returned in either JSON or EPrints XML.
Lists of eprint ids are returned in JSON.

DETAIL FOR EPrints Extended API
-------------------------------

The extended API is expected to be present on the local machine
at http://localhost:8484.  eputil will convert the command line
parameters into the appropriate URL encoding the command line as
necessary and return the values from the Extended API end points.

The format of the command working with the EPrints extended API
is `eputil REPO_ID END_POINT_NAME [PARAMETER ...]`
You must specify the repository id in the command. E.g.

```
    eputil caltechauthors keys
	eputil caltechauthors doi
	eputil caltechauthors doi "10.5062/F4NP22DV"
	eputil caltechauthors creator-name "Morrell" "Thomas"
	eputil caltechauthors grant-number 
	eputil caltechauthors grant-number "kzcx3-sw-147"
```

See website for a full list of available end points.
[ep3apid](ep3apid.html)

DETAIL FOR EPrints REST API
---------------------------

eputil parses XML content retrieved from 
EPrints 3.x. REST API. It will render 
results in JSON or XML.  With the `-raw`
option you can get the unmodified EPrintXML from the 
REST API otherwise the XML is parsed before final 
rendering as JSON or XML. It requires a basic knowledge
of the layout of EPrint 3.x's REST API. It supports
both unauthenticated and basic authentication access
to the API. The REST API authentication mechanism 
appears indepent of the primary website authentication
setup of the installed EPrints (at least at Caltech
Library). See the examples to start exploring the API.

```
  -auth string
    	set the authentication type for access, default is basic (default "basic")
  -document
    	Retrieve a document from the provided url
  -h	display help
  -help
    	display help
  -i string
    	input file name (read the URL connection string from the input file
  -input string
    	input file name (read the URL connection string from the input file
  -json
    	attempt to parse XML into generaic JSON structure
  -license
    	display license
  -newline
    	if true add a trailing newline
  -nl
    	if true add a trailing newline
  -o string
    	output file name
  -output string
    	output file name
  -password
    	Prompt for the password for authenticated access
  -quiet
    	suppress error messages
  -raw
    	get the raw EPrint REST API response
  -s	Return the object in a simplified JSON data structure.
  -simple
    	Return the object in a simplified JSON data structure.
  -u string
    	set the username for authenticated access
  -un string
    	set the username for authenticated access
  -user string
    	set the username for authenticated access
  -username string
    	set the username for authenticated access
  -verbose
    	verbose output
  -version
    	display version
```

Fetch the raw unmarshaled EPrint XML via the 
EPrint REST API for id 123.

```
    eputil -raw https://example.org/rest/eprint/123.xml
```

Fetch the EPrint XML marshaled as XML using the 
EPrints REST API for id 123.

```
    eputil https://example.org/rest/eprint/123.xml 
```

Fetch the EPrint XML marshaled as JSON using the
EPrints REST API for id 123.

```
    eputil -json https://example.org/rest/eprint/123.xml
```

Get a JSON array of eprint ids from the REST API

```
    eputil -json https://example.org/rest/eprint/ 
```

Get the last modified date for id 123 from REST API

```
    eputil -raw https://example.org/rest/eprint/123/lastmod.txt 
```

If the EPrint REST API is protected by basic authentication
you can pass the username and password via command line
options. You will be prompted for the password value.
or via the URL.  In this example the username is 
"user" and password is "secret". In this example you will
be prompted to enter your secret.

```
    eputil -username=user -password \
      https://example.org/rest/eprint/123.xml
```

You can also pass the username and secret via the URL
but this leaves you vunerable to the password being recorded
in your command history or if another person has access to
the process table. You SHOULD NOT use this approach on a
shared machine!

```
    eputil https://user:secret@example.org/rest/eprint/123.xml
```

Getting IDs doesn't typically require authentication but seeing
specific records may depending on the roles and security
setup implemented in the EPrint instance.




# USAGE

    epgo-genpages [OPTIONS]

## SYNOPSIS

epgo-genpages generates JSON documents in a htdocs directory tree.

## CONFIGURATION

epgo-genpages can be configured through setting the following environment
variables-

+ EPGO_DATASET this is the dataset and collection directory (e.g. dataset/eprints)
+ EPGO_HTDOCS this is the directory where the JSON documents will be written.

## OPTIONS

```
	-api-url	the EPrints API url
	-build-eprint-mirror	Build a mirror of EPrint content rendered as JSON documents
	-dataset	the dataset/collection name
	-h	display help
	-help	display help
	-htdocs	specify where to write the HTML files to
	-l	display license
	-license	display license
	-o	output filename (log message)
	-output	output filename (log message)
	-repository-path	specify the repository path to use for generated content
	-site-url	the website url
	-template-path	specify where to read the templates from
	-v	display version
	-version	display version
```

## EXAMPLE

```shell
	epgo-genpages 
```

Generates JSON documents in EPGO_HTDOCS from EPGO_DATASET.


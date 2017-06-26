
# ep-genfeeds

## USAGE

    ep-genfeeds [OPTIONS]

## SYNOPSIS

ep-genfeeds generates JSON documents (feeds) in a htdocs directory tree.

## CONFIGURATION

ep-genfeeds can be configured through setting the following environment
variables-

+ EP_DATASET this is the dataset and collection directory (e.g. dataset/eprints) 
+ EP_HTDOCS  this is the directory where the JSON documents will be written.

## OPTIONS

```
	-d	specify where to write the HTML files to
	-dataset	the dataset/collection name
	-docs	specify where to write the HTML files to
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-o	output filename (log message)
	-output	output filename (log message)
	-v	display version
	-version	display version
```


EXAMPLE

```shell
    ep-genfeeds 
```

Generates JSON documents in EP_HTDOCS from EP_DATASET.

ep-genfeeds v0.0.10-beta1

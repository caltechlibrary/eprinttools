eprinttools

eprinttools is a collection of command line tools written in Go, a Go
package and set of command line utilities for working with EPrints 3.x
EPrint XML and REST API written in Python 3. Eventually it is planned to
have this project become a pure Python project.

This project also hosts demonstration code to replicate a public facing
version of an EPrints repository outside of EPrints. Think of it as the
public views and landing pages.

Go base code

The command line programs

-   eputil is a command line utility for interacting (e.g. harvesting)
    JSON and XML from EPrints’ REST API
    -   uses minimal configuration because it does less!
    -   it superceded the ep command
-   epfmt is a command line utility to pretty print EPrints XML and
    convert to/from JSON
    -   in the process of pretty printing it also validates the EPrints
        XML against the eprinttools Go package definitions
-   doi2eprintxml is a command line program for turning metadata
    harvested from CrossRef and DataCite into an EPrint XML document
    based on one or more supplied DOI
-   eprintxml2json is a command line program for taking EPrint XML and
    turning it into JSON

The first two utilities can be configured from the environment or
command line options. The environment settings are overridden by command
line options. For details running either command envoke the tool name
with the ‘-help’ option.

Python base code

Python Modules

eprints3x

This python module wraps the eputil Go comand in Python. It makes it
trivial to implement harvesting an EPrints repository using the EPrints
REST API.

eprintviews

This python module uses py_dataset and the harvested content to generate
a htdocs directory similar to the URL layout of EPrints. It features
classes for working with Views, Users (needed to attribute names to
userid fields EPrint XML harvested from the REST API), Subjects (a way
to load the subjects text file from an EPrints archive and generate the
path to label mapping used when rendering views into an htdocs
directory) and Aggregator (this does the heavy lifting of processing a
dataset collection of harvested EPrint XML and generating the views as
JSON documents in the htdocs directory).

command line tools

harvester_full.py, harvester_recent.py

These two Python programs use eprints3x module to implement harvesters
of EPrint XML and any related digitl objects (e.g. PDFs, images) into a
dataset collection

genviews.py

This Python program processes a dataset collection and renders an htdocs
tree populating it with JSON documents and key lists. This skeleton of
metadata and directory structure can then be processed into a rendered
website mirroring the content from an EPrints repository. This module
relies on eprintviews.

indexer.py

This Python program indexes the contents of our replicated EPrints site
by creating scheme.json files along side the index.json files that
represent the landing pages for the replicated repository. These can
then be easily ingested into Lunr.js or Elasticsearch. Currently the
proof of concept targets Lunr.js. This module relies on eprintviews.

mk_website.py

This Python program creates the HTML pages from Markdown documents in
the static folder (e.g. home page and major landing pages) as well as
the individual views and abstracts from the JSON documents created by
genviews.py. The final result is a static website ready to serve out to
the public. This module relies on eprintviews.

publisher.py

This Python program copys (syncs) the content with an AWS S3 bucket via
the AWS command line tools.

Related GitHub projects

py_dataset

This Python module provides access to dataset collections which we use
as intermediate storage for JSON documents and related attachments.

AMES

The eprintools command line programs have been made available to Python
via the AMES project. This include support for both read and write to
EPrints repository systems.

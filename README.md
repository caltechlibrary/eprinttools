
# eprinttools

eprinttools is a Go package and set of command line utilities for working 
with EPrints 3.x EPrint XML and REST API. It is in the process of transitioning to a Python 3 based project.

## The command line programs

+ [eputil](docs/eputil.html) is a command line utility for interacting (e.g. harvesting) JSON and XML from EPrints' REST API
    + uses minimal configuration because it does less!
    + it superceded the _ep_ command
+ [epfmt](docs/epfmt.html) is a command line utility to pretty print EPrints XML and convert to/from JSON
    + in the process of pretty printing it also validates the EPrints XML against the eprinttools Go package definitions
+ [doi2eprintxml](docs/doi2eprintxml.html) is a command line program for turning metadata harvested from CrossRef and DataCite into an EPrint XML document based on one or more supplied DOI
+ [eprintxml2json](docs/eprintxml2json.html) is a command line program for taking EPrint XML and turning it into JSON 

The first two utilities can be configured from the environment or 
command line options. The environment settings are overridden by command 
line options. For details running either command envoke the
tool name with the '-help' option. 

## Python integration via AMES

The eprintools command line programs have been made available to Python
via the [AMES](https://github.com/caltechlibrary/ames) project. This include support for both read and write to EPrints repository systems.

## Demo of wrapping _eputil_ in Python

There is a proof of concept wrapper for for _eputil_ in Python 3
in the "eprints3x" folder.  There is an example "demo-harvester-full.py" 
and "demo-harvester-recent.py" that shows how you can easily implement 
a harvester for your EPrints 3.3x content.


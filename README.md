eprinttools
===========

This is a collection of command line tools written in Go, a Go for working with EPrints 3.3.x EPrint XML and with the EPrint REST API.

This project also hosts demonstration code to replicate a public facing version of an EPrints repository outside of EPrints. Think of it as the public views and landing pages.

Go base code
------------

The command line programs.

- [eputil](docs/eputil.html) is a command line utility for interacting (e.g. harvesting) JSON and XML from EPrints' REST API
    - minimal configuration (because it does so much less!)
- [epfmt](docs/epfmt.html) is a command line utility to pretty print EPrints XML and convert to/from JSON including a simplified JSON inspired by DataCite and Invenion 3
- [doi2eprintxml](docs/doi2eprintxml.html) is a command line program for turning metadata harvested from CrossRef and DataCite into an EPrint XML document based on one or more supplied DOI
- [ep3apid](docs/ep3apid.html) is a Unix style web service for interacting with an EPrint repository via a localhost proxy. It includes the ability to get restricted key lists as well as retrieve a simplified JSON record representing an EPrints record

Configuration is via a JSON "settings.json" file. The settings includes a repository id with "dsn" (Data Source Name) attribute for accessing EPrint's MySQL database(s) and "rest" attribute holding the base URL used to access the REST API. You can define more than one repository in "settings.json". Below is a simple example for "example.edu"'s authors repository.

```json
    {
        "authors": {
            "dsn": "USERNAME:SECRET@/authors",
            "base_url": "https://authors.example.edu",
            "write": true
        }
    }
```

In the "dsn" attribute __USERNAME:SECRET__ are the username/password for accessing the database. In the "rest" attribute the __USERNAME:SECRET__ are the username/password for accessing the REST API.


Related GitHub projects
-----------------------

- [py_dataset](https://github.com/caltechlibrary/py_dataset), This Python module provides access to dataset collections which we use as intermediate storage for JSON documents and related attachments.
- [AMES](https://github.com/caltechlibrary/ames), The eprintools command line programs have been made available to Python via the AMES project. This include support for both read and write to EPrints repository systems.


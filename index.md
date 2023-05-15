eprinttools
===========

This is a collection of command line tools and a web service written in Go for working with EPrints 3.3.x EPrint XML, the EPrint REST API and directly with the EPrints MySQL repository database(s). It is used by Caltech Library to render our <https://feeds.library.caltech.edu> website as well as for migrating content into a new repository system. Some of the command line tools maybe of more generatl interest while others are specific to Caltech Library's needs. Much of the test code presumes access to our repositories so is specific to our needs.

Go base code
------------

The programs:

- [eputil](eputil.1.md) is a command line utility for interacting (e.g. harvesting) JSON and XML from EPrints' REST API
    - minimal configuration (because it does so much less!)
- [epfmt](epfmt.1.md) is a command line utility to pretty print EPrints XML and convert to/from JSON including a simplified JSON inspired by DataCite and Invenion 3
- [doi2eprintxml](doi2eprintxml.1.md) is a command line program for turning metadata harvested from CrossRef and DataCite into an EPrint XML document based on one or more supplied DOI
- [ep3apid](ep3apid.1.md) is a Unix style web service for interacting with an EPrint repository via a localhost proxy. It includes the ability to get restricted key lists as well as retrieve a simplified JSON record representing an EPrints record
- [ep3harvester](ep3harvester.1.md) is an EPrints 3.x metadata harvesting tool working at the MySQL 8 level for EPrints content. It harvests the contents into a MySQL 8 database, one table per eprints repository storing the harvested metadata in JSON columns. This tool can also harvest CSV files with information for people and groups referenced in the EPrints repositories.
- [ep3genfeeds](ep3genfeeds.1.md) is used to genate the JSON documents that drive our feeds website.
- [ep3datasets](ep3datasets.1.md) is a tool to generate [dataset collections](https://github.com/caltechlibrary/dataset) from previously harvested EPrints repositories

Use cases
---------

Two primary use cases have driven development of EPrinttools

1. Reusing the metadata and content in our EPrints 3.3.16 repositories (see [Caltech Library Feeds](https://feeds.library.caltech.edu)
2. Populating our EPrints repository from standardize data sources (see [Acacia Project](https://github.com/caltechlibrary/Acacia)).


Related GitHub projects
-----------------------

- [py_dataset](https://github.com/caltechlibrary/py_dataset), This Python module provides access to dataset collections which we use as intermediate storage for JSON documents and related attachments.
- [AMES](https://github.com/caltechlibrary/ames), The eprintools command line programs have been made available to Python via the AMES project. This include support for both read and write to EPrints repository systems.


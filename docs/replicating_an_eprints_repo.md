
# Replicating an EPrints repository

This article is about using two Python modules that are
part of eprinttools to replicate the public facing content
of an EPrints repository into an S3 bucket hosing a
static website version of the content.

Taking this approach the EPrints software becomes primarily
a curration platform while the public view can function
indenpent of the EPrints collection (e.g. it can be up when
you have EPrints down for maintenance).

This approach combines several Caltech Library tools projects

+ [eprinttools](https://github.com/caltechlibrary/eprinttools/releases)
+ [dataset](https://github.com/caltechlibrary/dataset/releases)
+ [py_dataset](https://github.com/caltechlibrary/py_dataset/releases)
+ [mkpage](https://github.com/caltechlibrary/mkpage/releases)

## Setting up a staging service

There are several steps to setup a staging service to populate
your S3 bucket (or compatible object store).

1. Make sure you have access to the EPrints REST API by a service account
2. Setup a S3 bucket (or other compatible object store) to be your website
3. Install the following command line utilities __dataset__, __eprinttools__ (cli) and __mkpage__ into your a directory in your PATH (e.g. /usr/local/bin)
4. Install __py_dataset__ (e.g. pip install py_dataset)
5. copy __eprints3x__ and __eprintviews__ modules to the staging location
6. copy harvester_full.py, harvester_recent.py, genviews.py, indexer.py, mk_website.py and publisher.py to the staging location
    + these use __eprints3x__ and __eprintviews__ modules.
7. Create, update the necessary configuration files
8. Follow the instruction to run the build and deploy process
9. Initialize the dataset collection if needed

## Configuration and initialization

### JSON configuration file

A JSON file is used to run the example Python3 programs for 
replicating your public EPrints content. Here's an example
which I've named `config.json` in the examples below.

```json
    {
      "eprint_url": "https://USER:SECERT@eprints.example.edu",
      "dataset": "DATASET_COLLECTION",
      "number_of_days": -7,
      "control_item": "https://eprints.example.edu/cgi/users/home?screen=EPrint::View&eprintid=",
      "users": "users.json",
      "subjects": "subjects.txt",
      "views": "views.json",
      "organization": "Example Library IR",
      "site_title": "A EPrints Repository",
      "site_welcome": "Welcome to an EPrints Repository",
      "distribution_id": "",
      "bucket": "",
      "htdocs": "htdocs",
      "templates": "templates",
      "static": "static"
    }
```

Field explanations

eprint_url
: (required) is a URL with login credentials authorized to access the EPrints REST API

dataset
: (required) name of your dataset collection you've initialized with the dataset command

number_of_days
: Is an integer used to calculate the "recent" number of days to harvest

control_item
: (used in templates) This is a link (minus the EPrints ID) to use to access the EPrint page needed to edit/manage the item

users
: (required) Points at the JSON export of the users in the EPrint repository, it is needed to get the names of depositors

subjects
: (required) Points at a copy of a plain text file found in your EPrint under `/archives/REPO_ID/cfg/subject`, it is used to map subject views to paths and names

views
: (required) Points at a JSON file that shows a path part (e.g. ids) and label to use for that view (e.g. "Eprint ID").

organization
: (used in templates) The name of your organization (used by the Pandoc templates)

site_title
: (used in templates) The website title (used by the Pandoc templates)

site_welcome
: (used in templates) The website welcome statement (used by the Pandoc templates)

htdocs
: (has default value) This is the directory to host the website (or replicate from), it functions as your document root

templates
: (has default value) This is the directory that holds your website templates

static
: (has default value) This is the directory that holds your static files and assets (e.g. css, favicon, non-content images like logos)

bucket
: (optional) This is a URI To the S3 (or S3 like) bucket to host the public static website if you are using **publisher.py** and **invalidate_cloudfront.py** programs

distribution_id
: (optional) This is the ID number used by Cloudfront for invalidating CDN cache. It is only used by **includate_cloudfront.py**. It is only used by **includate_cloudfront.py**

### views JSON file content

This file controls what views get built built by the `genviews.py` command.
It is a key value of path and label. The path corresponds the to view
aggregation supported in `eprintviews/aggregator.py`.


```json
    {
        "ids": "Eprint ID", 
        "year": "Year",
        "person-az": "Person", 
        "event": "Conference", 
        "collection": "Collection", 
        "latest": "Latest Additions", 
        "publication": "Publication", 
        "issn": "ISSN", 
        "person": "Person", 
        "types": "Type", 
        "subjects": "Subjects"
    }
```

### subjects text file

This is a copy of the subjects file used in your EPrints 3.3.x 
installation. It is used to gather the table values that map
to the display names of the subjects.

This file is found by looking in the repository's configuration
(e.g. if the repository is called "instrepo" the path is probably
`archives/instrepo/cfg/subjects`).

### users JSON file

The users JSON file is created by exporting your current
users via the Admin user search found in your EPrints repository.
The search for select all users then export as JSON.

**IMPORTANT: This file should not be in a publicly readable location**

## Build Process

"config.json" is the name of the configuraiton file used for
example purposes. It can be have any name, the JSON content
is what is important.

1. harvest EPrints into a dataset collection
    a. `harvester_full.py config.json`
    b. `harvester_recent.py config.json`
2. generate the views and landing pages directory and metadata
    a. `genviews.py config.json`
3. index the metadata for search
    a. `indexer.py config.json`
4. make the website pages
    a. `mk_website.py config.json`
5. publisher your htdocs to your S3 bucket
    a. `publisher.py config.json`
    b. (if needed) `invalidate_cloudfront.py config.json`




# Release Notes

## Release Notes v0.1.7

+ Improved Aggregators, released via PyPI
+ Currently not harvesting digital documents in dataset collection


## Release Notes v0.1.4

+ Started eprinttools transition to Python
+ Added eprints3x Python module
+ Added eprintviews Python module
+ Added example programs for replicating EPrints public views as static website
    + harvester_full.py, harvester_recent.py
    + genviews.py
    + indexer.py
    + mk_website.py
    + publisher.py
    + invalidate_cloudfront.py

## Release Notes v0.1.1

+ Improve detection of primary object
+ Added ability to retrieve document if given full URL


## Release Notes v0.1.0

+ Dropped Go lang support for dataset integration
    + Removed `ep` command line tool
    + Removed github.com/caltechlibrary/eprinttools/harvest
    + Removed dependency and directory support for dataset collections
+ Copied github.com/caltechlibrary/rc into this repository
+ Added eprints3x object and harvester.py replacing `ep` tool


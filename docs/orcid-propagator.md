
# ORCID Propagation

The _orcid-propagator_ is a service that runs periodically from
a cronjob. It makes three steps over the collection.

## Step 1

Traverse the EPrints collection, create/update a dataset collection 
keyed by creator ID saving the metadata of eprint id, creator id and
orcid value

```json
    {
        "creator_id": "R-S-Doiel",
        "orcid": "0000-0003-0900-6903",
        "ambigious": false,
        "confirmed": true,
        "eprints": [
            {"eprint_id": 99999, "orcid": ""},
            {"eprint_id": 199999, "orcid": "0000-0003-0900-6903"}
        ]
    }
```



## Step 2 

Traverse the dataset collection of creator ids, review each
eprint id and associated orcid foudn for ambiguity, if orcid
is unambious then proceed to step 3

## Step 3

For each creator record's eprint entry that is missing an
unambious orcid retreive the full EPrint XML, update the creator
id's ORCID value and re-import as update the EPrint XML called
the epadmin appropriately.

## Step 4

After traversing the creator id collection  for each ambigious
ORCID assignment print out the EPrint Record ID link, the
Creator ID, a list of possible ORCIDs needing correction.


Table Notes
===========

This documents holds some notes related to how EPrint's SQL table structure is organized.

eprint
: Primary eprint record table. This table provides the core single values produced in the EPrint XML

eprint_*_*
: These tables handle "item" content in the EPrint XML. A good example is "eprint_creators_name" which holds name information about a creator.

file
: This table describes the files that make up a document attached to an EPrint record

document and document_*
: These tables describes the documents associatd with a specific eprint record.  In our repositories only document and document_relation_type have any rows

history
: Holds change descriptions for files and document, also used to generate email back to patrons when their deposits fail



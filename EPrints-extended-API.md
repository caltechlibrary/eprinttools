EPrints extended API
====================

The EPrints software package from University of Southampton
provides a rich internal Perl API along with a RESTful web API.
The latter has been used extensively by Caltech Library to facilitate 
content across campus for our various EPrints repositories. (See priorities 
two and three of the Caltech Library's strategic plan). The challenge now 
is to move beyond the present limitations.


Extending EPrints directly is error prone and cumbersome. Implementing 
features in Perl safely is only the start of trouble if we modify EPrints 
directly. In contrast EPrints' MySQL database structure has proven to be 
durable and predictable. MySQL can be leverage directly to extended API 
seeks to beyond our current constraints.

What should an extended web API look like?

Design considerations
---------------------

- The extended API should be web accessible to support data platforms such as feeds.library.caltech.edu as well as our growing cast of application hosted on apps.library.caltech.edu
- It needs in interact with MySQL's EPrints database safely, e.g. be read only
- Minimize the load on EPrints' MySQL database, favor simple SQL queries, perhaps limiting them to single table scans
- Be near zero management, it should run as a service that doesn't require on going interventions and integrate into DLD's existing monitoring infrastructure

An extended API should provide a limited web service that maps URL end 
points to simple MySQL queries run against the various EPrints databases.
The service should be easy to implement require minimal resources, e.g. one prepared SQL statement per end point.

Security and privacy should be for front and center when implementing
any web service. By returning EPrint ID only we limit the risk of exposing 
in appropriate metadata (e.g. author information). The EPrint ID is an 
integer without specific meaning. It does not give you access to sensitive 
information.


Unique IDs to EPrint IDs
------------------------

The following URL end points are intended to take one unique identifier and
map that to one or more EPrint IDs. This can be done because each unique ID 
targeted can be identified by querying a single table in EPrints.  In 
addition the scan can return the complete results since all EPrint IDs are 
integers and returning all EPrint IDs in any of our repositories is 
sufficiently small to be returned in a single HTTP request.

- `/<REPO_ID>/doi/<DOI>` with the adoption of EPrints "doi" field in the EPrint table it makes sense to have a quick translation of DOI to EPrint id for a given EPrints repository. 
- `/<REPO_ID>/creator-id/<CREATOR_ID>` scans the name creator id field associated with creators and returns a list of EPrint ID 
- `/<REPO_ID>/creator-orcid/<ORCID>` scans the "orcid" field associated with creators and returns a list of EPrint ID 
- `/<REPO_ID>/editor-id/<CREATOR_ID>` scans the name creator id field associated with editors and returns a list of EPrint ID 
- `/<REPO_ID>/contributor-id/<CONTRIBUTOR_ID>` scans the "id" field associated with a contributors and returns a list of EPrint ID 
- `/<REPO_ID>/advisor-id/<ADVISOR_ID>` scans the name advisor id field associated with advisors and returns a list of EPrint ID 
- `/<REPO_ID>/committee-id/<COMMITTEE_ID>` scans the committee id field associated with committee members and returns a list of EPrint ID
- `/<REPO_ID>/group-id/<GROUP_ID>` this scans group ID and returns a list of EPrint IDs associated with the group
- `/<REPO_ID>/funder-id/<FUNDER_ID>` returns a list of EPrint IDs associated with the funder's ROR
- `/<REPO_ID>/grant-number/<GRANT_NUMBER>` returns a list of EPrint IDs associated with the grant number

Change Events
-------------

The follow API end points would facilitate faster updates to our feeds 
platform as well as allow us to create a separate public view of our EPrint 
repository content.

- `/<REPO_ID>/updated/<TIMESTAMP>` returns a list of EPrint IDs updated starting at timestamp (timestamps should have a resolution to the minute, e.g. "YYYY-MM-DD HH:MM:SS")
- `/<REPO_ID>/deleted/<TIMESTAMP>` returns a list of EPrint IDs deleted starting at timestamp
- `/<REPO_ID>/pubdate/<APROX_DATESTAMP>` this query scans the EPrint table for publication date where the day may be approximate, e.g. just the year, just the year and month or year, month and day. It returns one or more EPrint IDs.

Nice to have end points
-----------------------

The following end points should also be easily implemented and should be considered if useful to other projects or staff.

- `/<REPO_ID>/editor-orcid/<ORCID>` scans the "orcid" field associated with a editors and returns a list of EPrint ID 
- `/<REPO_ID>/contributor-orcid/<ORCID>` scans the "orcid" field associated with a contributors and returns a list of EPrint ID 
- `/<REPO_ID>/advisor-orcid/<ORCID>` scans the "orcid" field associated with advisors and returns a list of EPrint ID
- `/<REPO_ID>/committee-orcid/<ORCID>` scans the "orcid" field associated with committee members and returns a list of EPrint ID
- `/<REPO_ID>/creator-name/<FAMILY_NAME>/<GIVEN_NAME>` scans the name fields associated with creators and returns a list of EPrint ID 
- `/<REPO_ID>/editor-name/<FAMILY_NAME>/<GIVEN_NAME>` scans the family and given name field associated with a editors and returns a list of EPrint ID 
- `/<REPO_ID>/contributor-name/<FAMILY_NAME>/<GIVEN_NAME>` scans the family and given name field associated with a contributors and returns a list of EPrint ID 
- `/<REPO_ID>/advisor-name/<FAMILY_NAME>/<GIVEN_NAME>` scans the name fields associated with advisors returns a list of EPrint ID 
- `/<REPO_ID>/committee-name/<FAMILY_NAME>/<GIVEN_NAME>` scans the family and given name fields associated with committee members and returns a list of EPrint ID
- `/<REPO_ID>/group-ror/<ROR>` this scans the local group ROR related fields and returns a list of EPrint ids.
- `/<REPO_ID>/funder-ror/<FUNDER_ROR>` returns a list of EPrint IDs associated with the funder's ROR
- `/<REPO_ID>/updated/<TIMESTAMP>/<TIMESTAMP>` returns a list of EPrint IDs updated from start timestamp to  end timestamp
- `/<REPO_ID>/deleted/<TIMESTAMP>/<TIMESTAMP>` returns a list of EPrint IDs deleted from start timestamp to end timestamp
- `/<REPO_ID>/pubmed/<PUBMED_ID>` returns a list of EPrint IDs associated with the PubMed ID
- `/<REPO_ID>/issn/<ISSN>` returns a list of EPrint IDs associated with the ISSN
- `/<REPO_ID>/isbn/<ISSN>` returns a list of EPrint IDs associated with the ISSN


Someday, maybe
--------------

EPrints XML is complex and hard to work with. A simplified data structure 
could make working with our repository data much easier. If user/role 
strictions were enforced in an extended EPrints API it could provide a 
clean JSON expression of a more general biblographic record. At that stage 
it might also be desirable to allow updates to existing EPrints records via 
the extended API.


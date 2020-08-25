
# Views in our various EPrints repositories

## Notes from Google Scholar

It looks likes /view/ids/ and /view/year/ maybe used by the GoogleBot
as Google Scholar understands Eprints' directory layout even though those aren't usually linked from the homepage.  We also want to generate a "latest additions" link so that larger repositories can have their additions and updates picked up quicker. Google Scholar docs are at https://scholar.google.com/intl/en/scholar/inclusion.html#indexing.

## Questions

+ People, person, persion-az include the creators.item values and sometimes include contributors or editors values, should this be consistent?
+ We have lots of pays to describe an EPrints item type, should be standardize? (e.g. Type vs. Document Types)

## Common views

ids
: is based on the eprint id, appears to be a default view but not usually included in browseviews.html (results is a list of list of one item)

publication
: is based on publication simple field, many repositories have this, labeling of view varies but path is consistent

issn
: is base on issn simple field, many repositories have this, labeling is uppercase of field name

group
: is based on local_group.item, many repostiroes, labeling is usually Group 

person
: is based on creators, is in many repositories but not often in browseview.html, seems to be the same as people view, in some cases (Campus Pubs) it appears to include people with other affilations such as editor and contributors

person-az
: is based on creators, is in most repositories, often in browseview labed Person

collection
: is based on collection simple field, should probably be in all repositories but isn't listed often in browseview.html

## CaltechAUTHORS

### browseview.html

+ Person -> person-az
+ Year -> year
+ Document Type -> types
+ Research Group -> group
+ Collection -> collection
+ Latest Additions -> cgi/latest (not sure that this is a standard view, it may do a live SQL query each time it is hit given it's slow response time)

### /view/

+ Eprint ID -> ids
+ Journal or Publication Title -> publication
+ ISSN -> issn
+ Year -> year
+ Group -> group
+ Person -> person
+ Type -> types
+ Person -> person-az
+ Collection -> collection

## CaltechTHESIS

### browseview.html

+ Author -> author (this appears to be the same as person-az in other repositories)
+ Committee Members -> committee
+ Research Advisor -> advisor
+ Option (Field of Study) -> option
+ Degree/Thesis Type -> degree
+ Year -> year
+ Latest Additions -> cgi/latest (like CaltechAUTHORS)
+ Research Group -> group

### /view/

+ Eprint ID -> ids
+ Year -> year
+ Option -> option  (this is a custom field in thesis)
+ Degree/Thesis Type -> degree (this appears semi-custom)
+ Author -> author (this is really based on the creator field and is probably the same code as person in the other repositories)
+ Advisor -> advisor (this is based on thesis_advisor.item)
+ Committee Member -> committee (this is based on thesis_committee.item)
+ Group -> group (this is based on local_group.item)

## CaltechLabNotes

### browseview.html

+ Project Information -> /information.html (this about page link)
+ Browse by Name -> person
+ Browse by Subject -> subjects
+ Search -> /searchoptions.html (MySQL based simple and advanced search pages)

### /view/

+ Eprint ID -> ids
+ Subject -> subjects
+ Name -> person

## Caltech Magazine (calteches)

### browseview.html

+ Year -> year
+ Item Category -> subjects
+ Author -> person-az
+ Latest Additionas -> /cgi/latest

### /view/

+ Eprint ID -> ids
+ Year -> year
+ Item Category -> subjects
+ Person -> person
+ Author -> person-az
+ Type -> types

## CaltechCONF

### browseview.html

+ Person -> person-az
+ Year -> year
+ Conference -> event
+ Collection -> collection
+ Latest Additions -> /cgi/latest

### /view/

+ Eprint ID -> ids
+ Journal or Publication Title -> publication
+ ISSN -> issn
+ Year -> year
+ Conference -> event
+ Person -> person
+ Type -> types
+ Person -> person-az
+ Collection -> collection

## Oral Histories (caltechoh)

### browseview.html

+ Person -> person-az
+ Year -> year
+ Subject -> subjects

### /view/

+ Eprint ID -> ids
+ Year -> year
+ Subject -> subjects
+ Type -> types
+ Person -> person-az

## CampusPubs

### browseview.html

+ Publication Title -> publication
+ Year -> year
+ Issuing Body -> issuing_body
+ Author/Contribor -> person-az
+ Document Type -> types
+ Latest Additions -> /cgi/latest

### /view/

+ Eprint ID -> ids
+ Publication Title -> publication
+ Issuing Body -> issuing_body
+ Type - types
+ Author/Contributor Name -> person-az

issuing_body 
: based on corp_creators.item where id is path and name is label 



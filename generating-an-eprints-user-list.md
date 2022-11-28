
# Generating an EPrints user list

## Problem

The depositor value for an EPrints uses the EPrints Users 
values. I don't think it is a good idea to expose this part
of the REST API.  On the other hand users are not frequently
added or remove and EPrints provides an easy way to dump
a user list which we can leverage when resolving the EPrint XML
with who created the item.

## Solution

Authenticate into your EPrints repository as a repository
administrator. Click on the "Admin" link, Click on the "Search Users"
button, Check all the User Type check boxes and press the search button.
This will give you a results list. For the export option from the
results pick JSON (or the whatever format you like working with).
The resulting page can be saved and is available to integrate
into your application when decoding the EPrint XML values
for `userid`.

    


Composit Views
--------------

Composit views are lists of available values that can then be used to get lists of eprint ids. These are the lists typically found in the "Browse" page of EPrints.

- '/{REPO_ID}/year' returns a list of years that eprint records cover
- '/{REPO_ID}/year/{YEAR}' returns a list of eprint ids related to year
- '/{REPO_ID}/person' returns a list of objects with a person id attribute s and relationship that are in eprint records. E.g. '{"peron_id": "Doe-J": "is_a": [ "creator", "contributor" ]}'
- '/{REPO_ID}/person/{PERSON_ID}' returns a list of eprint ids related to person id, maybe related as creator, editor, contributor, advisor or committee member.
- '/{REPO_ID}/collection' returns a list of collection id
- '/{REPO_ID}/collection/{COLLECTION}' returns a list eprint ids associated with the collection id
- '/{REPO_ID}/type' returns a list of document types in repository
- '/{REPO_ID}/type/{TYPE}' returns a list eprint id associated with document types in repository
- '/{REPO_ID}/publications' returns a list of publications
- '/{REPO_ID}/publications/{PUBLICATION}' returns a list eprint id of publications
- '/{REPO_ID}/publisher' returns a list of publisher
- '/{REPO_ID}/publisher/{PUBLISHER}' returns a list eprint id of publisher




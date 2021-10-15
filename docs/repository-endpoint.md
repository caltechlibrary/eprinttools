Repository (end point)
------------------

The end point executes a sequence of "SHOW" SQL statements to build
a JSON object with table names as attributes pointing at an array of column names. This is suitable to determine on a per repository bases the
related table and columnames representing an EPrint record.

- '/<REPO_ID>/repository' return the JSON representation
- '/<REPO_ID>/repository/help' return this documentation

Example
-------


```shell
   curl http://localhost:8485/lemurprints/tables
```

Would return a JSON expression similar to the expression below.
The object has attributes that map to the EPrint talbles and for
each table the attribute points at an array of column names.

The "eprint" table is the root of the object. Each other table
is a sub object or array. Tables containing a "pos" field 
render as an array of objects (e.g. the "item" elements in the
EPrint XML). If pos is missing then it is an object with
attributes and values.

Each table is relatated by the "eprintid" column.


```json
{
  "eprint": [ "abstract", "alt_url", "book_title", "collection", "commentary", "completion_time", "composition_type", "contact_email", "coverage_dates", "data_type", "date_day", "date_month", "date_type", "date_year", "datestamp_day", "datestamp_hour", "datestamp_minute", "datestamp_month", "datestamp_second", "datestamp_year", "department", "dir", "doi", "edit_lock_since", "edit_lock_until", "edit_lock_user", "edition", "eprint_status", "eprintid", "errata", "event_dates", "event_location", "event_title", "event_type", "fileinfo", "full_text_status", "hide_thesis_author_email", "id_number", "importid", "institution", "interviewdate", "interviewer", "isbn", "ispublished", "issn", "item_issues_count", "keywords", "lastmod_day", "lastmod_hour", "lastmod_minute", "lastmod_month", "lastmod_second", "lastmod_year", "latitude", "learning_level", "longitude", "metadata_visibility", "monograph_type", "note", "num_pieces", "number", "official_cit", "official_url", "output_media", "pagerange", "pages", "parent_url", "patent_applicant", "pedagogic_type", "place_of_pub", "pmc_id", "pres_type", "publication", "publisher", "refereed", "replacedby", "rev_number", "review_status", "reviewer", "rights", "series", "source", "status_changed_day", "status_changed_hour", "status_changed_minute", "status_changed_month", "status_changed_second", "status_changed_year", "succeeds", "suggestions", "sword_depositor", "sword_slug", "task_purpose", "thesis_approved_date_day", "thesis_approved_date_month", "thesis_approved_date_year", "thesis_author_email", "thesis_defense_date_day", "thesis_defense_date_month", "thesis_defense_date_year", "thesis_degree_date_day", "thesis_degree_date_month", "thesis_degree_date_year", "thesis_degree_grantor", "thesis_public_date_day", "thesis_public_date_month", "thesis_public_date_year", "thesis_type", "title", "toc", "type", "userid", "volume" ],
  "eprint_accompaniment": [ "accompaniment", "eprintid", "pos" ],
  "eprint_alt_title": [ "alt_title", "eprintid", "pos" ],
  "eprint_conductors_id": [ "conductors_id", "eprintid", "pos" ],
  "eprint_conductors_name": [ "conductors_name_family", "conductors_name_given", "conductors_name_honourific", "conductors_name_lineage", "eprintid", "pos" ],
  "eprint_conductors_uri": [ "conductors_uri", "eprintid", "pos" ],
  "eprint_conf_creators_id": [ "conf_creators_id", "eprintid", "pos" ],
  "eprint_conf_creators_name": [ "conf_creators_name", "eprintid", "pos" ],
  "eprint_conf_creators_uri": [ "conf_creators_uri", "eprintid", "pos" ],
  "eprint_contributors_id": [ "contributors_id", "eprintid", "pos" ],
  "eprint_contributors_name": [ "contributors_name_family", "contributors_name_given", "contributors_name_honourific", "contributors_name_lineage", "eprintid", "pos" ],
  "eprint_contributors_type": [ "contributors_type", "eprintid", "pos" ], "eprint_contributors_uri": [ "contributors_uri", "eprintid", "pos" ],
  "eprint_copyright_holders": [ "copyright_holders", "eprintid", "pos" ],
  "eprint_corp_creators_id": [ "corp_creators_id", "eprintid", "pos" ],
  "eprint_corp_creators_name": [ "corp_creators_name", "eprintid", "pos" ],
  "eprint_corp_creators_uri": [ "corp_creators_uri", "eprintid", "pos" ],
  "eprint_creators_id": [ "creators_id", "eprintid", "pos" ],
  "eprint_creators_name": [ "creators_name_family", "creators_name_given", "creators_name_honourific", "creators_name_lineage", "eprintid", "pos" ],
  "eprint_creators_uri": [ "creators_uri", "eprintid", "pos" ],
  "eprint_divisions": [ "divisions", "eprintid", "pos" ],
  "eprint_editors_id": [ "editors_id", "eprintid", "pos" ],
  "eprint_editors_name": [ "editors_name_family", "editors_name_given", "editors_name_honourific", "editors_name_lineage", "eprintid", "pos" ],
  "eprint_editors_uri": [ "editors_uri", "eprintid", "pos" ],
  "eprint_exhibitors_id": [ "eprintid", "exhibitors_id", "pos" ],
  "eprint_exhibitors_name": [ "eprintid", "exhibitors_name_family", "exhibitors_name_given", "exhibitors_name_honourific", "exhibitors_name_lineage", "pos" ],
  "eprint_exhibitors_uri": [ "eprintid", "exhibitors_uri", "pos" ],
  "eprint_funders_agency": [ "eprintid", "funders_agency", "pos" ],
  "eprint_funders_grant_number": [ "eprintid", "funders_grant_number", "pos" ],
  "eprint_item_issues_comment": [ "eprintid", "item_issues_comment", "pos" ],
  "eprint_item_issues_description": [ "eprintid", "item_issues_description", "pos" ],
  "eprint_item_issues_id": [ "eprintid", "item_issues_id", "pos" ],
  "eprint_item_issues_reported_by": [ "eprintid", "item_issues_reported_by", "pos" ],
  "eprint_item_issues_resolved_by": [ "eprintid", "item_issues_resolved_by", "pos" ],
  "eprint_item_issues_status": [ "eprintid", "item_issues_status", "pos" ],
  "eprint_item_issues_timestamp": [ "eprintid", "item_issues_timestamp_day", "item_issues_timestamp_hour", "item_issues_timestamp_minute", "item_issues_timestamp_month", "item_issues_timestamp_second", "item_issues_timestamp_year", "pos" ],
  "eprint_item_issues_type": [ "eprintid", "item_issues_type", "pos" ],
  "eprint_keyword": [ "eprintid", "keyword", "pos" ],
  "eprint_lyricists_id": [ "eprintid", "lyricists_id", "pos" ],
  "eprint_lyricists_name": [ "eprintid", "lyricists_name_family", "lyricists_name_given", "lyricists_name_honourific", "lyricists_name_lineage", "pos" ],
  "eprint_lyricists_uri": [ "eprintid", "lyricists_uri", "pos" ],
  "eprint_other_numbering_system_id": [ "eprintid", "other_numbering_system_id", "pos" ],
  "eprint_other_numbering_system_name": [ "eprintid", "other_numbering_system_name", "pos" ],
  "eprint_producers_id": [ "eprintid", "pos", "producers_id" ],
  "eprint_producers_name": [ "eprintid", "pos", "producers_name_family", "producers_name_given", "producers_name_honourific", "producers_name_lineage" ],
  "eprint_producers_uri": [ "eprintid", "pos", "producers_uri" ],
  "eprint_projects": [ "eprintid", "pos", "projects" ],
  "eprint_reference": [ "eprintid", "pos", "reference" ],
  "eprint_referencetext": [ "eprintid", "pos", "referencetext" ],
  "eprint_related_url_description": [ "eprintid", "pos", "related_url_description" ],
  "eprint_related_url_type": [ "eprintid", "pos", "related_url_type" ],
  "eprint_related_url_url": [ "eprintid", "pos", "related_url_url" ],
  "eprint_relation_type": [ "eprintid", "pos", "relation_type" ],
  "eprint_relation_uri": [ "eprintid", "pos", "relation_uri" ],
  "eprint_skill_areas": [ "eprintid", "pos", "skill_areas" ],
  "eprint_subjects": [ "eprintid", "pos", "subjects" ],
  "eprint_thesis_advisor_email": [ "eprintid", "pos", "thesis_advisor_email" ],
  "eprint_thesis_advisor_id": [ "eprintid", "pos", "thesis_advisor_id" ],
  "eprint_thesis_advisor_name": [ "eprintid", "pos", "thesis_advisor_name_family", "thesis_advisor_name_given", "thesis_advisor_name_honourific", "thesis_advisor_name_lineage" ],
  "eprint_thesis_advisor_uri": [ "eprintid", "pos", "thesis_advisor_uri" ],
  "eprint_thesis_committee_email": [ "eprintid", "pos", "thesis_committee_email" ],
  "eprint_thesis_committee_id": [ "eprintid", "pos", "thesis_committee_id" ],
  "eprint_thesis_committee_name": [ "eprintid", "pos", "thesis_committee_name_family", "thesis_committee_name_given", "thesis_committee_name_honourific", "thesis_committee_name_lineage" ],
  "eprint_thesis_committee_role": [ "eprintid", "pos", "thesis_committee_role" ],
  "eprint_thesis_committee_uri": [ "eprintid", "pos", "thesis_committee_uri" ]
}
```


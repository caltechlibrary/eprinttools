#!/bin/bash

cat <<SQL >not_published_not_submitted_report.sql
select eprintid, date_type
from eprint
where date_type != 'published' and date_type != 'submitted'
order by date_type
SQL

mysql --batch caltechauthors < not_published_not_submitted_report.sql | tab2csv >not_published_not_submitted_report.csv

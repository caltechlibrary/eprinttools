#!/bin/bash

cat <<SQL >inpress_report.sql
select eprintid, date_type
from eprint
where date_type = 'inpress'
SQL

mysql --batch caltechauthors < inpress_report.sql | tab2csv >inpress_report.csv
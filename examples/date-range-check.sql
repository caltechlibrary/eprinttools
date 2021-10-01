SELECT eprintid FROM caltechauthors
WHERE (CONCAT(lastmod_year, "-",
LPAD(IFNULL(lastmod_month, 1), 2, "0"), "-",
LPAD(IFNULL(lastmod_day, 1), 2,"0"), " ",
LPAD(IFNULL(lastmod_hour, 0), 2, "0"), ":",
LPAD(IFNULL(lastmod_minute, 0), 2, "0"), ":",
LPAD(IFNULL(lastmod_second, 0), 2, "0")) >= ?) AND
(CONCAT(lastmod_year, "-",
LPAD(IFNULL(lastmod_month, 12), 2, "0"), "-",
LPAD(IFNULL(lastmod_day, 28), 2, "0"), " ",
LPAD(IFNULL(lastmod_hour, 23), 2, "0"), ":",
LPAD(IFNULL(lastmod_minute, 59), 2, "0"), ":",
LPAD(IFNULL(lastmod_second, 59), 2,"0")) <= ?);


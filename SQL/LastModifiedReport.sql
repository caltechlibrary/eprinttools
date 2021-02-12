--
-- LastModifiedReport.sql generates a listing of keys and modified datestamps
-- from an EPrints 3.3.x database.
--
-- @Author R. S. Doiel <rsdoiel@library.caltech.edu>
--
-- Example MySQL Shell Usage
--   CALL Last_Modified_Report(2021,2,1);
--
-- Example Shell Usage to generate a CSV version of report
--
--  mysql --batch -e 'CALL Last_Modified_Report(2021, 02, 01)' caltechauthors \
--       tr '\t' ','
--
-- Tested with MySQL 5.1 (coda3) and MySQL 8 (eprints.library.caltech.edu)
--
DROP PROCEDURE IF EXISTS Last_Modified_Report;
delimiter //
CREATE PROCEDURE Last_Modified_Report (year INT, month INT, day INT)
BEGIN
    SELECT eprintid AS eprint_id, 
           CONCAT(lastmod_year,"-",
                  LPAD(lastmod_month, 2, "0"),"-",
                  LPAD(lastmod_day, 2, "0")," ",
                  LPAD(lastmod_hour, 2, "0"), ":",
                  LPAD(lastmod_minute, 2, "0"), ":00") AS lastmod_date
           FROM eprint
           WHERE (lastmod_year >= year) AND 
                 (lastmod_month >= month) AND (lastmod_day >= day)
           ORDER BY lastmod_year, lastmod_month, 
                    lastmod_day, lastmod_hour, lastmod_minute ASC;
END;
//
delimiter ;


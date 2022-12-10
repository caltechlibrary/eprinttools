package eprinttools

//
// ep3sql.go provides crosswalk methods to/from SQL
//
import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	// Caltech Packages
	"github.com/caltechlibrary/eprinttools/cleaner"
	"github.com/caltechlibrary/pairtree"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// timestamp holds the Format of a MySQL time field
	timestamp = `2006-01-02 15:04:05`
	datestamp = `2006-01-02`
)

//
// DB SQL functions.
//

// OpenConnections
func OpenConnections(config *Config) error {
	if config.Connections == nil {
		config.Connections = map[string]*sql.DB{}
	}
	for repoID, dataSource := range config.Repositories {
		dataSourceName := dataSource.DSN
		// Setup DB connection for target repository
		db, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			return fmt.Errorf("could not open MySQL connection for %s, %s", repoID, err)
		}
		config.Connections[repoID] = db
		dataSource.TableMap, err = eprintTablesAndColumns(db, repoID)
		if err != nil {
			return fmt.Errorf("failed to map table and columns for %q, %s", repoID, err)
		}
	}
	return nil
}

// CloseConnections
func CloseConnections(config *Config) error {
	var (
		errors []string
	)
	if config.Connections == nil {
		config.Connections = map[string]*sql.DB{}
		return fmt.Errorf("no connections defined")
	}
	for name, db := range config.Connections {
		if err := db.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to close %s, %s", name, err))
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}
	return nil
}

// sqlQueryInts takes a repostory ID, a SQL statement and returns
// intergers retrieved.
func sqlQueryInts(config *Config, repoID string, stmt string) ([]int, error) {
	if db, ok := config.Connections[repoID]; ok {
		rows, err := db.Query(stmt)
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		defer rows.Close()
		value := 0
		values := []int{}
		for rows.Next() {
			err := rows.Scan(&value)
			if err == nil {
				values = append(values, value)
			} else {
				return nil, fmt.Errorf("ERROR: scan error (%q), %s", repoID, err)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("ERROR: rows error (%q), %s", repoID, err)
		}
		return values, nil
	}
	return nil, fmt.Errorf("bad request")
}

// sqlQueryIntIDs takes a repostory ID, a SQL statement and applies
// the args returning a list of integer id or error.
func sqlQueryIntIDs(config *Config, repoID string, stmt string, args ...interface{}) ([]int, error) {
	if db, ok := config.Connections[repoID]; ok {
		rows, err := db.Query(stmt, args...)
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		defer rows.Close()
		value := 0
		values := []int{}
		for rows.Next() {
			err := rows.Scan(&value)
			if (err == nil) && (value > 0) {
				values = append(values, value)
			} else {
				return nil, fmt.Errorf("ERROR: scan error (%q), %s", repoID, err)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("ERROR: rows error (%q), %s", repoID, err)
		}
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		return values, nil
	}
	return nil, fmt.Errorf("bad request")
}

// sqlQueryStringIDs takes a repostory ID, a SQL statement and applies
// the args returning a list of string type id or error.
func sqlQueryStringIDs(config *Config, repoID string, stmt string, args ...interface{}) ([]string, error) {
	if db, ok := config.Connections[repoID]; ok {
		rows, err := db.Query(stmt, args...)
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		defer rows.Close()
		value := ``
		values := []string{}
		for rows.Next() {
			err := rows.Scan(&value)
			if err == nil {
				values = append(values, value)
			} else {
				return nil, fmt.Errorf("ERROR: scan error (%q), %q, %s", repoID, stmt, err)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("ERROR: rows error (%q), %s", repoID, err)
		}
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		return values, nil
	}
	return nil, fmt.Errorf("bad request")
}

// IsPublic takes an EPrintID and returns true if public, false otherwise
//
// Check if an EPrint record "is public"
func IsPublic(config *Config, repoID string, eprintid int) (bool, error) {
	if db, ok := config.Connections[repoID]; ok {
		stmt := `SELECT IFNULL(eprint_status, '') AS status, IFNULL(metadata_visibility, '') AS visibility FROM eprint WHERE eprintid = ? LIMIT 1`
		rows, err := db.Query(stmt, eprintid)
		if err != nil {
			return false, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		defer rows.Close()
		var (
			status, visibility string
		)
		for rows.Next() {
			err := rows.Scan(&status, &visibility)
			if err != nil {
				return false, fmt.Errorf("ERROR: scan error (%q), %q, %s", repoID, stmt, err)
			}
		}
		if err := rows.Err(); err != nil {
			return false, fmt.Errorf("ERROR: rows error (%q), %s", repoID, err)
		}
		if err != nil {
			return false, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		return ((status == "archive") && (visibility == "show")), nil
	}
	return false, fmt.Errorf("bad request")
}

//
// Expose EPrint meta data structure
//

func GetTablesAndColumns(config *Config, repoID string) (map[string][]string, error) {
	if config.Connections == nil {
		return nil, fmt.Errorf(`database access not configured`)
	}
	if db, ok := config.Connections[repoID]; ok == true {
		return eprintTablesAndColumns(db, repoID)
	}
	return nil, fmt.Errorf(`database connections not defined for %s`, repoID)
}

//
// EPrint User Info
//

// GetUsernames returns a list of all usernames in a repository
func GetUsernames(config *Config, repoID string, userids ...int) ([]string, error) {
	stmt := `SELECT username FROM user ORDER BY userid`
	return sqlQueryStringIDs(config, repoID, stmt)
}

// GetUserID takes a username and returns a list of userid
func GetUserID(config *Config, repoID string, username string) ([]int, error) {
	stmt := `SELECT userid FROM user WHERE username = ?`
	return sqlQueryIntIDs(config, repoID, stmt, username)
}

// GetUserBy takes a field name (e.g. userid, username) and value
// and returns an EPrintUser object.
func GetUserBy(config *Config, repoID string, queryField string, queryValue interface{}) (*EPrintUser, error) {
	var (
		year, month, day, hour, minute, second int
		hideEMail                              string
	)
	if db, ok := config.Connections[repoID]; ok {
		stmt := fmt.Sprintf(`SELECT userid, username, usertype, IFNULL(name_honourific, '') AS honourific, IFNULL(name_family, '') AS family, IFNULL(name_given, '') AS given, IFNULL(name_lineage, '') AS lineage, IFNULL(email, '') AS email, IFNULL(hideemail, '') AS hideemail, IFNULL(dept, '') AS dept, IFNULL(org, '') AS org, IFNULL(address, '') AS address, IFNULL(country, '') AS country, IFNULL(joined_year, 0) AS joined_year, IFNULL(joined_month, 0) AS joined_month, IFNULL(joined_day, 0) AS joined_day, IFNULL(joined_hour, 0) AS joined_hour, IFNULL(joined_minute, 0) AS joined_minute, IFNULL(joined_second, 0) AS joined_second FROM user WHERE %s = ? LIMIT 1`, queryField)
		rows, err := db.Query(stmt, queryValue)
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		defer rows.Close()
		// Map values back into our object.
		user := new(EPrintUser)
		user.Name = new(Name)
		for rows.Next() {
			err := rows.Scan(&user.UserID, &user.Username, &user.Type,
				&user.Name.Honourific, &user.Name.Family,
				&user.Name.Given, &user.Name.Lineage,
				&user.EMail, &hideEMail,
				&user.Dept, &user.Org,
				&user.Address, &user.Country,
				&year, &month, &day, &hour, &minute, &second)
			if err != nil {
				return nil, fmt.Errorf("ERROR: scan error (%q), %q, %s", repoID, stmt, err)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("ERROR: rows error (%q), %s", repoID, err)
		}
		if err != nil {
			return nil, fmt.Errorf("ERROR: query error (%q), %s", repoID, err)
		}
		user.Joined = fmt.Sprintf(`%04d-%02d-%02d %02d:%02d:%02d`, year, month, day, hour, minute, second)
		if strings.ToLower(hideEMail) == "true" {
			user.HideEMail = true
		} else {
			user.HideEMail = false
		}
		return user, nil
	}
	return nil, fmt.Errorf(`bad request`)
}

func SQLCreateUser(config *Config, repoID string, user *EPrintUser) (int, error) {
	var (
		year, month, day, hour, minute, second int
	)
	now := time.Now()
	if user.Name == nil {
		user.Name = &Name{
			Honourific: ``,
			Family:     ``,
			Given:      ``,
			Lineage:    ``,
		}
	}
	if user.Joined == "" {
		user.Joined = now.Format(timestamp)
		year = now.Year()
		month = int(now.Month())
		day = now.Day()
		hour = now.Hour()
		minute = now.Minute()
		second = now.Second()
	} else {
		if dt, err := time.Parse(timestamp, user.Joined); err == nil {
			year = dt.Year()
			month = int(dt.Month())
			day = dt.Day()
			hour = dt.Hour()
			minute = dt.Minute()
			second = dt.Second()
		}
	}
	hideEMail := "FALSE"
	if user.HideEMail {
		hideEMail = "TRUE"
	}
	if db, ok := config.Connections[repoID]; ok {
		// First generate new row for user.
		stmt := `INSERT INTO user (userid) (SELECT (IFNULL((SELECT userid FROM user ORDER BY userid DESC LIMIT 1), 0) + 1) AS userid)`
		_, err := db.Exec(stmt)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		stmt = `SELECT userid FROM user ORDER BY userid DESC LIMIT 1`
		rows, err := db.Query(stmt)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		id := 0
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				return 0, fmt.Errorf(`could not calculate the new userid value, %s`, err)
			}
		}
		rows.Close()
		if err != nil {
			return 0, fmt.Errorf(`SQL failed to get insert id, %s`, err)
		}
		// Update user
		user.UserID = id
		stmt = `REPLACE INTO user (userid, username, usertype, name_honourific, name_family, name_given, name_lineage, email, hideemail, dept, org, address, country, joined_year, joined_month, joined_day, joined_hour, joined_minute, joined_second) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Exec(stmt, user.UserID, user.Username, user.Type, user.Name.Honourific, user.Name.Family, user.Name.Given, user.Name.Lineage, user.EMail, hideEMail, user.Dept, user.Org, user.Address, user.Country, year, month, day, hour, minute, second)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		return user.UserID, nil
	}
	return 0, fmt.Errorf(`bad request`)
}

func SQLUpdateUser(config *Config, repoID string, user *EPrintUser) error {
	if db, ok := config.Connections[repoID]; ok {
		hideEMail := "FALSE"
		if user.HideEMail {
			hideEMail = "TRUE"
		}
		stmt := `UPDATE user SET username = ?, usertype = ?, name_honourific = ?, name_family = ?, name_given = ?, name_lineage = ?, email = ?, hideemail = ?, dept = ?, org = ?, address = ?, country = ? WHERE userid = ?`
		_, err := db.Exec(stmt, user.Username, user.Type, user.Name.Honourific, user.Name.Family, user.Name.Given, user.Name.Lineage, user.EMail, hideEMail, user.Dept, user.Org, user.Address, user.Country, user.UserID)
		if err != nil {
			return fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		return nil
	}
	return fmt.Errorf(`bad request`)
}

func SQLReadUser(config *Config, repoID string, userid int) (*EPrintUser, error) {
	return GetUserBy(config, repoID, `userid`, userid)
}

//
// EPrint ID Lists
//

// GetAllEPrintIDs return a list of all eprint ids in repository or error
func GetAllEPrintIDs(config *Config, repoID string) ([]int, error) {
	return sqlQueryIntIDs(config, repoID, `SELECT eprintid FROM eprint
ORDER BY date_year DESC, date_month DESC, date_day DESC`)
}

// GetAllEPrintIDsWithStatus return a list of all eprint ids in a repository with a given status or return error
func GetAllEPrintIDsWithStatus(config *Config, repoID string, status string) ([]int, error) {
	return sqlQueryIntIDs(config, repoID, `SELECT eprintid FROM eprint WHERE (eprint_status = ?) ORDER BY date_year DESC, date_month DESC, date_day DESC`, status)
}

// GetEPrintIDsInTimestampRange return a list of EPrintIDs in created timestamp range
// or return error. field maybe either "datestamp" (for created date), "lastmod" (for last modified date)
func GetEPrintIDsInTimestampRange(config *Config, repoID string, field string, start string, end string) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT eprintid FROM eprint WHERE
(CONCAT(%s_year, "-",
LPAD(IFNULL(%s_month, 1), 2, "0"), "-",
LPAD(IFNULL(%s_day, 1), 2, "0"), " ",
LPAD(IFNULL(%s_hour, 0), 2, "0"), ":",
LPAD(IFNULL(%s_minute, 0), 2, "0"), ":",
LPAD(IFNULL(%s_second, 0), 2, "0")) >= ?) AND
(CONCAT(%s_year, "-",
LPAD(IFNULL(%s_month, 12), 2, "0"), "-",
LPAD(IFNULL(%s_day, 28), 2, "0"), " ",
LPAD(IFNULL(%s_hour, 23), 2, "0"), ":",
LPAD(IFNULL(%s_minute, 59), 2, "0"), ":",
LPAD(IFNULL(%s_second, 59), 2, "0")) <= ?)
ORDER BY %s_year DESC, %s_month DESC, %s_day DESC, %s_hour DESC, %s_minute DESC, %s_second DESC`,
		field, field, field, field, field, field, field, field, field, field, field, field,
		field, field, field, field, field, field)
	return sqlQueryIntIDs(config, repoID, stmt, start, end)
}

// GetEPrintIDsWithStatusInTimestampRange return a list of EPrintIDs with eprint_status in field timestamp range
// or return error. field maybe either "datestamp" (for created date), "lastmod" (for last modified date)
func GetEPrintIDsWithStatusInTimestampRange(config *Config, repoID string, status string, field string, start string, end string) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT eprintid FROM eprint WHERE
(eprint_status = ?) AND
(CONCAT(%s_year, "-",
LPAD(IFNULL(%s_month, 1), 2, "0"), "-",
LPAD(IFNULL(%s_day, 1), 2, "0"), " ",
LPAD(IFNULL(%s_hour, 0), 2, "0"), ":",
LPAD(IFNULL(%s_minute, 0), 2, "0"), ":",
LPAD(IFNULL(%s_second, 0), 2, "0")) >= ?) AND
(CONCAT(%s_year, "-",
LPAD(IFNULL(%s_month, 12), 2, "0"), "-",
LPAD(IFNULL(%s_day, 28), 2, "0"), " ",
LPAD(IFNULL(%s_hour, 23), 2, "0"), ":",
LPAD(IFNULL(%s_minute, 59), 2, "0"), ":",
LPAD(IFNULL(%s_second, 59), 2, "0")) <= ?)
ORDER BY %s_year DESC, %s_month DESC, %s_day DESC, %s_hour DESC, %s_minute DESC, %s_second DESC`,
		field, field, field, field, field, field,
		field, field, field, field, field, field,
		field, field, field, field, field, field)
	return sqlQueryIntIDs(config, repoID, stmt, status, start, end)
}

// GetEPrintIDsWithStatus returns a list of eprints in a timestmap range for
// a given status or returns an error
func GetEPrintIDsWithStatus(config *Config, repoID string, status string, start string, end string) ([]int, error) {
	stmt := `SELECT eprintid FROM eprint WHERE (eprint_status = ?) AND
(CONCAT(lastmod_year, "-",
LPAD(IFNULL(lastmod_month, 1), 2, "0"), "-",
LPAD(IFNULL(lastmod_day, 1), 2, "0"), " ",
LPAD(IFNULL(lastmod_hour, 0), 2, "0"), ":",
LPAD(IFNULL(lastmod_minute, 0), 2, "0"), ":",
LPAD(IFNULL(lastmod_second, 0), 2, "0")) >= ?) AND
(CONCAT(lastmod_year, "-",
LPAD(IFNULL(lastmod_month, 12), 2, "0"), "-",
LPAD(IFNULL(lastmod_day, 28), 2, "0"), " ",
LPAD(IFNULL(lastmod_hour, 23), 2, "0"), ":",
LPAD(IFNULL(lastmod_minute, 59), 2, "0"), ":",
LPAD(IFNULL(lastmod_second, 59), 2, "0")) <= ?)
ORDER BY lastmod_year DESC, lastmod_month DESC, lastmod_day DESC,
         lastmod_hour DESC, lastmod_minute DESC, lastmod_minute DESC`
	return sqlQueryIntIDs(config, repoID, stmt, status, start, end)
}

// GetEPrintIDsForDateType returns list of eprints in date range
// or returns an error
func GetEPrintIDsForDateType(config *Config, repoID string, dateType string, start string, end string) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT eprintid FROM eprint
WHERE ((date_type) = ?) AND
(CONCAT(date_year, "-",
LPAD(IFNULL(date_month, 1), 2, "0"), "-",
LPAD(IFNULL(date_day, 1), 2, "0")) >= ?) AND
(CONCAT(date_year, "-",
LPAD(IFNULL(date_month, 12), 2, "0"), "-",
LPAD(IFNULL(date_day, 28), 2, "0")) <= ?)
ORDER BY date_year DESC, date_month DESC, date_day DESC
`)
	return sqlQueryIntIDs(config, repoID, stmt, dateType, start, end)
}

// GetEPrintIDsWithStatusForDateType returns list of eprints in
// date range for a given status or returns an error
func GetEPrintIDsWithStatusForDateType(config *Config, repoID string, status string, dateType string, start string, end string) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT eprintid FROM eprint
WHERE (eprint_status = ? ) AND (date_type = ?) AND
(CONCAT(date_year, "-",
LPAD(IFNULL(date_month, 1), 2, "0"), "-",
LPAD(IFNULL(date_day, 1), 2, "0")) >= ?) AND
(CONCAT(date_year, "-",
LPAD(IFNULL(date_month, 12), 2, "0"), "-",
LPAD(IFNULL(date_day, 28), 2, "0")) <= ?)
ORDER BY date_year DESC, date_month DESC, date_day DESC
`)
	return sqlQueryIntIDs(config, repoID, stmt, status, dateType, start, end)
}

// GetAllUniqueID return a list of unique id values in repository
func GetAllUniqueID(config *Config, repoID string, field string) ([]string, error) {
	stmt := fmt.Sprintf(`SELECT %s
FROM eprint
WHERE %s IS NOT NULL
GROUP BY %s ORDER BY %s`,
		field, field, field, field)
	return sqlQueryStringIDs(config, repoID, stmt)
}

// GetEPrintIDsForUniqueID return list of eprints for DOI
func GetEPrintIDsForUniqueID(config *Config, repoID string, field string, value string) ([]int, error) {
	// NOTE: There should only be one eprint per DOI but we have dirty data because the field is not contrained as Unique
	stmt := fmt.Sprintf(`SELECT eprintid FROM eprint WHERE LOWER(%s) = LOWER(?)`, field)
	return sqlQueryIntIDs(config, repoID, stmt, value)
}

// GetAllPersonOrOrgIDs return a list of creator ids or error
func GetAllPersonOrOrgIDs(config *Config, repoID string, field string) ([]string, error) {
	stmt := fmt.Sprintf(`SELECT %s_id FROM eprint_%s_id
WHERE %s_id IS NOT NULL
GROUP BY %s_id ORDER BY %s_id`, field, field, field, field, field)
	return sqlQueryStringIDs(config, repoID, stmt)
}

// GetEPrintIDForPersonOrOrgID return a list of eprint ids associated with the person or organization id
func GetEPrintIDsForPersonOrOrgID(config *Config, repoID string, personOrOrgType string, personOrOrgID string) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT eprint_%s_id.eprintid AS eprintid
FROM eprint_%s_id JOIN eprint ON (eprint_%s_id.eprintid = eprint.eprintid)
WHERE eprint_%s_id.%s_id = ?
ORDER BY date_year DESC, date_month DESC, date_day DESC`,
		personOrOrgType, personOrOrgType, personOrOrgType, personOrOrgType, personOrOrgType)
	return sqlQueryIntIDs(config, repoID, stmt, personOrOrgID)
}

// GetAllORCIDs return a list of all ORCID in repository
func GetAllORCIDs(config *Config, repoID string) ([]string, error) {
	values, err := sqlQueryStringIDs(config, repoID, `SELECT creators_orcid
    FROM eprint_creators_orcid
    WHERE creators_orcid IS NOT NULL
    GROUP BY creators_orcid ORDER BY creators_orcid`)
	return values, err
}

// GetEPrintIDsForORCID return a list of eprint ids associated with the ORCID
func GetEPrintIDsForORCID(config *Config, repoID string, orcid string) ([]int, error) {
	return sqlQueryIntIDs(config, repoID, `SELECT eprint.eprintid AS eprintid
FROM eprint_creators_orcid JOIN eprint ON (eprint_creators_orcid.eprintid = eprint.eprintid)
WHERE creators_orcid = ?
ORDER BY date_year DESC, date_month DESC, date_day DESC
`, orcid)
}

// GetAllItems returns a list of simple items (e.g. local_group)
func GetAllItems(config *Config, repoID string, field string) ([]string, error) {
	stmt := fmt.Sprintf(`SELECT %s
FROM eprint_%s
WHERE eprint_%s.%s IS NOT NULL
GROUP BY eprint_%s.%s ORDER BY eprint_%s.%s`,
		field, field, field, field, field, field, field, field)
	return sqlQueryStringIDs(config, repoID, stmt)
}

// GetEPrintIDsForItem
func GetEPrintIDsForItem(config *Config, repoID string, field string, value string) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT eprint.eprintid AS eprintid
FROM eprint_%s JOIN eprint ON (eprint_%s.eprintid = eprint.eprintid)
WHERE eprint_%s.%s = ?
ORDER BY eprint.date_year DESC, eprint.date_month DESC, eprint.date_day DESC`, field, field, field, field)
	return sqlQueryIntIDs(config, repoID, stmt, value)
}

// GetAllPersonNames return a list of person names in repository
func GetAllPersonNames(config *Config, repoID string, field string) ([]string, error) {
	stmt := fmt.Sprintf(`SELECT CONCAT(%s_family, "/", %s_given) AS %s
FROM eprint_%s
WHERE (%s_family IS NOT NULL) OR (%s_given IS NOT NULL)
GROUP BY %s_family, %s_given ORDER BY %s_family, %s_given`,
		field, field, field,
		field, field, field, field, field, field, field)
	return sqlQueryStringIDs(config, repoID, stmt)
}

// GetEPrintIDsForPersonName return a list of eprint id for a person's name (family, given)
func GetEPrintIDsForPersonName(config *Config, repoID, field string, family string, given string) ([]int, error) {
	conditions := []string{}
	if strings.Contains(family, "*") || strings.Contains(given, "%") {
		conditions = append(conditions, fmt.Sprintf(`%s_family LIKE ?`, field))
	} else if family != "" {
		conditions = append(conditions, fmt.Sprintf(`%s_family = ?`, field))
	}
	if strings.Contains(given, "*") || strings.Contains(given, "%") {
		conditions = append(conditions, fmt.Sprintf(`%s_given LIKE ?`, field))
	} else if given != "" {
		conditions = append(conditions, fmt.Sprintf(`%s_given = ?`, field))
	}
	stmt := fmt.Sprintf(`SELECT eprint.eprintid AS eprintid
FROM eprint_%s JOIN eprint ON (eprint_%s.eprintid = eprint.eprintid)
WHERE %s
ORDER BY %s_family ASC, %s_given ASC, eprint.date_year DESC, eprint.date_month DESC, eprint.date_day DESC`,
		field, field, strings.Join(conditions, " AND "), field, field)
	return sqlQueryIntIDs(config, repoID, stmt, family, given)
}

// GetAllYears returns the publication years found in a repository
func GetAllYears(config *Config, repoID string) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT date_year FROM eprint WHERE date_type = "published" AND date_year IS NOT NULL GROUP BY date_year ORDER BY date_year DESC`)
	return sqlQueryInts(config, repoID, stmt)
}

// GetEPrintsIDsForYear returns a list of published eprint IDs for a given
// year.
func GetEPrintIDsForYear(config *Config, repoID string, year int) ([]int, error) {
	stmt := fmt.Sprintf(`SELECT eprintid FROM eprint WHERE date_type = "published" AND date_year  = ? ORDER BY date_year DESC, date_month DESC, date_day DESC`)
	return sqlQueryIntIDs(config, repoID, stmt, year)
}

//
// EPrints Metadata Structure
//

// eprintTablesAndColumns takes a DB connection and repoID then builds a map[string][]string{}
// structure representing the tables and their columns available in a EPrints Repository
func eprintTablesAndColumns(db *sql.DB, repoID string) (map[string][]string, error) {
	data := map[string][]string{}
	stmt := `SHOW TABLES LIKE "eprint%"`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
	}
	tables := []string{}
	for rows.Next() {
		tableName := ""
		if err := rows.Scan(&tableName); err == nil {
			if !strings.Contains(tableName, "__") {
				tables = append(tables, tableName)
			}
		}
	}
	rows.Close()

	for _, tableName := range tables {
		data[tableName] = []string{}
		stmt := fmt.Sprintf(`SHOW COLUMNS IN %s`, tableName)
		cRows, err := db.Query(stmt)
		if err != nil {
			return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
		}
		columns := []string{}
		var (
			colName, f1, f2, f3, f5 string
			f4                      interface{}
		)
		for cRows.Next() {
			//colName, f1, f2, f3, f4, f5 = &"", &"", &"", &"", nil, &""
			if err := cRows.Scan(&colName, &f1, &f2, &f3, &f4, &f5); err != nil {
				return nil, fmt.Errorf("cRows.Scan() error: %s", err)
			} else {
				columns = append(columns, colName)
			}
		}
		data[tableName] = columns
		cRows.Close()
	}
	// We need to add the document set of tables too.
	stmt = `SHOW TABLES LIKE "document%"`
	rows, err = db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
	}
	tables = []string{}
	for rows.Next() {
		tableName := ""
		if err := rows.Scan(&tableName); err == nil {
			if !strings.Contains(tableName, "__") {
				tables = append(tables, tableName)
			}
		}
	}
	rows.Close()

	for _, tableName := range tables {
		data[tableName] = []string{}
		stmt := fmt.Sprintf(`SHOW COLUMNS IN %s`, tableName)
		cRows, err := db.Query(stmt)
		if err != nil {
			return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
		}
		columns := []string{}
		var (
			colName, f1, f2, f3, f5 string
			f4                      interface{}
		)
		for cRows.Next() {
			//colName, f1, f2, f3, f4, f5 = &"", &"", &"", &"", nil, &""
			if err := cRows.Scan(&colName, &f1, &f2, &f3, &f4, &f5); err != nil {
				return nil, fmt.Errorf("cRows.Scan() error: %s", err)
			} else {
				columns = append(columns, colName)
			}
		}
		data[tableName] = columns
		cRows.Close()
	}
	// We need to add the files set of tables too.
	stmt = `SHOW TABLES LIKE "file%"`
	rows, err = db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
	}
	tables = []string{}
	for rows.Next() {
		tableName := ""
		if err := rows.Scan(&tableName); err == nil {
			if !strings.Contains(tableName, "__") {
				tables = append(tables, tableName)
			}
		}
	}
	rows.Close()

	for _, tableName := range tables {
		data[tableName] = []string{}
		stmt := fmt.Sprintf(`SHOW COLUMNS IN %s`, tableName)
		cRows, err := db.Query(stmt)
		if err != nil {
			return nil, fmt.Errorf("SQL(%q), %s", repoID, err)
		}
		columns := []string{}
		var (
			colName, f1, f2, f3, f5 string
			f4                      interface{}
		)
		for cRows.Next() {
			//colName, f1, f2, f3, f4, f5 = &"", &"", &"", &"", nil, &""
			if err := cRows.Scan(&colName, &f1, &f2, &f3, &f4, &f5); err != nil {
				return nil, fmt.Errorf("cRows.Scan() error: %s", err)
			} else {
				columns = append(columns, colName)
			}
		}
		data[tableName] = columns
		cRows.Close()
	}
	return data, nil
}

/*
 * Column mapping for tables.
 */

// colExpr takes a column name, ifNull bool and default value.
// If the "ifNull" bool is true then the form expressed is
// `IFNULL(%s, %s) AS %s` otherwise just the column name
// is returned.
func colExpr(name string, ifNull bool, value string) string {
	if ifNull {
		return fmt.Sprintf(`IFNULL(%s, %s) AS %s`, name, value, name)
	}
	return name
}

// eprintToColumnsAndValues for a given EPrints struct generate a
// list of column names to query along with a recieving values array.
// Return a list of column names (with null handle and aliases) and values.
//
// The bool ifNull will control the type of expression of the column.
func eprintToColumnsAndValues(eprint *EPrint, columnsIn []string, ifNull bool) ([]string, []interface{}) {
	columnsOut := []string{}
	values := []interface{}{}
	for i, key := range columnsIn {
		switch key {
		case "eprintid":
			values = append(values, &eprint.EPrintID)
			columnsOut = append(columnsOut, key)
		case "rev_number":
			values = append(values, &eprint.RevNumber)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "eprint_status":
			values = append(values, &eprint.EPrintStatus)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "userid":
			values = append(values, &eprint.UserID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "dir":
			values = append(values, &eprint.Dir)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "datestamp_year":
			values = append(values, &eprint.DatestampYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_month":
			values = append(values, &eprint.DatestampMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_day":
			values = append(values, &eprint.DatestampDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_hour":
			values = append(values, &eprint.DatestampHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_minute":
			values = append(values, &eprint.DatestampMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datestamp_second":
			values = append(values, &eprint.DatestampSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_year":
			values = append(values, &eprint.LastModifiedYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_month":
			values = append(values, &eprint.LastModifiedMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_day":
			values = append(values, &eprint.LastModifiedDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_hour":
			values = append(values, &eprint.LastModifiedHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_minute":
			values = append(values, &eprint.LastModifiedMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "lastmod_second":
			values = append(values, &eprint.LastModifiedSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_year":
			values = append(values, &eprint.StatusChangedYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_month":
			values = append(values, &eprint.StatusChangedMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_day":
			values = append(values, &eprint.StatusChangedDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_hour":
			values = append(values, &eprint.StatusChangedHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_minute":
			values = append(values, &eprint.StatusChangedMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "status_changed_second":
			values = append(values, &eprint.StatusChangedSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "type":
			values = append(values, &eprint.Type)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "metadata_visibility":
			values = append(values, &eprint.MetadataVisibility)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "title":
			values = append(values, &eprint.Title)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "ispublished":
			values = append(values, &eprint.IsPublished)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "full_text_status":
			values = append(values, &eprint.FullTextStatus)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "keywords":
			values = append(values, &eprint.Keywords)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "note":
			values = append(values, &eprint.Note)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "abstract":
			values = append(values, &eprint.Abstract)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "date_year":
			values = append(values, &eprint.DateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_month":
			values = append(values, &eprint.DateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_day":
			values = append(values, &eprint.DateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_type":
			values = append(values, &eprint.DateType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "series":
			values = append(values, &eprint.Series)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "volume":
			values = append(values, &eprint.Volume)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "number":
			values = append(values, &eprint.Number)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "publication":
			values = append(values, &eprint.Publication)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "publisher":
			values = append(values, &eprint.Publisher)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "place_of_pub":
			values = append(values, &eprint.PlaceOfPub)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "edition":
			values = append(values, &eprint.Edition)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pagerange":
			values = append(values, &eprint.PageRange)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pages":
			values = append(values, &eprint.Pages)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "event_type":
			values = append(values, &eprint.EventType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "event_title":
			values = append(values, &eprint.EventTitle)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "event_location":
			values = append(values, &eprint.EventLocation)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "event_dates":
			values = append(values, &eprint.EventDates)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "id_number":
			values = append(values, &eprint.IDNumber)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "refereed":
			values = append(values, &eprint.Refereed)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "isbn":
			values = append(values, &eprint.ISBN)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "issn":
			values = append(values, &eprint.ISSN)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "book_title":
			values = append(values, &eprint.BookTitle)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "official_url":
			values = append(values, &eprint.OfficialURL)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "alt_url":
			values = append(values, &eprint.AltURL)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "rights":
			values = append(values, &eprint.Rights)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "collection":
			values = append(values, &eprint.Collection)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "reviewer":
			values = append(values, &eprint.Reviewer)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "official_cit":
			values = append(values, &eprint.OfficialCitation)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "monograph_type":
			values = append(values, &eprint.MonographType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "suggestions":
			values = append(values, &eprint.Suggestions)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pres_type":
			values = append(values, &eprint.PresType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "succeeds":
			values = append(values, &eprint.Succeeds)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "commentary":
			values = append(values, &eprint.Commentary)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "contact_email":
			values = append(values, &eprint.ContactEMail)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "fileinfo":
			values = append(values, &eprint.FileInfo)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "latitude":
			values = append(values, &eprint.Latitude)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0.0`))
		case "longitude":
			values = append(values, &eprint.Longitude)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0.0`))
		case "department":
			values = append(values, &eprint.Department)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "output_media":
			values = append(values, &eprint.OutputMedia)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "num_pieces":
			values = append(values, &eprint.NumPieces)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "composition_type":
			values = append(values, &eprint.CompositionType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "data_type":
			values = append(values, &eprint.DataType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pedagogic_type":
			values = append(values, new(string))
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "learning_level":
			values = append(values, &eprint.LearningLevelText)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "completion_time":
			values = append(values, &eprint.CompletionTime)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "task_purpose":
			values = append(values, &eprint.TaskPurpose)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "doi":
			values = append(values, &eprint.DOI)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pmc_id":
			values = append(values, &eprint.PMCID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "pmid":
			values = append(values, &eprint.PMID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "parent_url":
			values = append(values, &eprint.ParentURL)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "toc":
			values = append(values, &eprint.TOC)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "interviewer":
			values = append(values, &eprint.Interviewer)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "interviewdate":
			values = append(values, &eprint.InterviewDate)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "nonsubj_keywords":
			values = append(values, &eprint.NonSubjKeywords)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "season":
			values = append(values, &eprint.Season)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "classification_code":
			values = append(values, &eprint.ClassificationCode)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "sword_depositor":
			values = append(values, &eprint.SwordDepositor)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "sword_depository":
			values = append(values, &eprint.SwordDepository)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "sword_slug":
			values = append(values, &eprint.SwordSlug)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "importid":
			values = append(values, &eprint.ImportID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "patent_applicant":
			values = append(values, &eprint.PatentApplicant)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "patent_number":
			values = append(values, &eprint.PatentNumber)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "institution":
			values = append(values, &eprint.Institution)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_type":
			values = append(values, &eprint.ThesisType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_degree":
			values = append(values, &eprint.ThesisDegree)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_degree_grantor":
			values = append(values, &eprint.ThesisDegreeGrantor)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_degree_date_year":
			values = append(values, &eprint.ThesisDegreeDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_degree_date_month":
			values = append(values, &eprint.ThesisDegreeDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_degree_date_day":
			values = append(values, &eprint.ThesisDegreeDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_submitted_date_year":
			values = append(values, &eprint.ThesisSubmittedDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_submitted_date_month":
			values = append(values, &eprint.ThesisSubmittedDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_submitted_date_day":
			values = append(values, &eprint.ThesisSubmittedDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_defense_date":
			values = append(values, &eprint.ThesisDefenseDate)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "thesis_defense_date_year":
			values = append(values, &eprint.ThesisDefenseDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_defense_date_month":
			values = append(values, &eprint.ThesisDefenseDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_defense_date_day":
			values = append(values, &eprint.ThesisDefenseDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_approved_date_year":
			values = append(values, &eprint.ThesisApprovedDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_approved_date_month":
			values = append(values, &eprint.ThesisApprovedDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_approved_date_day":
			values = append(values, &eprint.ThesisApprovedDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_public_date_year":
			values = append(values, &eprint.ThesisPublicDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_public_date_month":
			values = append(values, &eprint.ThesisPublicDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_public_date_day":
			values = append(values, &eprint.ThesisPublicDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_author_email":
			values = append(values, &eprint.ThesisAuthorEMail)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "hide_thesis_author_email":
			values = append(values, &eprint.HideThesisAuthorEMail)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "gradofc_approval_date":
			values = append(values, &eprint.GradOfficeApprovalDate)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "gradofc_approval_date_year":
			values = append(values, &eprint.GradOfficeApprovalDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "gradofc_approval_date_month":
			values = append(values, &eprint.GradOfficeApprovalDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "gradofc_approval_date_day":
			values = append(values, &eprint.GradOfficeApprovalDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "thesis_awards":
			values = append(values, &eprint.ThesisAwards)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "review_status":
			values = append(values, &eprint.ReviewStatus)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "copyright_statement":
			values = append(values, &eprint.CopyrightStatement)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "source":
			values = append(values, &eprint.Source)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "replacedby":
			values = append(values, &eprint.ReplacedBy)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "item_issues_count":
			values = append(values, &eprint.ItemIssuesCount)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "errata":
			values = append(values, &eprint.ErrataText)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "coverage_dates":
			values = append(values, &eprint.CoverageDates)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "edit_lock_user":
			values = append(values, &eprint.EditLockUser)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "edit_lock_since":
			values = append(values, &eprint.EditLockSince)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "edit_lock_until":
			values = append(values, &eprint.EditLockUntil)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		// The follow values represent sub tables and processed separately.
		case "patent_classification":
			values = append(values, &eprint.PatentClassificationText)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "language":
			values = append(values, &eprint.Language)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "referencetext":
			// NOTE: This is just an ignore column, CaltechES has this as a VARCHAR column in the eprint
			// table and all values are NULL. None of our other repositories have this columns. Possibly
			// added in error or is a legacy EPrint column that was migrated from an older version but was not used.
		default:
			// Handle case where we have value that is unmapped or not available in EPrint struct
			log.Printf("could not map %q (col. %d, eprintid %d) into EPrint struct", key, i, eprint.EPrintID)
		}

	}
	return columnsOut, values
}

func documentToColumnsAndValues(eprintID int, document *Document, columns []string, ifNull bool) ([]string, []interface{}) {
	columnsOut := []string{}
	values := []interface{}{}
	for i, key := range columns {
		switch key {
		case "docid":
			values = append(values, &document.DocID)
			columnsOut = append(columnsOut, key)
		case "eprintid":
			values = append(values, &document.EPrintID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "pos":
			values = append(values, &document.Pos)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "rev_number":
			values = append(values, &document.RevNumber)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "format":
			values = append(values, &document.Format)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "formatdesc":
			values = append(values, &document.FormatDesc)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "language":
			values = append(values, &document.Language)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "security":
			values = append(values, &document.Security)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "license":
			values = append(values, &document.License)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "main":
			values = append(values, &document.Main)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "date_embargo_year":
			values = append(values, &document.DateEmbargoYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_embargo_month":
			values = append(values, &document.DateEmbargoMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "date_embargo_day":
			values = append(values, &document.DateEmbargoDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "content":
			values = append(values, &document.Content)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "placement":
			values = append(values, &document.Placement)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mime_type":
			values = append(values, &document.MimeType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "media_duration":
			values = append(values, &document.MediaDuration)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_audio_codec":
			values = append(values, &document.MediaAudioCodec)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "media_video_codec":
			values = append(values, &document.MediaVideoCodec)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "media_width":
			values = append(values, &document.MediaWidth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_height":
			values = append(values, &document.MediaHeight)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_aspect_ratio":
			values = append(values, &document.MediaAspectRatio)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_sample_start":
			values = append(values, &document.MediaSampleStart)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "media_sample_stop":
			values = append(values, &document.MediaSampleStop)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		default:
			log.Printf("%q (col. %d, eprintid %d) not found in document table", key, i, eprintID)
		}
	}
	return columnsOut, values
}

func fileToColumnsAndValues(eprintID int, file *File, columns []string, ifNull bool) ([]string, []interface{}) {
	columnsOut := []string{}
	values := []interface{}{}
	for i, key := range columns {
		switch key {
		case "fileid":
			values = append(values, &file.FileID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "datasetid":
			values = append(values, &file.DatasetID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "objectid":
			values = append(values, &file.ObjectID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "filename":
			values = append(values, &file.Filename)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "mime_type":
			values = append(values, &file.MimeType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "hash":
			values = append(values, &file.Hash)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "hash_type":
			values = append(values, &file.HashType)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `""`))
		case "filesize":
			values = append(values, &file.FileSize)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_year":
			values = append(values, &file.MTimeYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_month":
			values = append(values, &file.MTimeMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_day":
			values = append(values, &file.MTimeDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_hour":
			values = append(values, &file.MTimeHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_minute":
			values = append(values, &file.MTimeMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "mtime_second":
			values = append(values, &file.MTimeSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "pronomid":
			values = append(values, &file.PronomID)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "classification_date_year":
			values = append(values, &file.ClassificationDateYear)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "classification_date_month":
			values = append(values, &file.ClassificationDateMonth)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "classification_date_day":
			values = append(values, &file.ClassificationDateDay)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "classification_date_hour":
			values = append(values, &file.ClassificationDateHour)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "classification_date_minute":
			values = append(values, &file.ClassificationDateMinute)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "classification_date_second":
			values = append(values, &file.ClassificationDateSecond)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		case "classification_quality":
			values = append(values, &file.ClassificationQuality)
			columnsOut = append(columnsOut, colExpr(key, ifNull, `0`))
		default:
			log.Printf("(col. %d, eprintid %d) field %q not found in file table", i, eprintID, key)
		}
	}
	return columnsOut, values
}

/*
 * Document and files models
 */
func documentIDToFiles(repoID string, baseURL string, eprintID int, documentID int, pos int, db *sql.DB, tables map[string][]string) []*File {
	//FIXME: Need to figure out if I need to pay attention to
	// file_copies_plugin and file_copies_sourceid tables. This appear
	// to be related to the storage manager.  They don't appear to
	// be reflected in the EPrint XML. files_copies_plugin has a single
	// unique column called copies_pluginid which has a single value
	// of "Storage::Local" for all rows (about the same number of rows
	// as in the file table.  In the table file_copies_sourceid the column
	// copies_sourceid appears to be base filenames like in the
	// file.filename where file.datasetid = "document". I'm ignoring in
	// this processing for file table for now.
	tableName := "file"
	columns, ok := tables[tableName]
	if ok {
		files := []*File{}
		file := new(File)
		columnSQL, _ := fileToColumnsAndValues(eprintID, file, columns, true)
		//FIXME: This needs to be an ordered list.
		stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE datasetid = 'document' AND objectid = ?`, strings.Join(columnSQL, ", "), tableName)
		rows, err := db.Query(stmt, documentID)
		if err != nil {
			log.Printf("Query failed %q for %d in %q, doc ID %d, , %q,  %s", tableName, eprintID, repoID, documentID, stmt, err)
		} else {
			i := 0
			for rows.Next() {
				file = new(File)
				_, values := fileToColumnsAndValues(eprintID, file, columns, true)
				if err := rows.Scan(values...); err != nil {
					log.Printf("Could not scan %q for %d in %q, doc ID %d, %s", tableName, eprintID, repoID, documentID, err)
				} else {
					file.ID = fmt.Sprintf("%s/id/file/%d", baseURL, file.FileID)
					file.MTime = makeTimestamp(file.MTimeYear, file.MTimeMonth, file.MTimeDay, file.MTimeHour, file.MTimeMinute, file.MTimeSecond)
					file.URL = fmt.Sprintf("%s/%d/%d/%s", baseURL, eprintID, pos, file.Filename)
					files = append(files, file)
				}
				i++
			}
		}
		if len(files) > 0 {
			return files
		}
	}
	return nil
}

func documentIDToRelation(repoID string, baseURL string, documentID int, pos int, db *sql.DB, tables map[string][]string) *RelationItemList {
	typeTable := "document_relation_type"
	_, okTypeTable := tables[typeTable]
	uriTable := "document_relation_uri"
	_, okUriTable := tables[uriTable]

	if okTypeTable && okUriTable {
		itemList := new(RelationItemList)
		stmt := fmt.Sprintf(`SELECT document_relation_type.pos AS pos, document_relation_type.relation_type, document_relation_uri.relation_uri FROM %s JOIN %s ON ((%s.docid = %s.docid) AND (%s.pos = %s.pos)) WHERE (%s.docid = ?)`, typeTable, uriTable, typeTable, uriTable, typeTable, uriTable, typeTable)
		rows, err := db.Query(stmt, documentID)
		if err != nil {
			log.Printf("Query failed %q, doc id %d, pos %d, %s", stmt, documentID, pos, err)
		} else {
			i, pos := 0, 0
			for rows.Next() {
				var (
					relationType, relationURI string
				)
				if err := rows.Scan(&pos, &relationType, &relationURI); err != nil {
					log.Printf("Could not scan relation type and relation uri (%d), %q join %q, doc id %d and pos %d, %s", pos, typeTable, uriTable, documentID, pos, err)
				} else {
					resizeItemList(pos, itemList)
					item := itemList.IndexOf(pos)
					item.Pos = pos
					item.Type = relationType
					item.URI = fmt.Sprintf(`%s%s`, baseURL, relationURI)
				}
				i++
			}
			if itemList.Length() > 0 {
				return itemList
			}
		}
	}
	return nil
}

func userIDToName(repoID string, userID int, db *sql.DB) string {
	name := ""
	if userID > 0 {
		stmt := `SELECT TRIM(CONCAT_WS(' ', IFNULL(name_honourific, ''), IFNULL(name_given, ''), IFNULL(name_family, ''), IFNULL(name_lineage, ''))) AS name FROM user WHERE userid = ?`

		rows, err := db.Query(stmt, userID)
		if err != nil {
			log.Printf(`Query failed "user" for %d in %q, %q,  %s`, userID, repoID, stmt, err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&name); err == nil {
				return name
			}
		}
	}
	return ""
}

func eprintIDToDocumentList(repoID string, baseURL string, eprintID int, db *sql.DB, tables map[string][]string) *DocumentList {
	tableName := "document"
	columns, ok := tables[tableName]
	if ok {
		documentList := new(DocumentList)
		document := new(Document)
		// NOTE: Bind the values in document to the values array used by
		// rows.Scan().
		columnSQL, _ := documentToColumnsAndValues(eprintID, document, columns, true)
		stmt := fmt.Sprintf(`SELECT %s FROM %s WHERE eprintid = ? ORDER BY eprintid ASC, pos ASC, rev_number DESC`, strings.Join(columnSQL, ", "), tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Query failed %q for %d in %q, %q,  %s", tableName, eprintID, repoID, stmt, err)
		} else {
			i := 0
			for rows.Next() {
				document = new(Document)
				_, values := documentToColumnsAndValues(eprintID, document, columns, true)
				if err := rows.Scan(values...); err != nil {
					log.Printf("Could not scan %q for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					document.ID = fmt.Sprintf("%s/id/document/%d", baseURL, document.DocID)
					document.Files = documentIDToFiles(repoID, baseURL, eprintID, document.DocID, document.Pos, db, tables)
					document.Relation = documentIDToRelation(repoID, baseURL, document.DocID, document.Pos, db, tables)
					documentList.Append(document)
				}
				i++
			}
			rows.Close()
		}

		// NOTE: The document_permission_group the table is empty in our repositories
		if documentList.Length() > 0 {
			// Attach files to documents
			for i := 0; i < documentList.Length(); i++ {
				document := documentList.IndexOf(i)
				files := documentIDToFiles(repoID, baseURL, eprintID, document.DocID, document.Pos, db, tables)
				if (files != nil) && (len(files) > 0) {
					document.Files = files
				}
			}
			return documentList
		}
	}
	return nil
}

/*
 * Common models and help functions
 */

func makePersonName(given string, family string, honourific string, lineage string) *Name {
	name := new(Name)
	isFlat := true
	if s := strings.TrimSpace(given); s != "" {
		name.Given = s
		isFlat = false
	}
	if s := strings.TrimSpace(family); s != "" {
		name.Family = s
		isFlat = false
	}
	if s := strings.TrimSpace(honourific); s != "" {
		name.Honourific = s
		isFlat = false
	}
	if s := strings.TrimSpace(lineage); s != "" {
		name.Lineage = s
		isFlat = false
	}
	if isFlat {
		return nil
	}
	return name
}

// NOTE: Do to the funky way MySQL and EPrints works with UTF
// I can't use the time package to build my formatted date string.
// An odd edge case can make the year, month or day off by one.

func makeTimestamp(year int, month int, day int, hour int, minute int, second int) string {
	if year > 0 && (hour > 0 || minute > 0 || second > 0) {
		return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
	}
	// Oral Historied repositories has a "datestamp" without hour, minute, second
	if year > 0 {
		return fmt.Sprintf("%d-%02d-%02d", year, month, day)
	}
	return ""
}

func makeDatestamp(year int, month int, day int) string {
	parts := []string{}
	if year > 0 {
		parts = append(parts, fmt.Sprintf("%d", year))
	}
	if month > 0 {
		parts = append(parts, fmt.Sprintf("%02d", month))
	}
	if day > 0 {
		parts = append(parts, fmt.Sprintf("%02d", day))
	}
	if len(parts) == 3 {
		return strings.Join(parts, "-")
	}
	return ""
}

func makeApproxDate(year int, month int, day int) string {
	parts := []string{}
	if year > 0 {
		parts = append(parts, fmt.Sprintf("%d", year))
	}
	if month > 0 {
		parts = append(parts, fmt.Sprintf("%02d", month))
	}
	if day > 0 {
		parts = append(parts, fmt.Sprintf("%02d", day))
	}
	if len(parts) > 0 {
		return strings.Join(parts, "-")
	}
	return ""
}

func approxYMD(src string) (int, int, int) {
	layout := "2006-01-02"
	switch len(src) {
	case 4:
		layout = "2006"
	case 7:
		layout = "2006-01"
	case 10:
		layout = "2006-01-02"
	}
	dt, err := time.Parse(layout, src)
	if err == nil {
		switch len(src) {
		case 4:
			return dt.Year(), 0, 0
		case 7:
			return dt.Year(), int(dt.Month()), 0
		case 10:
			return dt.Year(), int(dt.Month()), dt.Day()
		}
	}
	return 0, 0, 0
}

func makeDirValue(ID int) string {
	return fmt.Sprintf(`disk0/%s`, strings.TrimSuffix(pairtree.Encode(fmt.Sprintf("%08d", ID)), "/"))
}

func resizeItemList(pos int, itemList ItemsInterface) {
	if (pos >= 0) && (pos >= itemList.Length()) {
		for j := itemList.Length(); j <= (pos + 1); j++ {
			item := new(Item)
			itemList.Append(item)
		}
	}
}

/*
 * PersonItemList model
 */
func eprintIDToPersonItemList(db *sql.DB, tables map[string][]string, repoID string, eprintID int, tablePrefix string, itemList ItemsInterface) int {
	var (
		pos                                       int
		value, honourific, given, family, lineage string
	)
	tableName := tablePrefix + `_name`
	_, ok := tables[tableName]
	if ok {
		columnPrefix := strings.TrimPrefix(tableName, `eprint_`)
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s_honourific, '') AS honourific, IFNULL(%s_given, '') AS given, IFNULL(%s_family, '') AS family, IFNULL(%s_lineage, '') AS lineage FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnPrefix, columnPrefix, columnPrefix, columnPrefix, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			for rows.Next() {
				if err := rows.Scan(&pos, &honourific, &given, &family, &lineage); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					// Check if we have enough items in our item list.
					resizeItemList(pos, itemList)
					if item := itemList.IndexOf(pos); item != nil {
						item.Pos = pos
						item.Name = makePersonName(given, family, honourific, lineage)
					} else {
						log.Printf("itemList too short, pos (%d) not found for %d in %s", pos, eprintID, tableName)
					}
				}
			}
			rows.Close()
		}
	}
	tablesAndColumn := map[string][2]string{}
	columnPrefix := strings.TrimPrefix(tablePrefix, `eprint_`)
	for _, suffix := range []string{"id", "orcid", "uri", "url", "role", "email", "show_email"} {
		key := fmt.Sprintf("%s_%s", tablePrefix, suffix)
		tablesAndColumn[key] = [2]string{
			// Column Name
			fmt.Sprintf("%s_%s", columnPrefix, suffix),
			// Column Alias
			suffix,
		}
	}
	i := 0
	for tableName, columnNames := range tablesAndColumn {
		columnName, columnAlias := columnNames[0], columnNames[1]
		_, ok := tables[tableName]
		if ok {
			stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnAlias, tableName)
			rows, err := db.Query(stmt, eprintID)
			if err != nil {
				log.Printf("Could not query (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
			} else {
				for rows.Next() {
					if err := rows.Scan(&pos, &value); err != nil {
						log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
					} else {
						itemList.SetAttributeOf(pos, columnAlias, value)
					}
				}
				rows.Close()
			}
		}
	}
	return itemList.Length()
}

func eprintIDToCreators(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CreatorItemList {
	tablePrefix := `eprint_creators`
	itemList := new(CreatorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToEditors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *EditorItemList {
	tablePrefix := `eprint_editors`
	itemList := new(EditorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToContributors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ContributorItemList {
	tablePrefix := `eprint_contributors`
	itemList := new(ContributorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToExhibitors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ExhibitorItemList {
	tablePrefix := `eprint_exhibitors`
	itemList := new(ExhibitorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToProducers(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ProducerItemList {
	tablePrefix := `eprint_producers`
	itemList := new(ProducerItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToConductors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ConductorItemList {
	tablePrefix := `eprint_conductors`
	itemList := new(ConductorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToLyricists(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *LyricistItemList {
	tablePrefix := `eprint_lyricists`
	itemList := new(LyricistItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToThesisAdvisors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ThesisAdvisorItemList {
	tablePrefix := `eprint_thesis_advisor`
	itemList := new(ThesisAdvisorItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToThesisCommittee(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ThesisCommitteeItemList {
	tablePrefix := `eprint_thesis_committee`
	itemList := new(ThesisCommitteeItemList)
	if count := eprintIDToPersonItemList(db, tables, repoID, eprintID, tablePrefix, itemList); count > 0 {
		return itemList
	}
	return nil
}

/*
 * SimpleItemList model
 */

func eprintIDToSimpleItemList(db *sql.DB, tables map[string][]string, repoID string, eprintID int, tableName string, itemList ItemsInterface) int {
	columnName := strings.TrimPrefix(tableName, `eprint_`)
	var (
		pos   int
		value string
	)
	_, ok := tables[tableName]
	if ok {
		stmt := fmt.Sprintf(`SELECT pos, TRIM(IFNULL(%s, '')) AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					if value != "" {
						resizeItemList(pos, itemList)
						if item := itemList.IndexOf(pos); item != nil {
							item.Pos = pos
							item.Value = value
						} else {

							log.Printf("Failed to resize item major/minor list (table: %s, eprintID: %d) scanned pos: %d, value: %q itemList.Length() %d -> %+v\n", tableName, eprintID, pos, value, itemList.Length(), item)
						}
					}
				}
				i++
			}
			rows.Close()
		}
	}
	return itemList.Length()
}

func eprintIDToLocalGroup(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *LocalGroupItemList {
	tableName := `eprint_local_group`
	itemList := new(LocalGroupItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToReferenceText(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ReferenceTextItemList {
	tableName := `eprint_referencetext`
	itemList := new(ReferenceTextItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToProjects(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ProjectItemList {
	tableName := `eprint_projects`
	itemList := new(ProjectItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToSubjects(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *SubjectItemList {
	tableName := `eprint_subjects`
	itemList := new(SubjectItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToAccompaniment(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *AccompanimentItemList {
	tableName := `eprint_accompaniment`
	itemList := new(AccompanimentItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToSkillAreas(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *SkillAreaItemList {
	tableName := `eprint_skill_areas`
	itemList := new(SkillAreaItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToCopyrightHolders(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CopyrightHolderItemList {
	tableName := `eprint_copyright_holders`
	itemList := new(CopyrightHolderItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToReference(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ReferenceItemList {
	tableName := `eprint_reference`
	itemList := new(ReferenceItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToAltTitle(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *AltTitleItemList {
	tableName := `eprint_alt_title`
	itemList := new(AltTitleItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToPatentAssignee(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *PatentAssigneeItemList {
	tableName := `eprint_patent_assignee`
	itemList := new(PatentAssigneeItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToRelatedPatents(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *RelatedPatentItemList {
	tableName := `eprint_related_patents`
	itemList := new(RelatedPatentItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToDivisions(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *DivisionItemList {
	tableName := `eprint_divisions`
	itemList := new(DivisionItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToOptionMajor(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *OptionMajorItemList {
	tableName := `eprint_option_major`
	itemList := new(OptionMajorItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

func eprintIDToOptionMinor(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *OptionMinorItemList {
	tableName := `eprint_option_minor`
	itemList := new(OptionMinorItemList)
	if count := eprintIDToSimpleItemList(db, tables, repoID, eprintID, tableName, itemList); count > 0 {
		return itemList
	}
	return nil
}

/*
 * Hetrogenous models
 */

func eprintIDToConfCreators(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ConfCreatorItemList {
	var (
		pos   int
		value string
	)
	tableName := `eprint_conf_creators_name`
	columnName := `conf_creators_name`
	_, ok := tables[tableName]
	if ok {
		itemList := new(ConfCreatorItemList)
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s
FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					resizeItemList(pos, itemList)
					item := itemList.IndexOf(pos)
					item.Pos = pos
					item.Name = new(Name)
					item.Name.Value = strings.TrimSpace(value)
				}
				i++
			}
			rows.Close()

			if itemList.Length() > 0 {
				tablesAndColumn := map[string]string{
					"eprint_conf_creators_id":  "conf_creators_id",
					"eprint_conf_creators_ror": "conf_creators_ror",
					"eprint_conf_creators_uri": "conf_creators_uri",
					"eprint_conf_creators":     "conf_creators",
				}
				for tableName, columnName := range tablesAndColumn {
					_, ok := tables[tableName]
					if ok {
						stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
						rows, err = db.Query(stmt, eprintID)
						if err != nil {
							log.Printf("Could not query (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
						} else {
							i := 0
							for rows.Next() {
								if err := rows.Scan(&pos, &value); err != nil {
									log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
								} else {
									for _, item := range itemList.Items {
										if item.Pos == pos && value != "" {
											switch columnName {
											case "conf_creators_id":
												item.ID = value
											case "conf_creators_ror":
												item.ROR = value
											case "conf_creators_uri":
												item.URI = value
											case "conf_creators":
												item.Name = new(Name)
												item.Name.Value = strings.TrimSpace(value)
											}
											break
										}
									}
								}
								i++
							}
							rows.Close()
						}
					}
					return itemList
				}
			}
		}
	}
	return nil
}

func eprintIDToCorpCreators(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CorpCreatorItemList {
	var (
		pos   int
		value string
	)
	itemList := new(CorpCreatorItemList)
	tableName := `eprint_corp_creators_name`
	columnName := `corp_creators_name`
	_, ok := tables[tableName]
	if ok {
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s
FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					resizeItemList(pos, itemList)
					item := itemList.IndexOf(pos)
					item.Pos = pos
					item.Name = new(Name)
					item.Name.Value = strings.TrimSpace(value)
				}
			}
			rows.Close()
		}
	}

	tablesAndColumn := map[string]string{
		"eprint_corp_creators_id":  "corp_creators_id",
		"eprint_corp_creators_ror": "corp_creators_ror",
		"eprint_corp_creators_uri": "corp_creators_uri",
		"eprint_corp_creators":     "corp_creators",
	}
	for tableName, columnName := range tablesAndColumn {
		_, ok := tables[tableName]
		if ok {
			stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
			rows, err := db.Query(stmt, eprintID)
			if err != nil {
				log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &value); err != nil {
						log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
					} else {
						resizeItemList(pos, itemList)
						item := itemList.IndexOf(pos)
						item.Pos = pos
						switch columnName {
						case "corp_creators_id":
							item.ID = value
						case "corp_creators_ror":
							item.ROR = value
						case "corp_creators_uri":
							item.URI = value
						case "corp_creators":
							item.Name = new(Name)
							item.Name.Value = strings.TrimSpace(value)
						}
						break
					}
					i++
				}
				rows.Close()
			}
		}
	}
	if itemList.Length() > 0 {
		return itemList
	}
	return nil
}

func eprintIDToCorpContributors(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *CorpContributorItemList {
	var (
		pos   int
		value string
	)
	itemList := new(CorpContributorItemList)
	tableName := `eprint_corp_contributors_name`
	columnName := `corp_contributors_name`
	_, ok := tables[tableName]
	if ok {
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s
FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					resizeItemList(pos, itemList)
					item := itemList.IndexOf(pos)
					item.Pos = pos
					item.Name = new(Name)
					item.Name.Value = strings.TrimSpace(value)
				}
			}
			rows.Close()
		}
	}

	tablesAndColumn := map[string]string{
		"eprint_corp_contributors_id":  "corp_contributors_id",
		"eprint_corp_contributors_ror": "corp_contributors_ror",
		"eprint_corp_contributors_uri": "corp_contributors_uri",
		"eprint_corp_contributors":     "corp_contributors",
	}
	for tableName, columnName := range tablesAndColumn {
		_, ok := tables[tableName]
		if ok {
			stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
			rows, err := db.Query(stmt, eprintID)
			if err != nil {
				log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					if err := rows.Scan(&pos, &value); err != nil {
						log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
					} else {
						resizeItemList(pos, itemList)
						item := itemList.IndexOf(pos)
						item.Pos = pos
						switch columnName {
						case "corp_contributors_id":
							item.ID = value
						case "corp_contributors_ror":
							item.ROR = value
						case "corp_contributors_uri":
							item.URI = value
						case "corp_contributors":
							item.Name = new(Name)
							item.Name.Value = strings.TrimSpace(value)
						}
						break
					}
					i++
				}
				rows.Close()
			}
		}
	}
	if itemList.Length() > 0 {
		return itemList
	}
	return nil
}

func eprintIDToFunders(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *FunderItemList {
	var (
		pos   int
		value string
	)
	tableName := `eprint_funders_agency`
	columnName := `funders_agency`
	_, ok := tables[tableName]
	if ok {
		// eprint_%_id is a known structure. eprintid, pos, contributors_id
		itemList := new(FunderItemList)
		stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &value); err != nil {
					log.Printf("Could not scan %s for %d in %q, %s", tableName, eprintID, repoID, err)
				} else {
					resizeItemList(pos, itemList)
					item := itemList.IndexOf(pos)
					item.Pos = pos
					if value != "" {
						item.Agency = value
					}
				}
				i++
			}
			rows.Close()
		}
		tablesAndColumns := map[string]string{
			"eprint_funders_grant_number": "funders_grant_number",
			"eprint_funders_ror":          "funders_ror",
		}
		if itemList.Length() > 0 {
			for tableName, columnName := range tablesAndColumns {
				if _, ok := tables[tableName]; ok {
					stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
					rows, err = db.Query(stmt, eprintID)
					if err != nil {
						log.Printf("Could not query %s for %d in %q, %s", tableName, eprintID, repoID, err)
					} else {
						i := 0
						for rows.Next() {
							if err := rows.Scan(&pos, &value); err != nil {
								log.Printf("Could not scan (%d) %s for %d in %q, %s", i, tableName, eprintID, repoID, err)
							} else {
								resizeItemList(pos, itemList)
								if value != "" {
									for _, item := range itemList.Items {
										if item.Pos == pos {
											switch columnName {
											case "funders_grant_number":
												item.GrantNumber = value
											case "funders_ror":
												item.ROR = value
											}
											break
										}
									}
								}
							}
							i++
						}
						rows.Close()
					}
				}
			}
			return itemList
		}
	}
	return nil
}

func eprintIDToRelatedURL(repoID string, baseURL string, eprintID int, db *sql.DB, tables map[string][]string) *RelatedURLItemList {
	tablesAndColumns := map[string]string{
		"eprint_related_url_url":         "related_url_url",
		"eprint_related_url_type":        "related_url_type",
		"eprint_related_url_description": "related_url_description",
	}
	itemList := new(RelatedURLItemList)
	for tableName, columnName := range tablesAndColumns {
		if _, ok := tables[tableName]; ok {
			stmt := fmt.Sprintf(`SELECT pos, IFNULL(%s, "") AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
			rows, err := db.Query(stmt, eprintID)
			if err != nil {
				log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
			} else {
				i := 0
				for rows.Next() {
					var (
						pos   int
						value string
					)
					if err := rows.Scan(&pos, &value); err != nil {
						log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
					} else {
						if value != "" {
							resizeItemList(pos, itemList)
							item := itemList.IndexOf(pos)
							item.Pos = pos
							switch columnName {
							case `related_url_url`:
								item.URL = value
							case `related_url_type`:
								item.Type = value
							case `related_url_description`:
								item.Description = value
							}
						}
					}
					i++
				}
				rows.Close()
			}
		}
	}
	if itemList.Length() > 0 {
		return itemList
	}
	return nil
}

func eprintIDToOtherNumberingSystem(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *OtherNumberingSystemItemList {
	tableNames := []string{`eprint_other_numbering_system_name`, `eprint_other_numbering_system_id`}
	ok := true
	for _, tableName := range tableNames {
		if _, hasTable := tables[tableName]; !hasTable {
			ok = false
			break
		}
	}
	if ok {
		itemList := new(OtherNumberingSystemItemList)
		stmt := fmt.Sprintf(`
SELECT %s.pos AS pos, IFNULL(other_numbering_system_name, '') AS name, IFNULL(other_numbering_system_id, '') AS systemid FROM %s JOIN %s ON (%s.eprintid = %s.eprintid AND %s.pos = %s.pos) WHERE %s.eprintid = ?`, tableNames[0], tableNames[0], tableNames[1], tableNames[0], tableNames[1], tableNames[0], tableNames[1], tableNames[0])
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
		} else {
			var (
				pos      int
				name, id string
			)
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &name, &id); err != nil {
					log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
				} else {
					resizeItemList(pos, itemList)
					item := itemList.IndexOf(pos)
					item.Pos = pos
					item.ID = id
					if name != "" {
						item.Name = new(Name)
						item.Name.Value = strings.TrimSpace(name)
					}
				}
				i++
			}
			rows.Close()
			if itemList.Length() > 0 {
				return itemList
			}
		}
	}
	return nil
}

func eprintIDToItemIssues(repoID string, eprintID int, db *sql.DB, tables map[string][]string) *ItemIssueItemList {
	tableName := `eprint_item_issues_timestamp`
	_, ok := tables[tableName]
	if ok {
		var (
			year, month, day, hour, minute, second, pos int
			value                                       string
		)
		itemList := new(ItemIssueItemList)
		stmt := fmt.Sprintf(`SELECT pos,
IFNULL(item_issues_timestamp_year, 0) AS year,
IFNULL(item_issues_timestamp_month, 0) AS month,
IFNULL(item_issues_timestamp_day, 0) AS day,
IFNULL(item_issues_timestamp_hour, 0) AS hour,
IFNULL(item_issues_timestamp_minute, 0) AS minute,
IFNULL(item_issues_timestamp_second, 0) AS second
FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, tableName)
		rows, err := db.Query(stmt, eprintID)
		if err != nil {
			log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
		} else {
			i := 0
			for rows.Next() {
				if err := rows.Scan(&pos, &year, &month, &day, &hour, &minute, &second); err != nil {
					log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
				} else {
					resizeItemList(pos, itemList)
					item := itemList.IndexOf(pos)
					item.Pos = pos
					item.Timestamp = makeTimestamp(year, month, day, hour, minute, second)
				}
				i++
			}
			rows.Close()
		}
		if itemList.Length() > 0 {
			tablesAndColumn := map[string]string{
				"eprint_item_issues_type":        "item_issues_type",
				"eprint_item_issues_status":      "item_issues_status",
				"eprint_item_issues_description": "item_issues_description",
				"eprint_item_issues_id":          "item_issues_id",
				"eprint_item_issues_resolved_by": "item_issues_resolved_by",
				"eprint_item_issues_reported_by": "item_issues_reported_by",

				"eprint_item_issues_comment": "item_issues_comment",
			}
			for tableName, columnName := range tablesAndColumn {
				if _, ok := tables[tableName]; ok {
					stmt = fmt.Sprintf(`SELECT pos, IFNULL(%s, '') AS %s FROM %s WHERE eprintid = ? ORDER BY eprintid, pos`, columnName, columnName, tableName)
					rows, err := db.Query(stmt, eprintID)
					if err != nil {
						log.Printf("Could not query %d in %q, %s", eprintID, repoID, err)
					} else {
						i := 0
						for rows.Next() {
							if err := rows.Scan(&pos, &value); err != nil {
								log.Printf("Could not scan (%d) %d for %s, %s", i, eprintID, repoID, err)
							} else {
								for _, item := range itemList.Items {
									if item.Pos == pos && strings.TrimSpace(value) != "" {
										switch strings.TrimPrefix(columnName, "item_issues_") {
										case "type":
											item.Type = value
										case "status":
											item.Status = value
										case "description":
											item.Description = value
										case "id":
											item.ID = value
										case "resolved_by":
											item.ResolvedBy = value
										case "reported_by":
											item.ReportedBy = value
										case "comment":
											item.Comment = value
										}
									}
								}
							}
							i++
						}
						rows.Close()
					}
				}
			}
			return itemList
		}
	}
	return nil
}

// SQLReadEPrint expects a repository map and EPrint ID
// and will generate a series of SELECT statements populating
// a new EPrint struct or return an error (e.g. "not found" if eprint id is not in repository)
func SQLReadEPrint(config *Config, repoID string, baseURL string, eprintID int) (*EPrint, error) {
	var (
		tables  map[string][]string
		columns []string
	)
	if eprintID == 0 {
		return nil, fmt.Errorf("not found, %d not in %q", eprintID, repoID)
	}
	_, ok := config.Repositories[repoID]
	if !ok {
		return nil, fmt.Errorf("not found, %q not defined", repoID)
	}
	tables = config.Repositories[repoID].TableMap
	columns, ok = tables["eprint"]
	if !ok {
		return nil, fmt.Errorf("not found, %q eprint table not defined", repoID)
	}
	db, ok := config.Connections[repoID]
	if !ok {
		return nil, fmt.Errorf("no database connection for %s", repoID)
	}

	// NOTE: since the specific subset of columns in a repository
	// are known only at run time we need to setup a generic pointer
	// array for the scan results based on our newly allocated
	// EPrint struct.

	eprint := new(EPrint) // Generate an empty EPrint struct
	eprint.EPrintID = eprintID

	// NOTE: The data is littered with NULLs in EPrints. We need to
	// generate both a map of values into the EPrint stucture and
	// aggregated the SQL Column definitions to deal with the NULL
	// values.
	columnSQL, values := eprintToColumnsAndValues(eprint, columns, true)

	// NOTE: With the "values" pointer array setup the query can be built
	// and executed in the usually SQL fashion.
	stmt := fmt.Sprintf(`SELECT %s FROM eprint WHERE eprintid = ? LIMIT 1`, strings.Join(columnSQL, `, `))
	rows, err := db.Query(stmt, eprintID)
	if err != nil {
		return nil, fmt.Errorf(`ERROR: query error (%q, %q), %s`, repoID, stmt, err)
	}
	cnt := 0
	for rows.Next() {
		// NOTE: Because values array holds the addresses into our
		// EPrint struct the "Scan" does the actual mapping.
		// This makes it sorta "auto-magical"
		if err := rows.Scan(values...); err != nil {
			log.Printf(`%s.eprint eprintid = %d, %s`, repoID, eprintID, err)
		}
		cnt++
	}
	rows.Close()
	// NOTE: need to handle zero rows returned!
	if cnt > 0 {
		// Normalize fields inferred from MySQL database tables.
		eprint.ID = fmt.Sprintf(`%s/id/eprint/%d`, baseURL, eprint.EPrintID)
		eprint.LastModified = makeTimestamp(eprint.LastModifiedYear, eprint.LastModifiedMonth, eprint.LastModifiedDay, eprint.LastModifiedHour, eprint.LastModifiedMinute, eprint.LastModifiedSecond)
		// NOTE: EPrint XML uses a datestamp for output but tracks a timestamp.
		eprint.Datestamp = makeTimestamp(eprint.DatestampYear, eprint.DatestampMonth, eprint.DatestampDay, eprint.DatestampHour, eprint.DatestampMinute, eprint.DatestampSecond)
		eprint.StatusChanged = makeTimestamp(eprint.StatusChangedYear, eprint.StatusChangedMonth, eprint.StatusChangedDay, eprint.StatusChangedHour, eprint.StatusChangedMinute, eprint.StatusChangedSecond)
		eprint.Date = makeApproxDate(eprint.DateYear, eprint.DateMonth, eprint.DateDay)

		// FIXME: Add Depository info (eprint.userid -> user* tables)
		//   deposited on, deposited by
		if eprint.UserID > 0 {
			eprint.DepositedBy = userIDToName(repoID, eprint.UserID, db)
			eprint.DepositedOn = makeTimestamp(eprint.DatestampYear, eprint.DatestampMonth, eprint.DatestampDay, eprint.DatestampHour, eprint.DatestampMinute, eprint.DatestampSecond)
		}

		// Used in CaltechTHESIS
		eprint.ThesisSubmittedDate = makeDatestamp(eprint.ThesisSubmittedDateYear, eprint.ThesisSubmittedDateMonth, eprint.ThesisSubmittedDateDay)
		eprint.ThesisDefenseDate = makeDatestamp(eprint.ThesisDefenseDateYear, eprint.ThesisDefenseDateMonth, eprint.ThesisDefenseDateDay)
		eprint.ThesisApprovedDate = makeDatestamp(eprint.ThesisApprovedDateYear, eprint.ThesisApprovedDateMonth, eprint.ThesisApprovedDateDay)
		eprint.ThesisPublicDate = makeDatestamp(eprint.ThesisPublicDateYear, eprint.ThesisPublicDateMonth, eprint.ThesisPublicDateDay)
		eprint.ThesisDegreeDate = makeDatestamp(eprint.ThesisDegreeDateYear, eprint.ThesisDegreeDateMonth, eprint.ThesisDegreeDateDay)
		eprint.GradOfficeApprovalDate = makeDatestamp(eprint.GradOfficeApprovalDateYear, eprint.GradOfficeApprovalDateMonth, eprint.GradOfficeApprovalDateDay)

		// CreatorsItemList
		eprint.Creators = eprintIDToCreators(repoID, eprintID, db, tables)
		// EditorsItemList
		eprint.Editors = eprintIDToEditors(repoID, eprintID, db, tables)
		// ContributorsItemList
		eprint.Contributors = eprintIDToContributors(repoID, eprintID, db, tables)

		// CorpCreators
		eprint.CorpCreators = eprintIDToCorpCreators(repoID, eprintID, db, tables)
		// CorpContributors
		eprint.CorpContributors = eprintIDToCorpContributors(repoID, eprintID, db, tables)

		// LocalGroupItemList (SimpleItemList)
		eprint.LocalGroup = eprintIDToLocalGroup(repoID, eprintID, db, tables)
		// FundersItemList (custom)
		eprint.Funders = eprintIDToFunders(repoID, eprintID, db, tables)
		// Documents (*DocumentList)
		eprint.Documents = eprintIDToDocumentList(repoID, baseURL, eprintID, db, tables)
		// RelatedURLs List
		eprint.RelatedURL = eprintIDToRelatedURL(repoID, baseURL, eprintID, db, tables)
		// ReferenceText (item list)
		eprint.ReferenceText = eprintIDToReferenceText(repoID, eprintID, db, tables)
		// Projects
		eprint.Projects = eprintIDToProjects(repoID, eprintID, db, tables)
		// OtherNumberingSystem (item list)
		eprint.OtherNumberingSystem = eprintIDToOtherNumberingSystem(repoID, eprintID, db, tables)
		// Subjects List
		eprint.Subjects = eprintIDToSubjects(repoID, eprintID, db, tables)
		// ItemIssues
		eprint.ItemIssues = eprintIDToItemIssues(repoID, eprintID, db, tables)

		// Exhibitors
		eprint.Exhibitors = eprintIDToExhibitors(repoID, eprintID, db, tables)
		// Producers
		eprint.Producers = eprintIDToProducers(repoID, eprintID, db, tables)
		// Conductors
		eprint.Conductors = eprintIDToConductors(repoID, eprintID, db, tables)

		// Lyricists
		eprint.Lyricists = eprintIDToLyricists(repoID, eprintID, db, tables)

		// Accompaniment
		eprint.Accompaniment = eprintIDToAccompaniment(repoID, eprintID, db, tables)
		// SkillAreas
		eprint.SkillAreas = eprintIDToSkillAreas(repoID, eprintID, db, tables)
		// CopyrightHolders
		eprint.CopyrightHolders = eprintIDToCopyrightHolders(repoID, eprintID, db, tables)
		// Reference
		eprint.Reference = eprintIDToReference(repoID, eprintID, db, tables)

		// ConfCreators
		eprint.ConfCreators = eprintIDToConfCreators(repoID, eprintID, db, tables)
		// AltTitle
		eprint.AltTitle = eprintIDToAltTitle(repoID, eprintID, db, tables)
		// PatentAssignee
		eprint.PatentAssignee = eprintIDToPatentAssignee(repoID, eprintID, db, tables)
		// RelatedPatents
		eprint.RelatedPatents = eprintIDToRelatedPatents(repoID, eprintID, db, tables)
		// Divisions
		eprint.Divisions = eprintIDToDivisions(repoID, eprintID, db, tables)
		// ThesisAdvisor
		eprint.ThesisAdvisor = eprintIDToThesisAdvisors(repoID, eprintID, db, tables)
		// ThesisCommittee
		eprint.ThesisCommittee = eprintIDToThesisCommittee(repoID, eprintID, db, tables)

		// OptionMajor
		eprint.OptionMajor = eprintIDToOptionMajor(repoID, eprintID, db, tables)
		// OptionMinor
		eprint.OptionMinor = eprintIDToOptionMinor(repoID, eprintID, db, tables)

		/*************************************************************
		    NOTE: These are notes about possible original implementation
		    errors or elements that did not survive the upgrade to
		    EPrints 3.3.16

		    eprint.LearningLevels (not an item list in EPrints) using LearningLevelText
		    GScholar, skipping not an item list, a 2010 plugin for EPRints 3.2.
		    eprint.GScholar = eprintIDToGScholar(repoID, eprintID, db, tables)
		    Shelves, a plugin, not replicating, not an item list
		    eprint.Shelves = eprintIDToSchelves(repoID, eprintID, db, tables)
		    eprint.PatentClassification is not not an item list, using eprint.PatentClassificationText
		    eprint.OtherURL appears to be an extraneous
		    eprint.CorpContributors apears to be an extraneous
		*************************************************************/
	} else {
		return nil, fmt.Errorf("not found")
	}

	return eprint, nil
}

// qmList generates an array of string where each element holds "?".
func qmList(length int) []string {
	list := []string{}
	for i := 0; i < length; i++ {
		list = append(list, `?`)
	}
	return list
}

// insertItemList takes an repoID, table name, list of columns and
// an EPrint datastructure then generates and executes a series of
// INSERT statement to create an Item List for the given table.
func insertItemList(db *sql.DB, repoID string, tableName string, columns []string, eprint *EPrint) error {
	var (
		itemList ItemsInterface
	)
	eprintid := eprint.EPrintID
	switch {
	case strings.HasPrefix(tableName, `eprint_creators_`):
		itemList = eprint.Creators
	case strings.HasPrefix(tableName, `eprint_editors_`):
		itemList = eprint.Editors
	case strings.HasPrefix(tableName, `eprint_contributors_`):
		itemList = eprint.Contributors
	case strings.HasPrefix(tableName, `eprint_corp_creators`):
		itemList = eprint.CorpCreators
	case strings.HasPrefix(tableName, `eprint_corp_contributors_`):
		itemList = eprint.CorpContributors
	case strings.HasPrefix(tableName, `eprint_thesis_advisor_`):
		itemList = eprint.ThesisAdvisor
	case strings.HasPrefix(tableName, `eprint_thesis_committee_`):
		itemList = eprint.ThesisCommittee
	case strings.HasPrefix(tableName, `eprint_item_issues_`):
		itemList = eprint.ItemIssues
	case strings.HasPrefix(tableName, `eprint_alt_title`):
		itemList = eprint.AltTitle
	case strings.HasPrefix(tableName, `eprint_conductors`):
		itemList = eprint.Conductors
	case strings.HasPrefix(tableName, `eprint_conf_creators_`):
		itemList = eprint.ConfCreators
	case strings.HasPrefix(tableName, `eprint_exhibitors_`):
		itemList = eprint.Exhibitors
	case strings.HasPrefix(tableName, `eprint_producers_`):
		itemList = eprint.Producers
	case strings.HasPrefix(tableName, `eprint_lyricists_`):
		itemList = eprint.Lyricists
	case strings.HasPrefix(tableName, `eprint_accompaniment`):
		itemList = eprint.Accompaniment
	case strings.HasPrefix(tableName, `eprint_subjec`):
		itemList = eprint.Subjects
	case strings.HasPrefix(tableName, `eprint_local_`):
		itemList = eprint.LocalGroup
	case strings.HasPrefix(tableName, `eprint_div`):
		itemList = eprint.Divisions
	case strings.HasPrefix(tableName, `eprint_option_maj`):
		itemList = eprint.OptionMajor
	case strings.HasPrefix(tableName, `eprint_option_min`):
		itemList = eprint.OptionMinor
	case strings.HasPrefix(tableName, `eprint_funders_`):
		itemList = eprint.Funders
	case strings.HasPrefix(tableName, `eprint_funders`):
		// Ignore, eprint_funders is empty in CaltechAUTHORS ...
		itemList = new(FunderItemList)
	case strings.HasPrefix(tableName, `eprint_other_numbering_system`):
		itemList = eprint.OtherNumberingSystem
	case strings.HasPrefix(tableName, `eprint_projects`):
		itemList = eprint.Projects
	case strings.HasPrefix(tableName, `eprint_referencetext`):
		itemList = eprint.ReferenceText
	case strings.HasPrefix(tableName, `eprint_related_url`):
		itemList = eprint.RelatedURL
	case strings.HasPrefix(tableName, `eprint_skill_areas`):
		itemList = eprint.SkillAreas
	case strings.HasPrefix(tableName, `eprint_patent_assignee`):
		itemList = eprint.PatentAssignee
	case strings.HasPrefix(tableName, `eprint_related_patents`):
		itemList = eprint.RelatedPatents
	case strings.HasPrefix(tableName, `eprint_referencetext`):
		itemList = eprint.ReferenceText
	case strings.HasPrefix(tableName, `eprint_accompaniment`):
		itemList = eprint.Accompaniment
	case strings.HasPrefix(tableName, `eprint_reference`):
		itemList = eprint.Reference
	case strings.HasPrefix(tableName, `eprint_copyright_holders`):
		itemList = eprint.CopyrightHolders
	case strings.HasPrefix(tableName, `eprint_related_patent`):
		itemList = eprint.RelatedPatents
	case strings.HasPrefix(tableName, `eprint_parent_assign`):
		itemList = eprint.PatentAssignee
	case strings.HasPrefix(tableName, `eprint_skill`):
		itemList = eprint.SkillAreas
	case strings.HasPrefix(tableName, `eprint_relation`):
		// NOTE: This is not the same as document_relation_*, it is a separate item list item list
		// it has the same structure with a uri and type. Our eprint implementations use a Relation
		itemList = eprint.Relation
	case strings.HasPrefix(tableName, `eprint_keyword`):
		// NOTE: this we appear to use the longtext of key in our eprint table. Not sure if this
		// is new or old structure. It is posssible that our longtext for keywords is a legacy structure.
		// itemList = eprint.Keyword
	default:
		return fmt.Errorf(`do not understand table %q, columns %s`, tableName, strings.Join(columns, `, `))
	}
	// Clear the list, then insert
	stmt := fmt.Sprintf(`DELETE FROM %s WHERE eprintid = ?`, tableName)
	_, err := db.Exec(stmt, eprint.EPrintID)
	if err != nil {
		return fmt.Errorf(`SQL error, %q, %s`, stmt, err)
	}
	for pos := 0; pos < itemList.Length(); pos++ {
		item := itemList.IndexOf(pos)
		item.Pos = pos
		values := []interface{}{}
		columnsSQL := []string{}
		for _, col := range columns {
			switch {
			case col == `eprintid`:
				values = append(values, eprintid)
				columnsSQL = append(columnsSQL, col)
			case col == `pos`:
				values = append(values, pos)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_id`):
				values = append(values, item.ID)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_type`):
				values = append(values, item.Type)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_family`):
				values = append(values, item.Name.Family)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_given`):
				values = append(values, item.Name.Given)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_honourific`):
				values = append(values, item.Name.Honourific)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_lineage`):
				values = append(values, item.Name.Lineage)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_name`):
				values = append(values, item.Name.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_show_email`):
				// NOTE: _show_email needs to be tested before _email
				values = append(values, item.ShowEMail)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_email`):
				// NOTE: _show_email needs to be tested before _email
				values = append(values, item.EMail)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_role`):
				values = append(values, item.Role)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_url`):
				values = append(values, item.URL)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `description`):
				values = append(values, item.Description)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_agency`):
				values = append(values, item.Agency)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_grant_number`):
				values = append(values, item.GrantNumber)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_uri`):
				values = append(values, item.URI)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_orcid`):
				values = append(values, item.ORCID)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_ror`):
				values = append(values, item.ROR)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_timestamp`):
				values = append(values, item.Timestamp)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_status`):
				values = append(values, item.Status)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_reported_by`):
				values = append(values, item.ReportedBy)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_resolved_by`):
				values = append(values, item.ResolvedBy)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_comment`):
				values = append(values, item.ResolvedBy)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_group`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_subjects`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_major`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_minor`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `_holders`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `divisions`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `subjects`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `referencetext`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `accompaniment`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `related_patents`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `patent_assignee`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `skill_areas`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			case strings.HasSuffix(col, `alt_title`):
				values = append(values, item.Value)
				columnsSQL = append(columnsSQL, col)
			default:
				return fmt.Errorf("do not understand column %s.%s\n", tableName, col)
			}
		}
		stmt := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tableName, strings.Join(columnsSQL, `, `), strings.Join(qmList(len(columnsSQL)), `, `))
		_, err := db.Exec(stmt, values...)
		if err != nil {
			return fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
	}
	return nil
}

// SQLCreateEPrint will read a EPrint structure and
// generate SQL INSERT, REPLACE and DELETE statements
// suitable for creating a new EPrint record in the repository.
func SQLCreateEPrint(config *Config, repoID string, ds *DataSource, eprint *EPrint) (int, error) {
	var (
		err error
	)
	db, ok := config.Connections[repoID]
	if !ok {
		return 0, fmt.Errorf(`no database connection for %s`, repoID)
	}
	// If eprint id is zero generate a sequence of INSERT statements
	// for the record. Others use generate the appropriate History
	// records and then delete insert the new record.
	tableName := `eprint`

	if columns, ok := ds.TableMap[tableName]; ok {
		// Generate an empty row and capture the id created.
		stmt := `INSERT INTO eprint (eprintid) (SELECT (IFNULL((SELECT eprintid FROM eprint ORDER BY eprintid DESC LIMIT 1), 0) + 1) AS eprintid)`
		_, err := db.Exec(stmt)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		stmt = `SELECT eprintid FROM eprint ORDER BY eprintid DESC LIMIT 1`
		rows, err := db.Query(stmt)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
		id := 0
		for rows.Next() {
			if err := rows.Scan(&id); err != nil {
				return 0, fmt.Errorf(`could not calculate the new eprintid value, %s`, err)
			}
		}
		rows.Close()
		if err != nil {
			return 0, fmt.Errorf(`SQL failed to get insert id, %s`, err)
		}
		eprint.EPrintID = int(id)
		eprint.Dir = makeDirValue(eprint.EPrintID)
		// FIXME: decide if the is automatic or if this should be
		// passed in with the data structure.
		// Generate minimal date and time stamps
		now := time.Now()
		if eprint.Datestamp == "" {
			eprint.Datestamp = now.Format(timestamp)
			eprint.DatestampYear = now.Year()
			eprint.DatestampMonth = int(now.Month())
			eprint.DatestampDay = now.Day()
			eprint.DatestampHour = now.Hour()
			eprint.DatestampMinute = now.Minute()
			eprint.DatestampSecond = now.Second()
		} else if dt, err := time.Parse(datestamp, eprint.Datestamp); err == nil {
			eprint.DatestampYear = dt.Year()
			eprint.DatestampMonth = int(dt.Month())
			eprint.DatestampDay = dt.Day()
		} else if dt, err := time.Parse(timestamp, eprint.Datestamp); err == nil {
			eprint.DatestampYear = dt.Year()
			eprint.DatestampMonth = int(dt.Month())
			eprint.DatestampDay = dt.Day()
			eprint.DatestampHour = dt.Hour()
			eprint.DatestampMinute = dt.Minute()
			eprint.DatestampSecond = dt.Second()
		}

		eprint.LastModified = now.Format(timestamp)
		eprint.LastModifiedYear = now.Year()
		eprint.LastModifiedMonth = int(now.Month())
		eprint.LastModifiedDay = now.Day()
		eprint.LastModifiedHour = now.Hour()
		eprint.LastModifiedMinute = now.Minute()
		eprint.LastModifiedSecond = now.Second()

		eprint.StatusChanged = now.Format(timestamp)
		eprint.StatusChangedYear = now.Year()
		eprint.StatusChangedMonth = int(now.Month())
		eprint.StatusChangedDay = now.Day()
		eprint.StatusChangedHour = now.Hour()
		eprint.StatusChangedMinute = now.Minute()
		eprint.StatusChangedSecond = now.Second()

		if eprint.Date != "" {
			eprint.DateYear, eprint.DateMonth, eprint.DateDay = approxYMD(eprint.Date)
		}
		if eprint.ThesisSubmittedDate != "" {
			eprint.ThesisSubmittedDateYear, eprint.ThesisSubmittedDateMonth, eprint.ThesisSubmittedDateDay = approxYMD(eprint.ThesisSubmittedDate)
		}
		if eprint.ThesisDefenseDate != "" {
			eprint.ThesisDefenseDateYear, eprint.ThesisDefenseDateMonth, eprint.ThesisDefenseDateDay = approxYMD(eprint.ThesisDefenseDate)
		}
		if eprint.ThesisApprovedDate != "" {
			eprint.ThesisApprovedDateYear, eprint.ThesisApprovedDateMonth, eprint.ThesisApprovedDateDay = approxYMD(eprint.ThesisApprovedDate)
		}
		if eprint.ThesisPublicDate != "" {
			eprint.ThesisPublicDateYear, eprint.ThesisPublicDateMonth, eprint.ThesisPublicDateDay = approxYMD(eprint.ThesisPublicDate)
		}
		if eprint.GradOfficeApprovalDate != "" {
			eprint.GradOfficeApprovalDateYear, eprint.GradOfficeApprovalDateMonth, eprint.GradOfficeApprovalDateDay = approxYMD(eprint.GradOfficeApprovalDate)
		}

		// Step two, write the rest of the date into the main table.
		columnsSQL, values := eprintToColumnsAndValues(eprint, columns, false)
		stmt = fmt.Sprintf(`REPLACE INTO %s (%s) VALUES (%s)`,
			tableName,
			strings.Join(columnsSQL, `, `),
			strings.Join(qmList(len(columnsSQL)), `, `))
		_, err = db.Exec(stmt, values...)
		if err != nil {
			return 0, fmt.Errorf(`SQL error, %q, %s`, stmt, err)
		}
	}
	if eprint.EPrintID != 0 {
		for tableName, columns := range ds.TableMap {
			// Handle the remaining tables, i.e. skip eprint table.
			switch {
			case tableName == `eprint`:
				// Skip eprint table, we've already processed it
			case tableName == `eprint_keyword`:
				// Skip eprint_keyword, our EPrints use keywords (longtext) in eprint table.
			case strings.HasPrefix(tableName, `document`):
				//log.Printf(`FIXME %s columns: %s`, tableName, strings.Join(columns, `, `))
			case strings.HasPrefix(tableName, `file`):
				//log.Printf(`FIXME %s columns: %s`, tableName, strings.Join(columns, `, `))
			default:
				// Insert new rows in associated table
				if err := insertItemList(db, repoID, tableName, columns, eprint); err != nil {
					return eprint.EPrintID, fmt.Errorf(`failed to insert eprintid %d in table %s for %s, %s`, eprint.EPrintID, tableName, repoID, err)
				}
			}
		}
		return eprint.EPrintID, nil
	}
	return 0, err
}

// ImportEPrints take an repository id, eprints structure.
// It is a re-implementation of EPrints Perl Import EPrint XML.
//
// ImportEPrints will create EPrint records using
// SQLCreateEPrint function. ImportEPrint is responsible for
// applly any rule sets and eprint status, SQLCreateEPRint is
// responsible for datestamps and timestamps.
//
// ImportEPrints returns a list of EPrint IDs created
// if successful and an error if something goes wrong.
func ImportEPrints(config *Config, repoID string, ds *DataSource, eprints *EPrints) ([]int, error) {
	var importErrors error
	ids := []int{}

	if config.Connections == nil {
		return nil, fmt.Errorf(`no databases are not configured`)
	}
	_, ok := config.Connections[repoID]
	if !ok {
		return nil, fmt.Errorf(`%s database connection not configured`, repoID)
	}

	// Check to make sure updates are allowed if non-Zero
	// eprint ids present.
	for _, eprint := range eprints.EPrint {
		if eprint.EPrintID != 0 {
			return nil, fmt.Errorf("create failed eprint id %d in %s", eprint.EPrintID, repoID)
		}
		if eprint.Collection == "" && ds.DefaultCollection != "" {
			eprint.Collection = DefaultCollection
		}
		if eprint.IDNumber == "" && ds.DefaultOfficialURL != "" {
			eprint.IDNumber = GenerateIDNumber(eprint)
		}
		if eprint.OfficialURL == "" && ds.DefaultOfficialURL != "" {
			eprint.OfficialURL = GenerateOfficialURL(eprint)
		}
		if eprint.Rights == "" && ds.DefaultRights != "" {
			eprint.Rights = ds.DefaultRights
		}
		if eprint.Refereed == "" && eprint.Type == "article" &&
			ds.DefaultRefereed != "" {
			eprint.Refereed = ds.DefaultRefereed
		}
		if eprint.EPrintStatus == "" && ds.DefaultStatus != "" {
			eprint.EPrintStatus = ds.DefaultStatus
		}
		if eprint.Abstract != "" && ds.StripTags {
			if cleaner.HasEncodedElements([]byte(eprint.Abstract)) {
				eprint.Abstract = string(cleaner.StripTags([]byte(eprint.Abstract)))
			}
		}
	}
	for _, eprint := range eprints.EPrint {
		id, err := SQLCreateEPrint(config, repoID, ds, eprint)
		if err != nil {
			if importErrors == nil {
				importErrors = err
			} else {
				importErrors = fmt.Errorf("%s; %s", importErrors, err)
			}
		}
		ids = append(ids, id)
	}
	return ids, importErrors
}

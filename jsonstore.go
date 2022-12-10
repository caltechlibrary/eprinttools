// jsonstore.go holds the operations for openning and close a JSON
// Document Store currently implemented in MySQL 8 using JSON columns.
package eprinttools

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// MySQL database support
	_ "github.com/go-sql-driver/mysql"
)

// OpenJSONStore
func OpenJSONStore(config *Config) error {
	if config.JSONStore == "" {
		return fmt.Errorf("JSONStore is not set")
	} else {
		// Setup DB connection for target repository
		db, err := sql.Open("mysql", config.JSONStore)
		if err != nil {
			return fmt.Errorf("Could not open MySQL connection for %s, %s", config.JSONStore, err)
		}
		config.Jdb = db
	}
	return nil
}

// CloseJSONStore
func CloseJSONStore(config *Config) error {
	if config.Jdb != nil {
		if err := config.Jdb.Close(); err != nil {
			return fmt.Errorf("Failed to close %s, %s", config.JSONStore, err)
		}
	}
	return nil
}

// Doc holds the data structure of our jsonstore row.
type Doc struct {
	ID           int    `json:"id"`
	Src          []byte `json:"src"`
	Action       string `json:"action,omitempty"`
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastmod,omitempty"`
	PubDate      string `json:"pubDate,omitempty"`
	Status       string `json:"status,omitempty"`
	IsPublic     bool   `json:"is_public,omitempty"`
	RecordType   string `json:"record_type,omitempty"`
}

// SaveJSONDocument takes a configuration, repoName, eprint id as integer and
// JSON source saving it to the appropriate JSON table.
func SaveJSONDocument(cfg *Config, repoName string, id int, src []byte, action string, created string, lastmod string, pubdate string, status string, isPublic bool, recordType string) error {

	stmt := fmt.Sprintf(`REPLACE INTO %s (id, src, action, created, lastmod, pubdate, status, is_public, record_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, repoName)
	doc := new(Doc)
	doc.ID = id
	doc.Src = src
	doc.Action = action
	doc.Created = created
	doc.LastModified = lastmod
	doc.PubDate = pubdate
	doc.Status = status
	doc.IsPublic = isPublic
	doc.RecordType = recordType
	_, err := cfg.Jdb.Exec(stmt, doc.ID, doc.Src, doc.Action, doc.Created, doc.LastModified, doc.PubDate, doc.Status, doc.IsPublic, doc.RecordType)
	if err != nil {
		return fmt.Errorf("sql failed for %d in %s, %s", id, repoName, err)
	}
	return nil
}

// GetJSONDocument takes a configuration, repoName, eprint id and returns
// the JSON source document.
func GetJSONDocument(cfg *Config, repoName string, id int) ([]byte, error) {
	stmt := fmt.Sprintf("SELECT src FROM %s WHERE id = ? LIMIT 1", repoName)
	rows, err := cfg.Jdb.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get %s id %d, %s", repoName, id, err)
	}
	defer rows.Close()
	var src []byte
	for rows.Next() {
		if err := rows.Scan(src); err != nil {
			return nil, fmt.Errorf("Failed to get row in %s for id %d, %s", repoName, id, err)
		}
	}
	err = rows.Err()
	return src, err
}

// GetJSONRow takes a configuration, repoName, eprint id and returns
// the table row as JSON source.
func GetJSONRow(cfg *Config, repoName string, id int) ([]byte, error) {
	stmt := fmt.Sprintf("SELECT id, src, action, created, lastmod, pubdate, status, is_public, record_type FROM %s WHERE id = ?", repoName)
	rows, err := cfg.Jdb.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get %s id %d, %s", repoName, id, err)
	}
	defer rows.Close()
	doc := new(Doc)
	for rows.Next() {
		if err := rows.Scan(doc.ID, doc.Src, doc.Action, doc.Created, doc.LastModified, doc.PubDate, doc.Status, doc.IsPublic, doc.RecordType); err != nil {
			return nil, fmt.Errorf("Failed to get row in %s for id %d, %s", repoName, id, err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return jsonEncode(doc)
}

func SavePersonJSON(cfg *Config, person *Person) error {
	stmt := `REPLACE INTO _people (person_id, cl_people_id,family_name,given_name, sort_name, thesis_id,advisor_id,authors_id,
        archivesspace_id,directory_id,viaf_id,lcnaf,
        isni,wikidata,snac,orcid,image,educated_at,caltech,jpl,faculty,alumn,
        status,directory_person_type,title,bio,division) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := cfg.Jdb.Exec(stmt, person.PersonID, person.CLPeopleID, person.FamilyName, person.GivenName, person.SortName, person.ThesisID, person.AdvisorID, person.AuthorsID,
		person.ArchivesSpaceID, person.DirectoryID, person.VIAF, person.LCNAF,
		person.ISNI, person.Wikidata, person.SNAC, person.ORCID, person.Image, person.EducatedAt, person.Caltech, person.JPL, person.Faculty, person.Alumn,
		person.Status, person.DirectoryPersonType, person.Title, person.Bio, person.Division)
	if err != nil {
		return err
	}
	return nil
}

func GetPerson(cfg *Config, personID string) (*Person, error) {
	stmt := `SELECT person_id, cl_people_id, family_name, given_name, sort_name, thesis_id, advisor_id, authors_id,
    archivesspace_id, directory_id, viaf_id, lcnaf, isni, wikidata, snac, orcid,
    image, educated_at, caltech, jpl, faculty, alumn, status, directory_person_type,
    title, bio, division, UNIX_TIMESTAMP(updated) FROM _people WHERE cl_people_id = ?`
	row, err := cfg.Jdb.Query(stmt, personID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var updated int64
	person := new(Person)
	if row.Next() {
		if err := row.Scan(&person.PersonID, &person.CLPeopleID, &person.FamilyName, &person.GivenName, &person.SortName, &person.ThesisID, &person.AdvisorID, &person.AuthorsID,
			&person.ArchivesSpaceID, &person.DirectoryID, &person.VIAF, &person.LCNAF, &person.ISNI, &person.Wikidata, &person.SNAC, &person.ORCID,
			&person.Image, &person.EducatedAt, &person.Caltech, &person.JPL, &person.Faculty, &person.Alumn, &person.Status, &person.DirectoryPersonType,
			&person.Title, &person.Bio, &person.Division, &updated); err != nil {
			return nil, err
		}
		person.Updated = time.Unix(updated, 0)
	}
	err = row.Err()
	return person, err
}

func GetPersonIDs(cfg *Config) ([]string, error) {
	stmt := `SELECT person_id FROM _people ORDER BY sort_name`
	rows, err := cfg.Jdb.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		if strings.TrimSpace(id) != "" {
			ids = append(ids, id)
		}
	}
	err = rows.Err()
	return ids, err
}

func SaveGroupJSON(cfg *Config, group *Group) error {
	stmt := `REPLACE INTO _groups (group_id,name,alternative,email,date,description,start,approx_start,activity,end,
    approx_end,website,pi,parent,prefix,grid,isni,ringold,viaf,ror) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := cfg.Jdb.Exec(stmt,
		group.GroupID, group.Name, group.Alternative, group.EMail, group.Date, group.Description, group.Start, group.ApproxStart, group.Activity, group.End,
		group.ApproxEnd, group.Website, group.PI, group.Parent, group.Prefix, group.GRID, group.ISNI, group.RinGold, group.VIAF, group.ROR)
	if err != nil {
		return err
	}
	return nil
}

func GetGroup(cfg *Config, groupID string) (*Group, error) {
	var updated int64
	stmt := `SELECT group_id, name, alternative, email, date, description, start,
    approx_start, activity, end, approx_end, website, pi,
    parent, prefix, grid, isni, ringold, viaf,
    ror, UNIX_TIMESTAMP(updated) FROM _groups WHERE group_id = ?`
	row, err := cfg.Jdb.Query(stmt, groupID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	group := new(Group)
	if row.Next() {
		if err := row.Scan(&group.GroupID, &group.Name, &group.Alternative, &group.EMail, &group.Date, &group.Description, &group.Start,
			&group.ApproxStart, &group.Activity, &group.End, &group.ApproxEnd, &group.Website, &group.PI,
			&group.Parent, &group.Prefix, &group.GRID, &group.ISNI, &group.RinGold, &group.VIAF,
			&group.ROR, &updated); err != nil {
			return nil, err
		}
	}
	err = row.Err()
	group.Updated = time.UnixMicro(updated)
	return group, err
}

func GetGroupIDByName(cfg *Config, groupName string) (string, error) {
	var groupID string
	stmt := `SELECT group_id FROM _groups WHERE name = ? OR (LOCATE(?, alternative) > 0) LIMIT 1`
	row, err := cfg.Jdb.Query(stmt, groupName, groupName)
	if err != nil {
		return "", err
	}
	defer row.Close()
	if row.Next() {
		if err := row.Scan(&groupID); err != nil {
			return "", err
		}
	}
	return groupID, nil
}

func GetGroupIDs(cfg *Config) ([]string, error) {
	stmt := `SELECT group_id FROM _groups ORDER BY name`
	rows, err := cfg.Jdb.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		if strings.TrimSpace(id) != "" {
			ids = append(ids, id)
		}
	}
	err = rows.Err()
	return ids, err
}

func GetGroupAggregations(cfg *Config, groupID string) (map[string]map[string][]int, error) {
	// Read the _aggregate_group to get the eprintid for group by decending publation date
	stmt := `SELECT repository, eprintid, record_type FROM _aggregate_groups WHERE group_id = ? ORDER BY repository, pubDate DESC`
	rows, err := cfg.Jdb.Query(stmt, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		id         int
		recordType string
		repoName   string
		repository string
	)
	m := map[string]map[string][]int{}
	for rows.Next() {
		if err := rows.Scan(&repository, &id, &recordType); err != nil {
			return nil, err
		}
		if repoName != repository {
			repoName = repository
			if _, ok := m[repoName]; !ok {
				m[repoName] = map[string][]int{}
			}
			m[repoName]["combined"] = []int{}
		}
		if _, ok := m[repoName][recordType]; !ok {
			m[repoName][recordType] = []int{}
		}
		m[repoName][recordType] = append(m[repoName][recordType], id)
		m[repoName]["combined"] = append(m[repoName]["combined"], id)
	}
	err = rows.Err()
	return m, err
}

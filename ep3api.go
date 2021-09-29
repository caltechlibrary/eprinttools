/**
 * ep3api.go defines an extended EPrints web server.
 */
package eprinttools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//
// Service configuration management
//

// Config holds a configuration file structure used by EPrints Extended API
// Configuration file is expected to be in JSON format.
type Config struct {
	// Hostname for running service
	Hostname string `json:"hostname"`

	// Repositories are defined by a REPO_ID (string)
	// that points at a MySQL Db connection string
	Repositories map[string]string `json:"repositories"`

	// Connections is a map to database connections
	Connections map[string]*sql.DB `json:"-"`

	// Routes holds the mapping of end points to repository id
	// instances.
	Routes map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error) `json:"-"`
}

const (
	// timestamp holds the Format of a MySQL time field
	timestamp = "2006-01-02 15:04:05"
)

var (
	config *Config
)

//
// End Point handles (route as defined `/<REPO_ID>/<END-POINT>/<ARGS>`)
//

func updatedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if db, ok := config.Connections[repoID]; ok == true {
		if (len(args) < 1) || (len(args) > 2) {
			return 400, fmt.Errorf("Bad Request")
		}
		var err error
		end := time.Now()
		start := time.Now()
		if (len(args) > 0) && (args[0] != "now") {
			start, err = time.Parse(timestamp, args[0])
			if err != nil {
				return 400, fmt.Errorf("Bad Request, %s", err)
			}
		}
		if (len(args) > 1) && (args[1] != "now") {
			end, err = time.Parse(timestamp, args[1])
			if err != nil {
				return 400, fmt.Errorf("Bad Request, %s", err)
			}
		}
		stmt := `SELECT eprintid FROM eprint WHERE (lastmod_year >= ? AND lastmod_year <= ?) AND (lastmod_month >= ? AND lastmod_month <= ?) AND (lastmod_day >= ? AND lastmod_day <= ?) AND (lastmod_hour >= ? AND lastmod_hour <= ?) AND (lastmod_minute >= ? AND lastmod_minute <=?)`
		rows, err := db.Query(stmt, start.Year(), end.Year(), start.Month(), end.Month(), start.Day(), end.Day(), start.Hour(), end.Hour(), start.Minute(), end.Minute())
		if err != nil {
			log.Printf("ERROR: query error (%q), %s", repoID, err)
			return 500, fmt.Errorf("Internal Server Error")
		}
		defer rows.Close()
		eprintid := 0
		eprintIDs := []int{}
		for rows.Next() {
			err := rows.Scan(&eprintid)
			if (err == nil) && (eprintid > 0) {
				eprintIDs = append(eprintIDs, eprintid)
			} else {
				log.Printf("ERROR: scan error (%q), %s", repoID, err)
				return 500, fmt.Errorf("Internal Server Error")
			}
		}
		if err := rows.Err(); err != nil {
			log.Printf("ERROR: rows error (%q), %s", repoID, err)
			return 500, fmt.Errorf("Internal Server Error")
		}
		src, err := json.MarshalIndent(eprintIDs, "", "  ")
		if err != nil {
			log.Printf("ERROR: marshal error (%q), %s", repoID, err)
			return 500, fmt.Errorf("Internal Server Error")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", src)
		return 200, nil
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func deletedEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if db, ok := config.Connections[repoID]; ok == true {
		if (len(args) < 1) || (len(args) > 2) {
			return 400, fmt.Errorf("Bad Request")
		}
		var err error
		end := time.Now()
		start := time.Now()
		if (len(args) > 0) && (args[0] != "now") {
			start, err = time.Parse(timestamp, args[0])
			if err != nil {
				return 400, fmt.Errorf("Bad Request, %s", err)
			}
		}
		if (len(args) > 1) && (args[1] != "now") {
			end, err = time.Parse(timestamp, args[1])
			if err != nil {
				return 400, fmt.Errorf("Bad Request, %s", err)
			}
		}
		stmt := `SELECT eprintid FROM eprint WHERE (lastmod_year >= ? AND lastmod_year <= ?) AND (lastmod_month >= ? AND lastmod_month <= ?) AND (lastmod_day >= ? AND lastmod_day <= ?) AND (lastmod_hour >= ? AND lastmod_hour <= ?) AND (lastmod_minute >= ? AND lastmod_minute <=?) AND (eprint_status = 'deletion')`
		rows, err := db.Query(stmt, start.Year(), end.Year(), start.Month(), end.Month(), start.Day(), end.Day(), start.Hour(), end.Hour(), start.Minute(), end.Minute())
		if err != nil {
			log.Printf("ERROR: query error (%q), %s", repoID, err)
			return 500, fmt.Errorf("Internal Server Error")
		}
		defer rows.Close()
		eprintid := 0
		eprintIDs := []int{}
		for rows.Next() {
			err := rows.Scan(&eprintid)
			if (err == nil) && (eprintid > 0) {
				eprintIDs = append(eprintIDs, eprintid)
			} else {
				log.Printf("ERROR: scan error (%q), %s", repoID, err)
				return 500, fmt.Errorf("Internal Server Error")
			}
		}
		if err := rows.Err(); err != nil {
			log.Printf("ERROR: rows error (%q), %s", repoID, err)
			return 500, fmt.Errorf("Internal Server Error")
		}
		src, err := json.MarshalIndent(eprintIDs, "", "  ")
		if err != nil {
			log.Printf("ERROR: JSON marshal error (%q), %s", repoID, err)
			return 500, fmt.Errorf("Internal Server Error")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", src)
		return 200, nil
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func pubdateEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func doiEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func creatorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func creatorORCIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func editorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func contributorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func advisorIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func committeeIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func groupIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func funderIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

func grantNumberIDEndPoint(w http.ResponseWriter, r *http.Request, repoID string, args []string) (int, error) {
	if _, ok := config.Connections[repoID]; ok == true {
		return 501, fmt.Errorf("Not Implemented (%s) %s", repoID, r.URL.Path)
	} else {
		return 500, fmt.Errorf("Internal Server Error")
	}
}

//
// The following define the API as a service handling errors,
// routes and logging.
//
func logRequest(r *http.Request, status int, err error) {
	q := r.URL.Query()
	errStr := "OK"
	if err != nil {
		errStr = fmt.Sprintf("%s", err)
	}
	if len(q) > 0 {
		log.Printf("%s %s RemoteAddr: %s UserAgent: %s Query: %+v Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, errStr)
	} else {
		log.Printf("%s %s RemoteAddr: %s UserAgent: %s Response: %d %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, errStr)
	}
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, `ERROR: %d %s`, statusCode, err)
}

func routeEndPoints(w http.ResponseWriter, r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		return 400, fmt.Errorf("Bad Request")
	}
	// Expected URL structure of `/<REPO_ID>/<END_POINT>/<ARGS>`
	repoID, endPoint, args := parts[1], parts[2], parts[3:]
	if routes, hasRepo := config.Routes[repoID]; hasRepo == true {
		if fn, hasRoute := routes[endPoint]; hasRoute == true {
			return fn(w, r, repoID, args)
		}
	}
	return 404, fmt.Errorf("Not Found")
}

func api(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		statusCode int
	)
	if r.Method != "GET" {
		statusCode, err = 405, fmt.Errorf("Method Not Allowed")
		handleError(w, statusCode, err)
	} else {
		switch r.URL.Path {
		case "/":
			statusCode, err = 200, nil
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, `EPrints 3.3.x extended API, eprinttools version %s`, Version)
		case "/index.html":
			statusCode, err = 200, nil
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, `EPrints 3.3.x extended API, eprinttools version %s`, Version)
		case "/favicon.ico":
			statusCode, err = 200, nil
			fmt.Fprintf(w, "")
			//statusCode, err = 404, fmt.Errorf("Not Found")
			//handleError(w, statusCode, err)
		default:
			statusCode, err = routeEndPoints(w, r)
			if err != nil {
				handleError(w, statusCode, err)
			}
		}
	}
	logRequest(r, statusCode, err)
}

func loadConfig(fname string) error {
	config = new(Config)
	config.Repositories = map[string]string{}
	if src, err := ioutil.ReadFile(fname); err != nil {
		return err
	} else {
		// Since we should be OK, unmarshal in into active config
		if err = json.Unmarshal(src, &config); err != nil {
			return fmt.Errorf("Unmarshaling %q failed, %s", fname, err)
		}
		if config.Hostname == "" {
			config.Hostname = "localhost:8484"
		}
	}
	return nil
}

func InitExtendedAPI(settings string) error {
	var err error
	if err = loadConfig(settings); err != nil {
		return fmt.Errorf("Failed to load %q, %s", settings, err)
	}
	if config == nil {
		return fmt.Errorf("Missing configuration")
	}
	if config.Hostname == "" {
		return fmt.Errorf("Hostings hostname for service")
	}
	if config.Repositories == nil || len(config.Repositories) < 1 {
		return fmt.Errorf(`Missing "repositories" configuration`)
	}
	if config.Connections == nil {
		config.Connections = map[string]*sql.DB{}
	}
	if config.Routes == nil {
		config.Routes = map[string]map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){}
	}
	// This is a map endpoints and point handlers.
	// This implements the registration design pattern
	routes := map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){
		"updated":        updatedEndPoint,
		"deleted":        deletedEndPoint,
		"pubdate":        pubdateEndPoint,
		"doi":            doiEndPoint,
		"creator-id":     creatorIDEndPoint,
		"creator-orcid":  creatorORCIDEndPoint,
		"editor-id":      editorIDEndPoint,
		"contributor-id": contributorIDEndPoint,
		"advisor-id":     advisorIDEndPoint,
		"committee-id":   committeeIDEndPoint,
		"group-id":       groupIDEndPoint,
		"funder-id":      funderIDEndPoint,
		"grant-number":   grantNumberIDEndPoint,
	}

	/* NOTE: We need a DB connection to MySQL for each
	   EPrints repository supported by the API
	   for access to MySQL */
	for repoID, dataSourceName := range config.Repositories {
		// Setup DB connection for target repository
		if db, err := sql.Open("mysql", dataSourceName); err != nil {
			return fmt.Errorf("Could not open MySQL conncetion for %s, %s", repoID, err)
		} else {
			//log.Printf("Setting  DB connection to %q", repoID)
			//db.Ping()
			config.Connections[repoID] = db
		}
		// Add routes (end points) for the target repository
		for route, fn := range routes {
			if config.Routes[repoID] == nil {
				config.Routes[repoID] = map[string]func(http.ResponseWriter, *http.Request, string, []string) (int, error){}
			}
			config.Routes[repoID][route] = fn
		}
	}
	return nil
}

func RunExtendedAPI(appName string) error {
	/* Setup web server */
	log.Printf(`
%s %s

EPrints 3.3.x Extended API

Listening on http://%s

Press ctl-c to terminate.
`, appName, Version, config.Hostname)
	http.HandleFunc("/", api)
	return http.ListenAndServe(config.Hostname, nil)
}

package eprinttools

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/dataset/v2"
)

//
//
// The "datasets" creates the pairtree datasets based on the
// EPrinttools settings.json and previously run harvests.
//

// generateDataset creates a dataset from the harvested repository.
func generateDataset(cfg *Config, repoName string, dsn string, projectDir string, keyList string, verbose bool) error {
	repoCfg, ok := cfg.Repositories[repoName]
	if !ok {
		return fmt.Errorf("%s not found in configuration", repoName)
	}

	cName := path.Join(projectDir, repoCfg.DefaultCollection) + ".ds"
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	// NOTE: We're creating a pairtree dataset collection so it can be downloaded as a ZIP
	c, err := dataset.Init(cName, dsn)
	if err != nil {
		return err
	}
	c.Close()
	c, err = dataset.Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()

	if keyList == "" {
		// Get all the rows from the harvested repository, then
		// write the JSON source to the dataset.
		var (
			cntStmt string
			stmt    string
			tot     int
		)
		if repoCfg.PublicOnly {
			cntStmt = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE is_public IS TRUE", repoName)
			stmt = fmt.Sprintf("SELECT id, src FROM %s WHERE is_public IS TRUE", repoName)
		} else {
			cntStmt = fmt.Sprintf("SELECT COUNT(*) FROM %s", repoName)
			stmt = fmt.Sprintf("SELECT id, src FROM %s", repoName)
		}
		rows, err := cfg.Jdb.Query(cntStmt)
		if err != nil {
			return err
		}
		if rows.Next() {
			if err := rows.Scan(&tot); err != nil {
				return err
			}
		}
		rows.Close()

		modValue := calcModValue(tot)

		rows, err = cfg.Jdb.Query(stmt)
		if err != nil {
			return err
		}
		defer rows.Close()
		i := 0
		t0 := time.Now()
		for rows.Next() {
			var (
				eprintid int
				src      []byte
			)
			if err := rows.Scan(&eprintid, &src); err != nil {
				log.Printf("failed to read row in %q, %s", repoName, err)
			} else {
				key := fmt.Sprintf("%d", eprintid)
				if err := c.CreateJSON(key, src); err != nil {
					log.Printf("failed to create %q in %q, %s", key, cName, err)
				}
				i += 1
				if verbose && ((i % modValue) == 0) {
					log.Printf("%s %d records added (%s)", path.Base(cName), i, progress(t0, i, tot))
				}
			}
		}
		if err := rows.Err(); err != nil {
			log.Printf("row error %s", err)
			return err
		}
		if verbose {
			log.Printf("dataset %s created in %v", cName, time.Since(t0).Truncate(time.Second))
		}
	} else {
		// Get key list rows from the harvested repository, then
		// write the JSON source to the dataset.
		var (
			stmt string
			tot  int
			ids  []int
		)
		fp, err := os.Open(keyList)
		if err != nil {
			return err
		}
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			tok := scanner.Text()
			if tok != "" {
				id, err := strconv.Atoi(strings.TrimSpace(tok))
				if err != nil {
					log.Printf("warning %q is not a valid eprintid", tok)
				} else {
					ids = append(ids, id)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		tot = len(ids)
		if verbose {
			log.Printf("Retrieved %d keys from %q", tot, keyList)
		}
		modValue := calcModValue(tot)

		stmt = fmt.Sprintf("SELECT id, src FROM %s WHERE id = ?", repoName)
		t0 := time.Now()
		for i, id := range ids {
			rows, err := cfg.Jdb.Query(stmt, id)
			if err != nil {
				log.Printf("WARNING query failed for %d, line %d, %q", id, i, err)
				continue
			}
			defer rows.Close()
			for rows.Next() {
				var (
					eprintid int
					src      []byte
				)
				if err := rows.Scan(&eprintid, &src); err != nil {
					log.Printf("WARNING failed to read row in %q, id %d, line %d, %s", repoName, id, i, err)
				} else {
					key := fmt.Sprintf("%d", eprintid)
					if err := c.CreateJSON(key, src); err != nil {
						log.Printf("failed to create %q in %q, %s", key, cName, err)
					}
				}
			}
			if err := rows.Err(); err != nil {
				log.Printf("WARNING row error for id %d, line %d, %s", id, i, err)
			}
			if verbose && ((i % modValue) == 0) {
				log.Printf("%s %d records added (%s)", path.Base(cName), i, progress(t0, i, tot))
			}
		}
		if verbose {
			log.Printf("dataset %s created in %v", cName, time.Since(t0).Truncate(time.Second))
		}
	}
	return nil
}

// RunDataset will use the eprinttools settings.jons config file
// and a repository ID (e.g. caltechauthors) and render a
// dataset collection based on the previously harvested contents.
func RunDataset(cfgName string, repoName string, dsn string, keyList string, verbose bool) error {
	t0 := time.Now()
	appName := path.Base(os.Args[0])
	// Read in the configuration for this harvester instance.
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("could not create a configuration object")
	}
	if err := OpenJSONStore(cfg); err != nil {
		return err
	}
	defer cfg.Jdb.Close()
	if cfg.ProjectDir != "" && cfg.ProjectDir != "." {
		if _, err := os.Stat(cfg.ProjectDir); os.IsNotExist(err) {
			if err := os.MkdirAll(cfg.ProjectDir, 0775); err != nil {
				log.Printf("%s", err)
				os.Exit(1)
			}
		}
	}
	if verbose {
		log.Printf("%s (%s) started %v", appName, repoName, time.Now().Sub(t0).Truncate(time.Second))
	}
	if err := generateDataset(cfg, repoName, dsn, cfg.ProjectDir, keyList, verbose); err != nil {
		return err
	}
	if verbose {
		log.Printf("%s (%s) run time %v", appName, repoName, time.Since(t0).Truncate(time.Second))
	}
	return nil
}

// RunDatasets will use the eprinttools settings.jons config file
// and reader dataset collections based on the contents.
func RunDatasets(cfgName string, dsn string, keyList string, verbose bool) error {
	t0 := time.Now()
	appName := path.Base(os.Args[0])
	// Read in the configuration for this harvester instance.
	cfg, err := LoadConfig(cfgName)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("could not create a configuration object")
	}
	if err := OpenJSONStore(cfg); err != nil {
		return err
	}
	defer cfg.Jdb.Close()
	if _, err := os.Stat(cfg.ProjectDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.ProjectDir, 0775); err != nil {
			log.Printf("%s", err)
			os.Exit(1)
		}
	}
	if verbose {
		log.Printf("%s started %v", appName, time.Now().Sub(t0).Truncate(time.Second))
	}
	for repoName := range cfg.Repositories {
		if err := generateDataset(cfg, repoName, dsn, cfg.ProjectDir, keyList, verbose); err != nil {
			return err
		}
	}
	if verbose {
		log.Printf("%s run time %v", appName, time.Since(t0).Truncate(time.Second))
	}
	return nil
}

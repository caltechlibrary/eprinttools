//
// Package ep is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"unicode/utf8"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/ep"
)

var (
	usage = `USAGE: %s [OPTIONS]`

	description = `
SYNOPSIS

%s generates JSON documents (feeds) in a htdocs directory tree.

CONFIGURATION

%s can be configured through setting the following environment
variables-

EP_DATASET this is the dataset and collection directory (e.g. dataset/eprints)

EP_HTDOCS  this is the directory where the JSON documents will be written.`

	examples = `
EXAMPLE

    %s 

Generates JSON documents in EP_HTDOCS from EP_DATASET.`

	// Standard Options
	showHelp    bool
	showVersion bool
	showLicense bool
	outputFName string

	// App Options
	htdocs      string
	datasetName string
)

// slugify ensures we have a path friendly name or returns an error.
// NOTE: The web server does not expect to look on disc for URL Encoded paths, instead
// we need to ensure the name does not have a slash or other path unfriendly value.
func slugify(s string) (string, error) {
	if utf8.RuneCountInString(s) > 200 {
		return "", fmt.Errorf("string to long (%d), %q", utf8.RuneCountInString(s), s)
	}
	if strings.Contains(s, "/") == true {
		return "", fmt.Errorf("string contains a slash and cannot be a directory name, %q", s)
	}
	if strings.Contains(s, `\`) == true {
		return "", fmt.Errorf("string contains a back slash and should be a directory name, %q", s)
	}
	return s, nil
}

// buildSite generates a website based on the contents of the exported EPrints data.
// The site builder needs to know the name of the BoltDB, the root directory
// for the website and directory to find the templates
func buildSite(api *ep.EPrintsAPI, feedSize int) error {
	var err error

	if feedSize < 1 {
		feedSize = ep.DefaultFeedSize
	}

	// FIXME: This could be replaced by copying all the records in dataset/COLLECTION
	// that are public and published.

	// Collect the recent publications (all types)
	log.Printf("Building Recently Published (feed size %d)", feedSize)
	err = api.BuildFeed(feedSize, "Recently Published", path.Join("recent", "publications"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
		return api.GetPublications(0, feedSize)
	})
	if err != nil {
		log.Printf("error: %s", err)
	}
	// Collect the rencently published  articles
	log.Printf("Building Recent Articles")
	err = api.BuildFeed(feedSize, "Recent Articles", path.Join("recent", "articles"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
		return api.GetArticles(0, feedSize)
	})
	if err != nil {
		log.Printf("error: %s", err)
	}

	// Collect EPrints by Group/Affiliation
	log.Printf("Building Local Groups")
	groupNames, err := api.GetLocalGroups()
	if err != nil {
		log.Printf("error: %s", err)
	} else {
		log.Printf("Found %d groups\n", len(groupNames))
		for _, groupName := range groupNames {
			// Build recently for each affiliation
			slug, err := slugify(groupName)
			if err != nil {
				log.Printf("Skipping %q, %s\n", groupName, err)
			} else {
				// Build complete list for each affiliation
				err = api.BuildFeed(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "publications"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
					return api.GetLocalGroupPublications(groupName, 0, -1)
				})
				if err != nil {
					log.Printf("Skipped: %s", err)
				} else {
					err = api.BuildFeed(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "recent", "publications"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
						return api.GetLocalGroupPublications(groupName, start, count)
					})
					if err != nil {
						log.Printf("Skipped: %s", err)
					}
					// Build complete list of articles for each affiliation
					err = api.BuildFeed(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "articles"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
						return api.GetLocalGroupArticles(groupName, 0, -1)
					})
					if err != nil {
						log.Printf("Skipped: %s", err)
					} else {
						// Build recent articles for each affiliation
						err = api.BuildFeed(-1, fmt.Sprintf("%s", groupName), path.Join("affiliation", slug, "recent", "articles"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
							return api.GetLocalGroupArticles(groupName, start, count)
						})
						if err != nil {
							log.Printf("Skipped: %s", err)
						}
					}
				}
			}
		}
	}

	/*
		// Collect EPrints by Funders
		log.Printf("Building Funders")
		funderNames, err := api.GetFunders()
		if err != nil {
			log.Printf("error: %s", err)
		} else {
			log.Printf("Found %d records with funders\n", len(funderNames))
			for _, funderName := range funderNames {
				slug, err := slugify(funderName)
				if err != nil {
					log.Printf("Skipping %q, %s\n", funderName, err)
				} else {
					// Build complete list for each funder
					err = api.BuildFeed(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "publications"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
						return api.GetFunderPublications(funderName, 0, -1)
					})
					if err != nil {
						log.Printf("Skipped: %s", err)
					} else {
						// Build recently for each funder
						err = api.BuildFeed(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "recent", "publications"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
							return api.GetFunderPublications(funderName, start, count)
						})
						if err != nil {
							log.Printf("Skipped: %s", err)
						}
						// Build complete list of articles for each funder
						err = api.BuildFeed(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "articles"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
							return api.GetFunderArticles(funderName, 0, -1)
						})
						if err != nil {
							log.Printf("Skipped: %s", err)
						} else {
							// Build recent articles for each funder
							err = api.BuildFeed(-1, fmt.Sprintf("%s", funderName), path.Join("funder", slug, "recent", "articles"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
								return api.GetFunderArticles(funderName, start, count)
							})
							if err != nil {
								log.Printf("Skipped: %s", err)
							}
						}
					}
				}
			}
		}
	*/

	// Collect EPrints by orcid ID and publish
	log.Printf("Building Person (orcid) works")
	orcids, err := api.GetORCIDs()
	if err != nil {
		log.Printf("error: %s", err)
	} else {
		log.Printf("Found %d orcids\n", len(orcids))
		for _, orcid := range orcids {

			// Build complete list for each orcid
			err = api.BuildFeed(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "publications"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
				return api.GetORCIDPublications(orcid, 0, -1)
			})
			if err != nil {
				log.Printf("Skipped: %s", err)
			} else {
				// Build a list of recent ORCID Publications
				err = api.BuildFeed(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "recent", "publications"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
					return api.GetORCIDPublications(orcid, start, count)
				})
				if err != nil {
					log.Printf("Skipped: %s", err)
				}
				// Build complete list of articles for each ORCID
				err = api.BuildFeed(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "articles"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
					return api.GetORCIDArticles(orcid, 0, -1)
				})
				if err != nil {
					log.Printf("Skipped: %s", err)
				} else {
					// Build a list of recent ORCID Articles
					err = api.BuildFeed(-1, fmt.Sprintf("ORCID: %s", orcid), path.Join("person", fmt.Sprintf("%s", orcid), "recent", "articles"), func(api *ep.EPrintsAPI, start, count int) ([]*ep.Record, error) {
						return api.GetORCIDArticles(orcid, start, count)
					})
					if err != nil {
						log.Printf("Skipped: %s", err)
					}
				}
			}
		}
	}

	return nil
}

func check(cfg *cli.Config, key, value string) string {
	if value == "" {
		log.Fatalf("Missing %s_%s", cfg.EnvPrefix, strings.ToUpper(key))
		return ""
	}
	return value
}

func init() {
	// Setup options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.StringVar(&outputFName, "o", "", "output filename (log message)")
	flag.StringVar(&outputFName, "output", "", "output filename (log message)")

	// App Specific options
	// NOTE: htdocs uses "d" and "docs" to align with mkpage option practice
	flag.StringVar(&htdocs, "d", "", "specify where to write the HTML files to")
	flag.StringVar(&htdocs, "docs", "", "specify where to write the HTML files to")

	// NOTE: "d" is taken so I am only including a long form
	flag.StringVar(&datasetName, "dataset", "", "the dataset/collection name")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	cfg := cli.New(appName, "EP", fmt.Sprintf(ep.LicenseText, appName, ep.Version), ep.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName)

	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	out, err := cli.Create(outputFName, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer cli.CloseFile(outputFName, out)

	// Log to out
	log.SetOutput(out)

	// Check to see we can merge the required fields are merged.
	htdocs = check(cfg, "htdocs", cfg.MergeEnv("htdocs", htdocs))
	datasetName = check(cfg, "dataset", cfg.MergeEnv("dataset", datasetName))

	if htdocs != "" {
		if _, err := os.Stat(htdocs); os.IsNotExist(err) {
			os.MkdirAll(htdocs, 0775)
		}
	}

	// Create an API instance
	api, err := ep.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//
	// Read the dataset indicated in configuration and
	// render pages in the various formats supported.
	//
	t0 := time.Now()
	log.Printf("%s %s\n", appName, ep.Version)
	log.Printf("Rendering pages from %s\n", datasetName)
	err = buildSite(api, -1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rendering complete")
	t1 := time.Now()
	log.Printf("Running time %v", t1.Sub(t0))
}

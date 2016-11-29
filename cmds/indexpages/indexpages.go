//
// Package epgo is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	// Caltech Libraries packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/epgo"

	// 3rd Party packages
	"github.com/blevesearch/bleve"
)

var (
	usage = `USAGE: %s [OPTIONS] [BLEVE_INDEX_NAME]`

	description = `
 SYNOPSIS

 %s is a command line utility to indexes content in the htdocs directory.
 It produces a Bleve search index used by servepages web service.
 Configuration is done through environmental variables.

 CONFIGURATION

 %s relies on the following environment variables for
 configuration when overriding the defaults:

    EPGO_HTDOCS       This should be the path to the directory tree
                        containings the content (e.g. JSON files) to be index.
                        This is generally populated with the caitpage command.
						Defaults to ./htdocs.

    EPGO_BLEVE        This is is the directory that will contain all the Bleve
                        indexes. Defaults to ./htdocs.bleve

`

	license = `
%s %s

Copyright (c) 2016, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of epgo nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

`
	// cli options
	showHelp    bool
	showVersion bool
	showLicense bool

	// additional options
	replaceIndex   bool
	htdocs         string
	indexName      string
	dbName         string
	apiURL         string
	siteURL        string
	templatePath   string
	repositoryPath string
	batchSize      int

	// internal counters
	dirCount  int
	fileCount int
)

func handleSignals() {
	signalChannel := make(chan os.Signal, 3)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			//handle SIGINT
			log.Println("SIGINT received, shutting down")
			os.Exit(0)
		case syscall.SIGTERM:
			//handle SIGTERM
			log.Println("SIGTERM received, shutting down")
			os.Exit(0)
		case syscall.SIGHUP:
			//FIXME: this maybe a good choice for closing and re-opening the index with bringing down the web service
			log.Println("SIGHUP received, shutting down")
			os.Exit(0)
		}
	}()
}

func check(cfg *cli.Config, key, value string) string {
	if value == "" {
		log.Fatal("Missing %s_%s", cfg.EnvPrefix, strings.ToUpper(key))
		return ""
	}
	return value
}

func init() {
	// Log to standard out
	log.SetOutput(os.Stdout)

	// Standard options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")

	// App specific options
	flag.StringVar(&htdocs, "htdocs", "", "The document root for the website")
	flag.StringVar(&indexName, "bleve", "", "The name of the Bleve index")
	flag.BoolVar(&replaceIndex, "r", false, "Replace the index if it exists")
	flag.StringVar(&repositoryPath, "repository-path", "", "Path of rendered repository content")
	flag.IntVar(&batchSize, "batch", batchSize, "Set the batch index size")
}

func createIndex(indexName string) (bleve.Index, error) {
	log.Printf("Creating Bleve index at %s\n", indexName)

	log.Println("Setting up index...")
	indexMapping := bleve.NewIndexMapping()
	// Add EPrint as a specific document map
	eprintMapping := bleve.NewDocumentMapping()

	/*
		EPrintID:     jsonDoc.ID,
		Type:         jsonDoc.Type,
		OfficialURL:  jsonDoc.OfficialURL,
		Title:        jsonDoc.Title,
		Abstract:     jsonDoc.Abstract,
		Keywords:     jsonDoc.Keywords,
		ISSN:         jsonDoc.ISSN,
		Publication:  jsonDoc.Publication,
		Note:         jsonDoc.Note,
		Authors:      jsonDoc.Creators.ToNames(),
		ORCIDs:       jsonDoc.Creators.ToORCIDs(),
		ISNIs:        jsonDoc.Creators.ToISNIs(),
		Rights:       jsonDoc.Rights,
		Funders:      jsonDoc.Funders.ToAgencies(),
		GrantNumbers: jsonDoc.Funders.ToGrantNumbers(),
		PubDate:      jsonDoc.PubDate(),
		LocalGroup:  jsonDoc.LocalGroup,
	*/
	// Now add specific eprint fields
	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = "en"
	titleMapping.Store = true
	titleMapping.Index = true
	eprintMapping.AddFieldMappingsAt("Title", titleMapping)

	abstractMapping := bleve.NewTextFieldMapping()
	abstractMapping.Analyzer = "en"
	abstractMapping.Store = true
	abstractMapping.Index = true
	eprintMapping.AddFieldMappingsAt("Abstract", abstractMapping)

	publicationMapping := bleve.NewTextFieldMapping()
	publicationMapping.Analyzer = "en"
	publicationMapping.Store = true
	publicationMapping.Index = true
	eprintMapping.AddFieldMappingsAt("Publication", publicationMapping)

	subjectsMapping := bleve.NewTextFieldMapping()
	subjectsMapping.Analyzer = "en"
	subjectsMapping.Store = true
	subjectsMapping.Index = true
	subjectsMapping.IncludeTermVectors = true
	eprintMapping.AddFieldMappingsAt("Subjects", subjectsMapping)

	keywordsMapping := bleve.NewTextFieldMapping()
	keywordsMapping.Analyzer = "en"
	keywordsMapping.Store = true
	keywordsMapping.Index = true
	keywordsMapping.IncludeTermVectors = true
	eprintMapping.AddFieldMappingsAt("Keywords", keywordsMapping)

	typeMapping := bleve.NewTextFieldMapping()
	typeMapping.Analyzer = "en"
	typeMapping.Store = true
	typeMapping.Index = true
	eprintMapping.AddFieldMappingsAt("Type", typeMapping)

	localGroupMapping := bleve.NewTextFieldMapping()
	localGroupMapping.Analyzer = "en"
	localGroupMapping.Store = true
	localGroupMapping.Index = true
	eprintMapping.AddFieldMappingsAt("LocalGroup", localGroupMapping)

	fundersMapping := bleve.NewTextFieldMapping()
	fundersMapping.Analyzer = "en"
	eprintMapping.AddFieldMappingsAt("Funders", fundersMapping)

	creatorsMapping := bleve.NewTextFieldMapping()
	creatorsMapping.Analyzer = "en"
	creatorsMapping.Store = true
	creatorsMapping.Index = true
	creatorsMapping.IncludeTermVectors = true
	eprintMapping.AddFieldMappingsAt("Authors", creatorsMapping)

	orcidMapping := bleve.NewTextFieldMapping()
	orcidMapping.Analyzer = "en"
	orcidMapping.Store = true
	orcidMapping.Index = true
	eprintMapping.AddFieldMappingsAt("ORCIDs", orcidMapping)

	isniMapping := bleve.NewTextFieldMapping()
	isniMapping.Analyzer = "en"
	isniMapping.Store = true
	isniMapping.Index = true
	eprintMapping.AddFieldMappingsAt("ISNIs", isniMapping)

	createdMapping := bleve.NewDateTimeFieldMapping()
	createdMapping.Store = true
	createdMapping.Index = false
	eprintMapping.AddFieldMappingsAt("created", createdMapping)

	// Finally add this mapping to the main index mapping
	indexMapping.AddDocumentMapping("eprint", eprintMapping)

	log.Printf("Creating Bleve index at %s\n", indexName)
	index, err := bleve.New(indexName, indexMapping)
	if err != nil {
		return nil, fmt.Errorf("Can't create new bleve index %q, %s", indexName, err)
	}
	return index, nil
}

func getIndex(indexName string) (bleve.Index, error) {
	//FIXME: we want to create a new fresh index, then swap the alias to the old one
	if _, err := os.Stat(indexName); os.IsNotExist(err) {
		return createIndex(indexName)
	}
	log.Printf("Opening Bleve index at %s\n", indexName)
	index, err := bleve.Open(indexName)
	if err != nil {
		return nil, fmt.Errorf("Can't open new bleve index %q, %s", indexName, err)
	}
	return index, nil
}

func indexSite(htdocs, eprintsDotJSON string, index bleve.Index, batchSize int) error {
	var (
		uris []string
	)

	startT := time.Now()
	count := 0
	batchNo := 1

	log.Printf("Reading %s", eprintsDotJSON)
	src, err := ioutil.ReadFile(eprintsDotJSON)
	if err != nil {
		return err
	}

	log.Printf("Decoding %s", eprintsDotJSON)
	err = json.Unmarshal(src, &uris)
	if err != nil {
		return err
	}
	total := len(uris)

	log.Printf("%d eprints found", total)
	batch := index.NewBatch()
	log.Printf("Indexed: %d of %d, batch size %d, run time %s", count, total, batchSize, time.Now().Sub(startT))
	for _, uri := range uris {
		p := path.Join(htdocs, uri)
		src, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		jsonDoc := new(epgo.Record)
		err = json.Unmarshal(src, &jsonDoc)
		if err != nil {
			log.Printf("error indexing %s, %s", uri, err)
		} else {
			batch.Index(uri, struct {
				EPrintID     int
				OfficialURL  string
				Title        string
				Abstract     string
				Keywords     string
				ISSN         string
				Publication  string
				Note         string
				Authors      []string
				ORCIDs       []string
				ISNIs        []string
				Type         string
				Rights       string
				Funders      []string
				GrantNumbers []string
				PubDate      string
				LocalGroup   []string
				Subjects     []string
			}{
				EPrintID:     jsonDoc.ID,
				Type:         jsonDoc.Type,
				OfficialURL:  jsonDoc.OfficialURL,
				Title:        jsonDoc.Title,
				Abstract:     jsonDoc.Abstract,
				Keywords:     jsonDoc.Keywords,
				ISSN:         jsonDoc.ISSN,
				Publication:  jsonDoc.Publication,
				Note:         jsonDoc.Note,
				Authors:      jsonDoc.Creators.ToNames(),
				ORCIDs:       jsonDoc.Creators.ToORCIDs(),
				ISNIs:        jsonDoc.Creators.ToISNIs(),
				Rights:       jsonDoc.Rights,
				Funders:      jsonDoc.Funders.ToAgencies(),
				GrantNumbers: jsonDoc.Funders.ToGrantNumbers(),
				PubDate:      jsonDoc.PubDate(),
				LocalGroup:   jsonDoc.LocalGroup,
				Subjects:     jsonDoc.Subjects,
			})
			// Flag the memory used by this jsonDoc to be garbage collected
			jsonDoc = nil
		}
		if batch.Size() >= batchSize {
			log.Printf("Indexing batch %d", batchNo)
			batchNo++
			err := index.Batch(batch)
			if err != nil {
				return err
			}
			count += batch.Size()
			batch.Reset()
			log.Printf("Indexed: %d of %d, batch size %d, run time %s", count, total, batchSize, time.Now().Sub(startT))
		}
	}
	if batch.Size() > 0 {
		log.Printf("Indexing batch %d", batchNo)
		err := index.Batch(batch)
		if err != nil {
			return err
		}
		count += batch.Size()
		log.Printf("Indexed: %d of %d, batch size %d, run time %s", count, total, batchSize, time.Now().Sub(startT))
	}
	return nil
}

func main() {
	var err error

	appName := path.Base(os.Args[0])
	cfg := cli.New(appName, "EPGO", fmt.Sprintf(license, appName, epgo.Version), epgo.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.OptionsText = "OPTIONS\n"

	flag.Parse()
	args := flag.Args()
	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}

	if batchSize == 0 {
		batchSize = 500
	}

	if len(args) > 0 {
		indexName = strings.Join(args, ":")
	}

	// Required fields
	dbName = check(cfg, "dbname", cfg.MergeEnv("dbname", dbName))
	names := check(cfg, "bleve", cfg.MergeEnv("bleve", indexName))
	if strings.Contains(names, ":") {
		indexName = (strings.Split(names, ":"))[0]
	}
	htdocs = check(cfg, "htdocs", cfg.MergeEnv("htdocs", htdocs))
	siteURL = check(cfg, "site_url", cfg.MergeEnv("site_url", siteURL))
	repositoryPath = check(cfg, "repository_path", cfg.MergeEnv("repository_path", repositoryPath))

	// Optional fields
	apiURL = cfg.MergeEnv("api_url", apiURL)
	templatePath = cfg.MergeEnv("template_path", templatePath)

	// Now log what we're running
	log.Printf("%s %s", appName, epgo.Version)

	if replaceIndex == true {
		log.Printf("Clearing index")
		err := os.RemoveAll(indexName)
		if err != nil {
			log.Fatalf("Could not removed %q, %s", indexName, err)
		}
	}

	handleSignals()

	index, err := getIndex(indexName)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	// Walk our data import tree and index things
	log.Printf("Start indexing contents of %s as %s\n", path.Join(htdocs, repositoryPath, "eprints.json"), indexName)
	err = indexSite(htdocs, path.Join(htdocs, repositoryPath, "eprints.json"), index, batchSize)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Finished")
}

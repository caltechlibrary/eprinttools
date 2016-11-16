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
	"path"
	"time"

	// Caltech Libraries packages
	"github.com/caltechlibrary/epgo"

	// 3rd Party packages
	"github.com/blevesearch/bleve"
)

var (
	description = `
 USAGE: %s [OPTIONS]

 SYNOPSIS

 %s is a command line utility to indexes content in the htdocs directory.
 It produces a Bleve search index used by servepages web service.
 Configuration is done through environmental variables.

 OPTIONS
`

	configuration = `

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
	maxBatchSize   int

	// internal counters
	dirCount  int
	fileCount int
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func usage(appName, version string) {
	fmt.Printf(description, appName, appName)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Printf(configuration, appName)
	fmt.Printf("%s %s\n", appName, version)
	os.Exit(0)
}

func init() {
	// standard options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")

	// app specific options
	flag.StringVar(&htdocs, "htdocs", "", "The document root for the website")
	flag.StringVar(&indexName, "bleve", "", "The name of the Bleve index")
	flag.BoolVar(&replaceIndex, "r", false, "Replace the index if it exists")
	flag.StringVar(&repositoryPath, "repository-path", "", "Path of rendered repository content")
	flag.IntVar(&maxBatchSize, "batch", maxBatchSize, "Set the maximum batch index size")
}

func createIndex(indexName string) (bleve.Index, error) {
	log.Printf("Creating Bleve index at %s\n", indexName)

	log.Println("Setting up index...")
	indexMapping := bleve.NewIndexMapping()
	// Add EPrint as a specific document map
	eprintMapping := bleve.NewDocumentMapping()

	// Now add specific eprint fields
	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = "en"
	titleMapping.Store = true
	titleMapping.Index = true
	eprintMapping.AddFieldMappingsAt("title", titleMapping)

	abstractMapping := bleve.NewTextFieldMapping()
	abstractMapping.Analyzer = "en"
	abstractMapping.Store = true
	abstractMapping.Index = true
	eprintMapping.AddFieldMappingsAt("abstract", abstractMapping)

	publicationMapping := bleve.NewTextFieldMapping()
	publicationMapping.Analyzer = "en"
	publicationMapping.Store = true
	publicationMapping.Index = true
	eprintMapping.AddFieldMappingsAt("publication", publicationMapping)

	subjectsMapping := bleve.NewTextFieldMapping()
	subjectsMapping.Analyzer = "en"
	subjectsMapping.Store = true
	subjectsMapping.Index = true
	subjectsMapping.IncludeTermVectors = true
	eprintMapping.AddFieldMappingsAt("subjects", subjectsMapping)

	typeMapping := bleve.NewTextFieldMapping()
	typeMapping.Analyzer = "en"
	typeMapping.Store = true
	typeMapping.Index = true
	eprintMapping.AddFieldMappingsAt("type", typeMapping)

	localGroupMapping := bleve.NewTextFieldMapping()
	localGroupMapping.Analyzer = "en"
	eprintMapping.AddFieldMappingsAt("local_group", localGroupMapping)

	fundersMapping := bleve.NewTextFieldMapping()
	fundersMapping.Analyzer = "en"
	eprintMapping.AddFieldMappingsAt("funders", fundersMapping)

	creatorsMapping := bleve.NewTextFieldMapping()
	creatorsMapping.Analyzer = "en"
	creatorsMapping.IncludeTermVectors = true
	eprintMapping.AddFieldMappingsAt("creators", creatorsMapping)

	orcidMapping := bleve.NewTextFieldMapping()
	orcidMapping.Analyzer = "en"
	eprintMapping.AddFieldMappingsAt("orcid", orcidMapping)

	isniMapping := bleve.NewTextFieldMapping()
	isniMapping.Analyzer = "en"
	eprintMapping.AddFieldMappingsAt("isni", isniMapping)

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

func indexSite(htdocs, eprintsDotJSON string, index bleve.Index, maxBatchSize int) error {
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
	batchSize := 10
	batch := index.NewBatch()
	log.Printf("Indexed: %d of %d, batch size %d, run time %s", count, total, batchSize, time.Now().Sub(startT))
	for _, uri := range uris {
		p := path.Join(htdocs, uri)
		src, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		var jsonDoc interface{}
		err = json.Unmarshal(src, &jsonDoc)
		if err != nil {
			log.Printf("error indexing %s, %s", uri, err)
		} else {
			batch.Index(uri, jsonDoc)
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
			if batchSize < maxBatchSize {
				batchSize += 10
			}
			if batchSize > maxBatchSize {
				batchSize = maxBatchSize
			}
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

	flag.Parse()
	if showHelp == true {
		usage(appName, epgo.Version)
	}
	if showVersion == true {
		fmt.Printf("%s %s\n", appName, epgo.Version)
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Printf(license, appName, epgo.Version)
	}

	if maxBatchSize == 0 {
		maxBatchSize = 100
	}

	var cfg epgo.Config

	// Required fields
	check(cfg.MergeEnv("EPGO", "DBNAME", dbName))
	check(cfg.MergeEnv("EPGO", "BLEVE", indexName))
	check(cfg.MergeEnv("EPGO", "HTDOCS", htdocs))
	check(cfg.MergeEnv("EPGO", "SITE_URL", siteURL))
	check(cfg.MergeEnv("EPGO", "REPOSITORY_PATH", repositoryPath))
	// Optional fields
	cfg.MergeEnv("EPGO", "API_URL", apiURL)
	cfg.MergeEnv("EPGO", "TEMPLATE_PATH", templatePath)

	// Now log what we're running
	log.Printf("%s %s", appName, epgo.Version)

	if replaceIndex == true {
		log.Printf("Clearing index")
		err := os.RemoveAll(cfg.Get("bleve"))
		if err != nil {
			log.Fatalf("Could not removed %s, %s", cfg.Get("bleve"), err)
		}
	}

	index, err := getIndex(cfg.Get("bleve"))
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer index.Close()

	// Walk our data import tree and index things
	log.Printf("Start indexing contents of %s as %s\n", path.Join(cfg.Get("htdocs"), cfg.Get("repository_path"), "eprints.json"), cfg.Get("bleve"))
	err = indexSite(cfg.Get("htdocs"), path.Join(cfg.Get("htdocs"), cfg.Get("repository_path"), "eprints.json"), index, maxBatchSize)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Finished")
}

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
	"flag"
	"fmt"
	"log"
	"os"
	"path"

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
	replaceIndex bool
	htdocs       string
	indexName    string
	dbName       string
	apiURL       string
	siteURL      string
	templatePath string

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
}

func getIndex(indexName string) (bleve.Index, error) {
	if _, err := os.Stat(indexName); os.IsNotExist(err) {
		log.Printf("Creating Bleve index at %s\n", indexName)

		log.Println("Setting up index...")
		indexMapping := bleve.NewIndexMapping()
		// Add Accession as a specific document map
		eprintMapping := bleve.NewDocumentMapping()

		// Now add specific accession fields
		titleMapping := bleve.NewTextFieldMapping()
		titleMapping.Analyzer = "en"
		titleMapping.Store = true
		titleMapping.Index = true
		eprintMapping.AddFieldMappingsAt("title", titleMapping)

		descriptionMapping := bleve.NewTextFieldMapping()
		descriptionMapping.Analyzer = "en"
		descriptionMapping.Store = true
		descriptionMapping.Index = true
		eprintMapping.AddFieldMappingsAt("description", descriptionMapping)

		subjectsMapping := bleve.NewTextFieldMapping()
		subjectsMapping.Analyzer = "en"
		subjectsMapping.Store = true
		subjectsMapping.Index = true
		subjectsMapping.IncludeTermVectors = true
		eprintMapping.AddFieldMappingsAt("subject", subjectsMapping)

		datesMapping := bleve.NewTextFieldMapping()
		datesMapping.Store = true
		datesMapping.Index = false
		eprintMapping.AddFieldMappingsAt("date", datesMapping)

		createdMapping := bleve.NewDateTimeFieldMapping()
		createdMapping.Store = true
		createdMapping.Index = false
		eprintMapping.AddFieldMappingsAt("created", createdMapping)

		// Finally add this mapping to the main index mapping
		indexMapping.AddDocumentMapping("eprint", eprintMapping)

		index, err := bleve.New(indexName, indexMapping)
		if err != nil {
			return nil, fmt.Errorf("Can't create new bleve index %s, %s", indexName, err)
		}
		return index, nil
	}
	log.Printf("Opening Bleve index at %s\n", indexName)
	index, err := bleve.Open(indexName)
	if err != nil {
		return nil, fmt.Errorf("Can't create new bleve index %s, %s", indexName, err)
	}
	return index, nil
}

func indexSite(index bleve.Index, batchSize int) error {
	return fmt.Errorf("indexSite() not implemented.")
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

	var cfg epgo.Config

	// Required fields
	check(cfg.MergeEnv("EPGO", "DBNAME", dbName))
	check(cfg.MergeEnv("EPGO", "BLEVE", indexName))
	check(cfg.MergeEnv("EPGO", "HTDOCS", htdocs))
	check(cfg.MergeEnv("EPGO", "SITE_URL", siteURL))
	// Optional fields
	cfg.MergeEnv("EPGO", "API_URL", apiURL)
	cfg.MergeEnv("EPGO", "TEMPLATE_PATH", templatePath)

	if replaceIndex == true {
		err := os.RemoveAll(indexName)
		if err != nil {
			log.Fatalf("Could not removed %s, %s", indexName, err)
		}
	}

	index, err := getIndex(indexName)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer index.Close()

	// Walk our data import tree and index things
	log.Printf("Start indexing of %s in %s\n", htdocs, indexName)
	indexSite(index, 1000)
	log.Printf("Finished")
}

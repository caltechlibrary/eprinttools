package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	// Caltech Library packages
	"github.com/caltechlibrary/epgo"
)

var (
	description = `
 USAGE: %s [OPTIONS]

 OVERVIEW

	%s generates HTML, .include pages, BibTeX and normalized JSON based 
	on the JSON output form epgo and templates associated with 
	the command.

 OPTIONS
`
	configuration = `

 CONFIGURATION

    %s can be configured through setting the following environment
	variables-

    EPGO_DBNAME    this is the BoltDB filename.

    EPGO_TEMPLATE_PATH  this is the directory that contains the templates
                   used used to generate the static content of the website.

    EPGO_HTDOCS    this is the directory where the HTML files are written.

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
	showHelp    bool
	showVersion bool
	showLicense bool

	htdocs         string
	dbName         string
	bleveName      string
	templatePath   string
	apiURL         string
	siteURL        string
	repositoryPath string
)

func usage(appName, version string) {
	fmt.Printf(description, appName, appName)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Printf(configuration, appName)
	fmt.Printf("%s %s\n", appName, version)
	os.Exit(0)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")

	flag.StringVar(&htdocs, "htdocs", "", "specify where to write the HTML files to")
	flag.StringVar(&dbName, "dbname", "", "the BoltDB name")
	flag.StringVar(&bleveName, "bleve", "", "the Bleve index/db name")
	flag.StringVar(&apiURL, "api-url", "", "the URL of the EPrints API")
	flag.StringVar(&siteURL, "site-url", "", "the website url")
	flag.StringVar(&templatePath, "templates", "", "specify where to read the templates from")
	flag.StringVar(&repositoryPath, "repository-path", "", "specify the repository path to use for generated content")
}

func main() {
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
		os.Exit(0)
	}

	var cfg epgo.Config

	// Check to see we can merge the required fields are merged.
	check(cfg.MergeEnv("EPGO", "HTDOCS", htdocs))
	check(cfg.MergeEnv("EPGO", "DBNAME", dbName))
	check(cfg.MergeEnv("EPGO", "TEMPLATE_PATH", templatePath))
	check(cfg.MergeEnv("EPGO", "SITE_URL", siteURL))
	// Merge any optional data
	cfg.MergeEnv("EPGO", "BLEVE", bleveName)
	cfg.MergeEnv("EPGO", "API_URL", apiURL)
	cfg.MergeEnv("EPGO", "REPOSITORY_PATH", repositoryPath)

	if cfg.Htdocs != "" {
		if _, err := os.Stat(htdocs); os.IsNotExist(err) {
			os.MkdirAll(htdocs, 0775)
		}
	}

	// Create an API instance
	api, err := epgo.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//
	// Read the boltdb indicated in configuration and
	// render pages in the various formats supported.
	//
	log.Printf("%s %s\n", appName, epgo.Version)
	log.Printf("Rendering pages from %s\n", cfg.DBName)
	err = api.BuildSite(-1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rendering complete")
}

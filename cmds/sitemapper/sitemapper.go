//
// sitemapper generates a sitemap.xml file by crawling the content generate with genpages
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
	"path/filepath"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/epgo"
)

type locInfo struct {
	Loc     string
	LastMod string
}

var (
	description = `
USAGE: %s [OPTIONS] HTDOCS_PATH MAP_FILENAME PUBLIC_BASE_URL

OVERVIEW

Generates a sitemap for the accession pages.
`

	configuration = `
`

	examples = `
EXAMPLE

    %s htdocs htdocs/sitemap.xml http://eprints.example.edu

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
	// Standard cli options
	showHelp    bool
	showVersion bool
	showLicense bool

	// App specific options
	apiURL       string
	dbName       string
	bleveName    string
	htdocs       string
	templatePath string
	siteURL      string

	changefreq string
	locList    []*locInfo
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func usage(description, configuration, examples, appName, version string) {
	fmt.Printf(description, appName)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Printf(configuration, appName)
	fmt.Printf(examples, appName)
	fmt.Printf("\n%s %s\n", appName, version)
}

func init() {
	// standard cli options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")

	// app specific cli options
	flag.StringVar(&changefreq, "u", "daily", "Set the change frequencely value, e.g. daily, weekly, monthly")
	flag.StringVar(&changefreq, "update-frequency", "daily", "Set the change frequencely value, e.g. daily, weekly, monthly")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()
	if showHelp == true {
		usage(description, configuration, examples, appName, epgo.Version)
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Printf("%s %s\n", appName, epgo.Version)
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Printf(license, appName, epgo.Version)
		os.Exit(0)
	}

	if len(args) != 3 {
		fmt.Printf("%s requires 3 parameters, see %s --help\n", appName, appName)
		os.Exit(1)
	}

	var cfg epgo.Config

	// Required
	check(cfg.MergeEnv("EPGO", "HTDOCS", htdocs))
	check(cfg.MergeEnv("EPGO", "SITE_URL", siteURL))

	// Optional
	cfg.MergeEnv("EPGO", "API_URL", apiURL)
	cfg.MergeEnv("EPGO", "DBNAME", dbName)
	cfg.MergeEnv("EPGO", "BLEVE", bleveName)
	cfg.MergeEnv("EPGO", "TEMPLATE_PATH", templatePath)

	if changefreq == "" {
		changefreq = "daily"
	}

	log.Printf("Starting map of %s\n", args[0])
	filepath.Walk(args[0], func(p string, info os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".html") {
			fname := path.Base(p)
			//NOTE: You can skip the eror pages in the sitemap
			if strings.HasPrefix(fname, "50") == false && strings.HasPrefix(p, "40") == false {
				finfo := new(locInfo)
				finfo.Loc = fmt.Sprintf("%s%s", args[2], strings.TrimPrefix(p, args[0]))
				yr, mn, dy := info.ModTime().Date()
				finfo.LastMod = fmt.Sprintf("%d-%0.2d-%0.2d", yr, mn, dy)
				log.Printf("Adding %s\n", finfo.Loc)
				locList = append(locList, finfo)
			}
		}
		return nil
	})
	fmt.Printf("Writing %s\n", args[1])
	fp, err := os.OpenFile(args[1], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatalf("Can't create %s, %s\n", args[1], err)
	}
	defer fp.Close()
	fp.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`))
	for _, item := range locList {
		fp.WriteString(fmt.Sprintf(`
    <url>
            <loc>%s</loc>
            <lastmod>%s</lastmod>
            <changefreq>%s</changefreq>
    </url>
`, item.Loc, item.LastMod, changefreq))
	}
	fp.Write([]byte(`
</urlset>
`))
}

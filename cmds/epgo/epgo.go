package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	// my package
	"github.com/rsdoiel/epgo"
)

var (
	// CLI options
	showHelp        bool
	showVersion     bool
	restAPI         bool
	prettyPrint     bool
	exportToDB      bool
	buildSite       bool
	publishedOldest int
	publishedNewest int

	// Configuration variables
	baseURL      string
	dbName       string
	htdocs       string
	templatesDir string
)

func init() {
	baseURL = os.Getenv("EPGO_BASE_URL")
	dbName = os.Getenv("EPGO_DBNAME")
	htdocs = os.Getenv("EPGO_HTDOCS")
	templatesDir = os.Getenv("EPGO_TEMPLATES")
	publishedNewest = 0
	publishedOldest = 0

	flag.BoolVar(&showHelp, "h", false, "display help info")
	flag.BoolVar(&showVersion, "v", false, "display version info")
	flag.BoolVar(&prettyPrint, "p", false, "pretty print JSON output")
	flag.BoolVar(&restAPI, "api", false, "read the contents from the API without saving in the database")
	flag.BoolVar(&exportToDB, "export", false, "export EPrints to database")
	flag.BoolVar(&buildSite, "build", false, "build pages and feeds from database")
	flag.IntVar(&publishedOldest, "published-oldest", 0, "list the N oldest published items")
	flag.IntVar(&publishedNewest, "published-newest", 0, "list the N newest published items")
}

func main() {
	flag.Parse()
	if showHelp == true {
		fmt.Println(`
 USAGE: epgo [OPTIONS] [EPRINT_URI]

 epgo wraps the REST API for E-Prints 3.3 or better. It can return a list of uri,
 a JSON view of the XML presentation as well as generates feeds and web pages.

 epgo can be configured with following environment variables

 + EPGO_BASE_URL (required) the URL to your E-Prints installation
 + EPGO_DBNAME   (required) the BoltDB name for exporting, site building, and content retrieval
 + EPGO_HTDOCS   (optional) the htdocs root for site building
 + EPGO_TEMPLATES (optional) the template directory to use for site building

 If EPRINT_URI is provided then an individual EPrint is return as
 a JSON structure (e.g. /rest/eprint/34.xml). Otherwise a list of EPrint paths are
 returned.

 OPTIONS

`)
		flag.PrintDefaults()

		fmt.Printf(`
 Version %s
`, epgo.Version)

		os.Exit(0)
	}

	if showVersion == true {
		fmt.Printf("Version %s\n", epgo.Version)
		os.Exit(0)
	}

	// This will read in any settings from the environment
	api, err := epgo.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	args := flag.Args()
	if exportToDB == true {
		if err := api.ExportEPrints(); err != nil {
			log.Fatalf("%s", err)
		}
		if buildSite == false {
			os.Exit(0)
		}
	}

	if buildSite == true {
		if err := api.BuildSite("recently-published"); err != nil {
			log.Fatalf("%s", err)
		}
		os.Exit(0)
	}

	//
	// Generate JSON output
	//
	var (
		src  []byte
		data interface{}
	)
	switch {
	case publishedNewest > 0:
		data, err = api.GetPublishedRecords(0, publishedNewest, epgo.Descending)
	case publishedOldest > 0:
		data, err = api.GetPublishedRecords(0, publishedOldest, epgo.Ascending)
	case restAPI == true:
		if len(args) == 1 {
			data, err = api.GetEPrint(args[0])
		} else {
			data, err = api.ListEPrintsURI()
		}
	default:
		if len(args) == 1 {
			data, err = api.Get(args[0])
		} else {
			data, err = api.ListURI(0, 1000000)
		}
	}

	if err != nil {
		log.Fatalf("%s", err)
	}

	if prettyPrint == true {
		src, _ = json.MarshalIndent(data, "", "    ")
	} else {
		src, _ = json.Marshal(data)
	}
	fmt.Printf("%s", src)
}

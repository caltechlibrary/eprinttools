package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rsdoiel/epgo"
	"log"
	"os"
)

var (
	showHelp    bool
	showVersion bool
	prettyPrint bool
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "display help info")
	flag.BoolVar(&showVersion, "v", false, "display version info")
	flag.BoolVar(&prettyPrint, "p", false, "pretty print JSON output")
}

func main() {
	flag.Parse()
	if showHelp == true {
		fmt.Println(`
 USAGE: epgo [OPTIONS] [EPRINT_ID_PATH]

 epgo provides a basic REST API interactions with E-Prints 3.3 or better.
 It uses the following environment variables if available.

 + EPGO_BASE_URL (required) the URL to your E-Prints installation
 + EPGO_USERNAME (optional) Username if you have the E-Prints REST API
                 retricted to authorized users
 + EPGO_PASSWORD (optional) Password if you have the E-Prints REST API
                 retricted to authorized users

 If EPRINT_ID_PATH is provided then an individual EPrint is return as
 a JSON structure. Otherwise a list of EPrint paths are returned.

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

	api, err := epgo.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	var (
		src  []byte
		data interface{}
	)

	args := flag.Args()
	if len(args) == 1 {
		data, err = api.GetEPrint(args[0])
	} else {
		data, err = api.ListEPrintsURI()
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

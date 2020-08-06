// apiexplorer is a demo program showing how to use and restclient package
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/rc"
)

var (
	synopsis = `
_apiexplorer_ is a demo of rc package for accessing REST API
`

	description = `
_apiexploere_ is a demo program exercising the rc Golang package.
`

	examples = `
_apiexplorer_ accessing a public DataCite API for DOI
"10.22002/D1.924".

` + "```" + `
    apiexplorer -method GET \
    "https://api.datacite.org/works/10.22002/D1.924?email=jdoe@example.edu"
` + "```" + `

	`

	// Standard Options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	outputFName      string
	generateMarkdown bool
	generateManPage  bool

	// Application Options
	authMethod string
	method     string
	userName   string
	userSecret string
	payload    string
	asJSON     bool
)

func main() {
	appName := path.Base(os.Args[0])

	// Configuration and command line interation
	app := cli.NewCli(rc.Version)
	app.AddHelp("synopsis", []byte(synopsis))
	app.AddHelp("description", []byte(description))
	app.AddHelp("examples", []byte(examples))
	app.AddHelp("license", []byte(fmt.Sprintf(rc.LicenseText, appName, rc.Version)))

	// Standard Options
	app.BoolVar(&showHelp, "h, help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.StringVar(&outputFName, "o,output", "", "output filename")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// Application Options
	app.StringVar(&authMethod, "auth", "", "set authorization type (e.g. oauth, shib)")
	app.StringVar(&userName, "un,username", "", "set username for authentication")
	app.StringVar(&userSecret, "pw,password", "", "set user secret to use for authentication")
	app.StringVar(&method, "method", "GET", "set the http method to use for request")
	app.StringVar(&payload, "payload", "", "A JSON structure holding the payload data")
	app.BoolVar(&asJSON, "as-json", false, "Convert XML to JSON before output")

	app.Parse()
	args := app.Args()

	// Pull environment if anything is unset
	if userName == "" {
		userName = os.Getenv("USERNAME")
	}
	if userSecret == "" {
		userSecret = os.Getenv("PASSWORD")
	}
	if authMethod == "" {
		authMethod = os.Getenv("AUTH_METHOD")
	}

	// Process options
	if generateMarkdown {
		app.GenerateMarkdown(os.Stdout)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(os.Stdout)
		os.Exit(0)
	}
	if showHelp {
		if len(args) > 0 {
			fmt.Fprintf(os.Stdout, app.Help(args...))
		} else {
			app.Usage(os.Stdout)
		}
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintln(os.Stdout, app.License())
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintln(os.Stdout, app.Version())
		os.Exit(0)
	}

	out, err := cli.Create(outputFName, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer cli.CloseFile(outputFName, out)

	authType := rc.AuthNone
	switch strings.TrimSpace(strings.ToLower(authMethod)) {
	case "basic":
		authType = rc.BasicAuth
	case "oath":
		authType = rc.OAuth
	case "shib":
		authType = rc.Shibboleth
	}

	data := map[string]string{}
	if payload != "" {
		if err := json.Unmarshal([]byte(payload), &data); err != nil {
			log.Fatal(err)
		}
	}

	for _, arg := range args {
		u, err := url.Parse(arg)
		if err != nil {
			log.Fatal(err)
		}
		api, err := rc.New(arg, authType, userName, userSecret)
		if err != nil {
			log.Fatal(err)
		}
		if src, err := api.Request(method, u.Path, data); err == nil {
			if asJSON == true {
				if bytes.HasPrefix(src, []byte("<")) {
					m := map[string]interface{}{}
					if err := xml.Unmarshal(src, &m); err != nil {
						log.Fatal(err)
					} else {
						s, err := json.Marshal(m)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Sprintf("%s\n", s)
					}
				} else {
					fmt.Fprintf(out, "%s\n", src)
					os.Exit(0)
				}
			}
			fmt.Fprintf(out, "%s\n", src)
			os.Exit(0)
		} else {
			log.Fatal(err)
		}
	}
}

/**
 * ep3api.go defines an extended EPrints web server.
 */
package eprinttools

import (
	"fmt"
	"log"
	"net/http"
)

var (
	AppName string
	Cfg     *Config
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
EPrints 3.x Extended API
%s %s
`, AppName, Version)
}

func InitExtendedAPI(appName string, cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("Missing configuration")
	}
	if appName == "" {
		return fmt.Errorf("Missing application name")
	}
	Cfg = cfg
	AppName = appName
	return nil
}

func RunExtendedAPI() error {
	/* Setup web server */
	http.HandleFunc("/", indexPage)

	fmt.Printf(`
EPrints 3.x Extended API
%s %s

Listening on %s

Press ctl-c to terminate.
`, AppName, Version, Cfg.Hostname)
	return http.ListenAndServe(Cfg.Hostname, nil)
}

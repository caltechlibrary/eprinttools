package rc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Version = `v0.0.7`

	LicenseText = `
%s %s

Copyright (c) 2017, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
)

//FIXME: Need to handle Basic auth, OAuth and Shibboleth
const (
	AuthNone = iota
	BasicAuth
	OAuth
	Shibboleth
)

type RestAPI struct {
	u        *url.URL
	id       string
	secret   string
	authType int
	token    string
	headers  map[string]string

	// Timeout is the client time out period, default is 10 seconds
	Timeout time.Duration
}

// New creates a new Rest Client RestAPI instance
// If clientID and clientSecret and empty and authType
func New(apiURL string, authType int, clientID, clientSecret string) (*RestAPI, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	if len(clientID) == 0 && len(clientSecret) == 0 && u.User != nil {
		clientID = u.User.Username()
		if pword, ok := u.User.Password(); ok == true {
			clientSecret = pword
		}
	}
	return &RestAPI{
		u:        u,
		id:       clientID,
		secret:   clientSecret,
		authType: authType,
		token:    "",
		Timeout:  10 * time.Second,
	}, nil
}

// AddHeader sets the header strings to send with the request
func (api *RestAPI) AddHeader(ky, value string) {
	if api.headers == nil {
		api.headers = map[string]string{}
	}
	api.headers[ky] = value
}

func (api *RestAPI) Login() error {
	switch api.authType {
	case AuthNone:
		return nil
	case BasicAuth:
		return api.basicAuthLogin()
	case OAuth:
		return api.oAuthLogin()
	case Shibboleth:
		return api.shibbolethLogin()
	default:
		return fmt.Errorf("auth type not supported")
	}
}

func (api *RestAPI) oAuthLogin() error {
	if api.token == "" {
		client := &http.Client{
			Timeout: api.Timeout,
		}
		u := api.u
		u.Path = "/oauth/token"

		// OAuth2 authentication is usually done with a POST, need to setup the form values
		// and URL encode the results.
		payload := map[string]string{
			"client_id":     api.id,
			"client_secret": api.secret,
			"scope":         "/read-public",
			"grant_type":    "client_credentials",
		}
		form := u.Query()
		for key, value := range payload {
			form.Add(key, value)
		}

		// OK, we're ready to setup our request, send it and get on our way.
		req, err := http.NewRequest("POST", u.String(), strings.NewReader(form.Encode()))
		if err != nil {
			return err
		}
		// Need to set the mime type for the content we're sending to the RestAPI
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Add any additional custom headers
		for k, v := range api.headers {
			req.Header.Add(k, v)
		}

		// Get the text response for RestAPI
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		src, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// We have a throw away response object from authenticating in
		data := &struct {
			AccessToken  string `json:"access_token"`
			Bearer       string `json:"bearer"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int    `json:"expires_in"`
			Scope        string `json:"scope"`
		}{}
		err = json.Unmarshal(src, &data)
		if err != nil {
			return err
		}
		if data.AccessToken == "" {
			return fmt.Errorf("Could not decode access token from %q, %+v", src, data)
		}
		api.token = data.AccessToken
	}
	return nil
}

// basicAuthLogin() implement basic HTTP Authorization
func (api *RestAPI) basicAuthLogin() error {
	if api.id == "" || api.secret == "" {
		return fmt.Errorf("Missing username or password for Basic Auth request")
	}
	return nil
}

// shibbolethLogin() implement Shibboleth HTTP Authorization
func (api *RestAPI) shibbolethLogin() error {
	return fmt.Errorf("shibbolethLogin() not implemented")
}

// Request contacts the Rest API and returns the full read response body, and error
// payload is the used to build the URL Query object (e.g. ?key=value&key1=value...)
func (api *RestAPI) Request(method, docPath string, payload map[string]string) ([]byte, error) {
	var (
		req *http.Request
		err error
	)
	// Create a http client
	client := &http.Client{
		Timeout: api.Timeout,
	}

	// NOT: if api.token not set we should just go ahead and oAuthLogin.
	if api.token == "" {
		if err := api.Login(); err != nil {
			return nil, err
		}
	}

	// NOTE: we want a copy the URL in Rest API object and update copy with the docPath
	u := api.u
	u.Path = docPath

	// NOTE: Based the HTTP method we want, we build our request appropriately
	switch strings.ToUpper(method) {
	case "GET":
		req, err = http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, err
		}
		// NOTE: If we're using Basic Auth setup the request with it
		if api.authType == BasicAuth {
			req.SetBasicAuth(api.id, api.secret)
		}
		// NOTE: If we've authenticated we need to path the auth token
		if len(api.token) > 0 {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.token))
		}
		// NOTE: We need to indicate the format we want
		req.Header.Add("Accept", "application/json")

		// NOTE: Build our payload to pass in the URL since this is a GET
		qry := req.URL.Query()
		for key, value := range payload {
			qry.Add(key, value)
		}
		req.URL.RawQuery = qry.Encode()
	default:
		return nil, fmt.Errorf("Do not know how to make a %s request", method)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		src, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return src, nil
	}
	return nil, fmt.Errorf("%s for %s", resp.Status, u.String())
}

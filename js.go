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
package epgo

//
// js.go adds support for easily creating JavaScript REPL services
//
import (
	"fmt"
	"log"
	"strconv"

	// Caltech library packages
	"github.com/caltechlibrary/ostdlib"

	// 3rd Party packages
	"github.com/robertkrimen/otto"
)

// AddExtensions creates a *otto.Otto (JavaScript VM) with functions added to integrate
// epgo REPL.
func (api *EPrintsAPI) AddExtensions(js *ostdlib.JavaScriptVM) *otto.Otto {
	vm := js.VM
	errorObject := func(obj *otto.Object, msg string) otto.Value {
		if obj == nil {
			obj, _ = vm.Object(`({})`)
		}
		log.Println(msg)
		obj.Set("status", "error")
		obj.Set("error", msg)
		return obj.Value()
	}

	// responseObject := func(data interface{}) otto.Value {
	// 	src, _ := json.Marshal(data)
	// 	obj, _ := vm.Object(fmt.Sprintf(`(%s)`, src))
	// 	return obj.Value()
	// }

	// EPrints REST API methods
	apiObj, _ := vm.Object(`api = {}`)
	apiObj.Set("listEPrintsURI", func(call otto.FunctionCall) otto.Value {
		uris, err := api.ListEPrintsURI()
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listEPrintsURI() failed %s, %s", call.CallerLocation(), err))
		}
		result, err := vm.ToValue(uris)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listEPrintsURI() failed %s, %s", call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("getEPrint", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("Missing uri for getEPrint() %s", call.CallerLocation()))
		}
		uri := call.Argument(0).String()
		record, err := api.GetEPrint(uri)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getEPrint(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(record)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getEPrint(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("exportEPrints", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("Missing export count for exportEPrints(n) %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		cnt, _ := strconv.Atoi(s)
		err := api.ExportEPrints(cnt)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("exportEprints(%d) failed %s, %s", cnt, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(true)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("exportEprints(%d) failed %s, %s", cnt, call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("listURI", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			return errorObject(nil, fmt.Sprintf("listURI(start, count) missing parameters %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		start, _ := strconv.Atoi(s)
		s = call.Argument(1).String()
		count, _ := strconv.Atoi(s)
		uris, err := api.ListURI(start, count)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listURI(%d, %d) failed %s, %s", start, count, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(uris)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("listURI(%d, %d) failed %s, %s", start, count, call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("get", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("get(uri) missing parameters %s", call.CallerLocation()))
		}
		uri := call.Argument(0).String()
		record, err := api.Get(uri)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("get(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(record)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("get(%q) failed %s, %s", uri, call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("getPublishedRecords", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			return errorObject(nil, fmt.Sprintf("getPublishedRecords(start, count, direction) missing parameters %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		start, _ := strconv.Atoi(s)
		s = call.Argument(1).String()
		count, _ := strconv.Atoi(s)
		s = call.Argument(2).String()
		direction, _ := strconv.Atoi(s)
		records, err := api.GetPublishedRecords(start, count, direction)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedRecords(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(records)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedRecords(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("getPublishedArticles", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			return errorObject(nil, fmt.Sprintf("getPublishedArticles(start, count, direction) missing parameters %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		start, _ := strconv.Atoi(s)
		s = call.Argument(1).String()
		count, _ := strconv.Atoi(s)
		s = call.Argument(2).String()
		direction, _ := strconv.Atoi(s)
		records, err := api.GetPublishedArticles(start, count, direction)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedArticles(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(records)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("getPublishedArticles(%d, %d, %d) failed %s, %s", start, count, direction, call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("renderDocuments", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 4 {
			return errorObject(nil, fmt.Sprintf("renderDocuments(docTitle, docDescription, basepath, records) missing parameters %s", call.CallerLocation()))
		}
		docTitle := call.Argument(0).String()
		docDescription := call.Argument(1).String()
		basepath := call.Argument(2).String()
		records := []*Record{}
		err := ostdlib.ToStruct(call.Argument(3), records)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("renderDocuments(%q, %q, %q, records) records error %s, %s", docTitle, docDescription, basepath, call.CallerLocation(), err))
		}
		if err := api.RenderDocuments(docTitle, docDescription, basepath, records); err != nil {
			return errorObject(nil, fmt.Sprintf("renderDocuments(%q, %q, %q, records) failed %s, %s", docTitle, docDescription, basepath, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(true)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("renderDocuments(%q, %q, %q, records) failed %s, %s", docTitle, docDescription, basepath, call.CallerLocation(), err))
		}
		return result
	})

	apiObj.Set("buildSite", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			return errorObject(nil, fmt.Sprintf("buildSite(feed_size) missing parameter %s", call.CallerLocation()))
		}
		s := call.Argument(0).String()
		feedSize, _ := strconv.Atoi(s)
		err := api.BuildSite(feedSize)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("buildSite(%d) failed %s, %s", feedSize, call.CallerLocation(), err))
		}
		result, err := vm.ToValue(true)
		if err != nil {
			return errorObject(nil, fmt.Sprintf("buildSite(%d) failed %s, %s", feedSize, call.CallerLocation(), err))
		}
		return result
	})
	return vm
}

// AddHelp builds the help structures for use in the REPL
func (api *EPrintsAPI) AddHelp(js *ostdlib.JavaScriptVM) {
	js.SetHelp("api", "listEPrintsURI", []string{}, "returns and array of EPrint uris")
	js.SetHelp("api", "getEPrint", []string{"uri string"}, "given an EPrint uri (e.g. /eprint/2026) return the eprint as a JavaScript object")
	js.SetHelp("api", "exportEPrints", []string{"N int"}, "Export N items from eprints. Exports highest value id first. If N == -1 then export everything")
	js.SetHelp("api", "listURI", []string{}, "return a list of URI saved in the boltdb")
	js.SetHelp("api", "get", []string{"uri string"}, "given a uri, return the results saved in the boltdb")
	js.SetHelp("api", "getPublishedRecords", []string{"start int", "count int", "direction int"}, "generate a list of published eprints starting at start for count by direction. Direction is 0 for ascending, 1 for descending")
	js.SetHelp("api", "getPublishedArticles", []string{"start int", "count int", "direction int"}, "generate a list of published articles eprints starting at start for count by direction. Direction is 0 for ascending, 1 for descending")
	js.SetHelp("api", "renderDocuments", []string{"docTitle string", "docDescription string", "basepath string", "records an array of eprints"},
		"render the eprint records list as an HTML file, HTML include, RSS and JSON documents")
	js.SetHelp("api", "buildSite", []string{"feedSize int"}, "build feeds recently-published and recent-articles. feedSize indicates the maximun number of items included")
}

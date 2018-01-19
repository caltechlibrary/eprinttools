//
// Package eprinttools is a collection of structures and functions for working with the E-Prints REST API
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
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
package eprinttools

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

// CustomFields defines a repository's field customization. This should prove helpful
// when deploying development instances and keeping them in sync with out know instances.
type CustomFields struct {
	XMLName        xml.Name                `json:"-"`
	Created        time.Time               `xml:"created" json:"created"`
	LastModified   time.Time               `xml:"last_modified" json:"last_modified"`
	RepositoryName string                  `xml:"repository_name" json:"repository_name"`
	RepositoryID   string                  `xml:"repository_id" json:"repostiory_id"`
	EPrintsPath    string                  `xml:"eprints_path,omitempty" json:"eprints_path,omitempty"`
	CustomFields   map[string]*CustomField `xml:"custom_fields,omitempty" json:"custom_fields,omitempty"`
}

// CustomField holds a repository's custom field definition.
type CustomField struct {
	XMLName                   xml.Name  `json:"-"`
	SystemName                string    `xml:"system_name" json:"system_name"`
	Public                    bool      `xml:"public" json:"public"`
	Type                      string    `xml:"type" json:"type"`
	IsCompound                bool      `xml:"is_compound" json:"is_compound"`
	IsRepeatable              bool      `xml:"is_repeatable" json:"is_repeatable"`
	IncludeInAbstract         bool      `xml:"include_in_abstract" json:"include_in_abstract"`
	IncludeInSubmissionForm   bool      `xml:"include_in_submission_form" json:"include_in_submission_form"`
	CollapsedInSubmissionForm bool      `xml:"collapsed_in_submission_form,omitempty" json:"collapsed_in_submission_form,omitempty"`
	IsRequired                bool      `xml:"is_required" json:"is_required"`
	DefaultValue              string    `xml:"default_value,omitempty" json:"default_value,omitempty"`
	ControlledVocabulary      []string  `xml:"controlled_vocabulary,omitempty" json:"controlled_vocabulary,omitempty"`
	ValidationRequired        bool      `xml:"validation_required,omitempty" json:"validation_required,omitempty"`
	AjaxLookupFilename        string    `xml:"ajax_lookup_filename,omitempty" json:"ajax_lookup_filename,omitempty"`
	Created                   time.Time `xml:"created" json:"created"`
	LastModified              time.Time `xml:"last_modified" json:"last_modified"`
}

// String render a custom fields struct as a string (using JSON notation)
func (cfields *CustomFields) String() string {
	src, _ := json.MarshalIndent(cfields, "", "  ")
	return fmt.Sprintf("%s", src)
}

// String render a custom field struct as a string (using JSON notation)
func (cf *CustomField) String() string {
	src, _ := json.MarshalIndent(cf, "", "  ")
	return fmt.Sprintf("%s", src)
}

// NewCustomFields creates a new, empty, custom fields stuct
func NewCustomFields() (*CustomFields, error) {
	return nil, fmt.Errorf("NewCustomFields() not implemented")
}

// NewCustomField creates a new, empty, custom field stuct
func NewCustomField() (*CustomField, error) {
	return nil, fmt.Errorf("NewCustomField() not implemented")
}

// Prompts returns a map array of prompt values by struct field name
func (cfields *CustomFields) Prompts() map[string]string {
	return map[string]string{
		"error": "Prompts() not implemented",
	}
}

// Update maps the "answers" into their field values for the custom fields struct
// It returns an error if there is a problem. Note it doesn't add/update individual
// fields, just the outer struct. For adding/updating a custom field itself use
// the *CustomField version of Add(), Update(), Depreciate()
func (cfields *CustomFields) Update(answers map[string]interface{}) error {
	return fmt.Errorf("Update() not implemented")
}

// Deploy process a CustomFields struct and updates an EPrints instance based on
// the state defined in the struct. IMPORTANT: as this changes files on disc you
// shouldn't use this in production (without regorous testing and a back up of
// disc and MySQL).
func (cfields *CustomFields) Deploy() error {
	return fmt.Errorf("Deploy() not implemented")
}

// Prompts maps a list of of custom field attributes names to individual prompts.
// It doesn't require an existing custom field system_name since all fields should
// be have the same definition structure.
func (cf *CustomField) Prompts() map[string]string {
	return map[string]string{
		"error": "Prompts() not implemented",
	}
}

// Update adds or updates a custom field based on answers map
func (cf *CustomField) Update(answers map[string]interface{}) error {
	return fmt.Errorf("Update() not implemented")
}

// Depreciate flags a field as depreciated and hides UI references to the field (as
// you can't really delelte fields from EPrints)
func (cf *CustomField) Depreciate(systemName string) error {
	return fmt.Errorf("Depreciate() not implemented")
}

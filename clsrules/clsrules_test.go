//
// clsrules_test.go (Caltech Library Specific Rules) is a package for
// implementing Caltech Library Specific features to processing
// and creating EPrint XML. Currently these include things like
// trimming prefixed "The " from titles, dropping series information,
// changing how the date is derived and very idiosencratic handling
// of Author and DOI references.
//
package clsrules

import (
	"testing"
)

func TestApply(t *testing.T) {
	//FIXME: generate unmodified DOI to EPrintXML list
	//FIXME: Apply clsrules to EPrintXML list
	//FIXME: Check if rules applied correctly.
	t.Errorf("clsrules.Apply() not implemented")
}

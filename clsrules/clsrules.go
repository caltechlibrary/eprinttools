//
// clsrules.go (Caltech Library Specific Rules) is a package for
// implementing Caltech Library Specific features to processing
// and creating EPrint XML. Currently these include things like
// trimming prefixed "The " from titles, dropping series information,
// changing how the date is derived and very idiosencratic handling
// of Author and DOI references.
//
package clsrules

import (
	"github.com/caltechlibrary/eprinttools"
)

func Apply(eprintsList *eprinttools.EPrints) (*eprinttools.EPrints, error) {
	return eprintsList, nil
}

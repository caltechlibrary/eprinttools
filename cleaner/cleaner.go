//
// Cleaner implements a naive HTML/XML sanitizer
//
package cleaner

import (
	"bytes"
	"regexp"
)

var (
	// RE Found on stack overflow: https://stackoverflow.com/questions/39329607/regex-find-all-xml-tags#47511815
	reOpenElement  = regexp.MustCompile(`<([^\/>]+)[/]*>`)
	reCloseElement = regexp.MustCompile(`</([^\/>]+)[/]*>`)
	reEntity       = regexp.MustCompile(`&\w+;`)

	entities = map[string]string{
		`&amp;`:  `&`,
		`&lt;`:   `<`,
		`&gt;`:   `>`,
		`&copy;`: `©`,
	}
)

// HasEncodedElements checks a byte array to determine if it has any &lt; and &gt; or &amp;
func HasEncodedElements(src []byte) bool {
	if reEntity.Match(src) {
		return true
	}
	return reOpenElement.Match(src) || reCloseElement.Match(src)
}

// StripTags removes HTML/XML tags identified through the regexp or entities
func StripTags(src []byte) []byte {
	if reEntity.Match(src) {
		for target, dest := range entities {
			src = bytes.ReplaceAll(src, []byte(target), []byte(dest))
		}
	}
	if reOpenElement.Match(src) {
		src = reOpenElement.ReplaceAll(src, []byte(``))
	}
	if reCloseElement.Match(src) {
		src = reCloseElement.ReplaceAll(src, []byte(``))
	}
	return src
}

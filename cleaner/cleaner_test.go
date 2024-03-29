package cleaner

import (
	"bytes"
	"testing"
)

func assertBool(t *testing.T, expected bool, got bool, msg string, args ...interface{}) {
	if expected != got {
		t.Errorf(msg, args...)
	}
}

func assertBytes(t *testing.T, expected []byte, got []byte, msg string, args ...interface{}) {
	if bytes.Compare(expected, got) != 0 {
		t.Errorf(msg, args...)
	}
}

func TestHasEncodedElements(t *testing.T) {
	var (
		src           []byte
		expected, got bool
	)
	src = []byte(`
This text has no encoded elements. cleaner.HasEncodedElements(src) should fail.

This text has no encoded elements. cleaner.HasEncodedElements(src) should fail.
`)
	expected, got = false, HasEncodedElements(src)
	assertBool(t, expected, got, `expected %t, got %t, cleaner.HasEncodedElements(src): %s`, expected, got, src)

	src = []byte(`&lt;jt&gt;This is a jats thing&lt;/jt&gt;`)
	expected, got = true, HasEncodedElements(src)
	assertBool(t, expected, got, `expected %t, got %t, cleaner.HasEncodedElements(src): %s`, expected, got, src)

	src = []byte(`
<i>This text has</i> encoded elements. cleaner.HasEncodedElements(src) should return true!
`)
	expected, got = true, HasEncodedElements(src)
	assertBool(t, expected, got, `expected %t, got %t, cleaner.HasEncodedElements(src): %s`, expected, got, src)

	src = []byte(`
This text has <span id="3" class="wicked-thing">encoded elements</span>. cleaner.HasEncodedElements(src) should return true!
`)
	expected, got = true, HasEncodedElements(src)
	assertBool(t, expected, got, `expected %t, got %t, cleaner.HasEncodedElements(src): %s`, expected, got, src)
}

func TestStripTags(t *testing.T) {
	var (
		src, expected, got []byte
	)
	src = []byte(`
This text has no encoded elements. cleaner.StripTags(src) should not change.

This text has no encoded elements. cleaner.StripTags(src) should not change.
`)
	expected = src[:]
	got = StripTags(src)
	assertBytes(t, expected, got, `expected %q, got %q, cleaner.StripTags(src): %s`, expected, got, src)

	src = []byte(`&lt;jt&gt;This is a jats thing&lt;/jt&gt;`)
	expected = src[:]
	for _, target := range []string{`&lt;jt&gt;`, `&lt;/jt&gt;`} {
		expected = bytes.ReplaceAll(expected, []byte(target), []byte(``))
	}
	got = StripTags(src)
	assertBytes(t, expected, got, `expected %q, got %q, cleaner.StripTags(src): %s`, expected, got, src)

	src = []byte(`
<i>This text has</i> encoded elements. cleaner.StripTags(src) should return a clean byte array!
`)
	expected = src[:]
	for _, target := range []string{`<i>`, `</i>`} {
		expected = bytes.ReplaceAll(expected, []byte(target), []byte(``))
	}
	got = StripTags(src)
	assertBytes(t, expected, got, `expected %q, got %q, cleaner.StripTags(src): %s`, expected, got, src)

	src = []byte(`
This text has <span id="3" class="wicked-thing">encoded elements</span>. cleaner.StripTags(src) should return a clean byte array!!
`)
	expected = src[:]
	for _, target := range []string{`<span id="3" class="wicked-thing">`, `</span>`} {
		expected = bytes.ReplaceAll(expected, []byte(target), []byte(``))
	}
	got = StripTags(src)
	assertBytes(t, expected, got, `expected %q, got %q, cleaner.StripTags(src): %s`, expected, got, src)

	src = []byte(`<jats:title>Abstract</jats:title><jats:p>C1q/TNF-related protein 1 (CTRP1) is a CTRP family member that has collagenous and globular C1q-like domains. The secreted form of CTRP1 is known to be associated with cardiovascular and metabolic diseases, but its cellular roles have not yet been elucidated. Here, we showed that cytosolic CTRP1 localizes to the endoplasmic reticulum (ER) membrane and that knockout or depletion of CTRP1 leads to mitochondrial fission defects, as demonstrated by mitochondrial elongation. Mitochondrial fission events are known to occur through an interaction between mitochondria and the ER, but we do not know whether the ER and/or its associated proteins participate directly in the entire mitochondrial fission event. Interestingly, we herein showed that ablation of CTRP1 suppresses the recruitment of DRP1 to mitochondria and provided evidence suggesting that the ER–mitochondrion interaction is required for the proper regulation of mitochondrial morphology. We further report that CTRP1 inactivation-induced mitochondrial fission defects induce apoptotic resistance and neuronal degeneration, which are also associated with ablation of DRP1. These results demonstrate for the first time that cytosolic CTRP1 is an ER transmembrane protein that acts as a key regulator of mitochondrial fission, providing new insight into the etiology of metabolic and neurodegenerative disorders.</jats:p>`)

	expected = src[:]
	for _, target := range []string{`<jats:title>`, `</jats:title>`, `<jats:p>`, `</jats:p>`} {
		expected = bytes.ReplaceAll(expected, []byte(target), []byte(``))
	}
	got = StripTags(src)
	assertBytes(t, expected, got, `expected %q, got %q, cleaner.StripTags(src): %s`, expected, got, src)
}

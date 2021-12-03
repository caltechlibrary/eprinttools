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
}

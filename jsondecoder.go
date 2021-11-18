package eprinttools

import (
	"bytes"
	"encoding/json"
	"io"
)

// Custom JSON decoder so we can treat numbers easier
func jsonDecode(src []byte, obj interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	dec.UseNumber()
	err := dec.Decode(&obj)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

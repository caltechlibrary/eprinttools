package eprinttools

import (
	"encoding/json"
	"bytes"
	"os"
	"io"
)

// Custom JSON encoder so we can treat numbers easier
func jsonEncode(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "    ")
}

// Custom JSON encoder with write to file
func jsonEncodeToFile(fName string, obj interface{}, perms os.FileMode) error {
	src, err := jsonEncode(obj)
	if err != nil {
		return err
	}
	if err := os.WriteFile(fName, src, perms); err != nil {
		return err
	}
	return nil
}

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

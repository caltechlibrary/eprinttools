package eprinttools

import (
	"encoding/json"
)

// Custom JSON decoder so we can treat numbers easier
func jsonEncode(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "    ")
}

package eprinttools

// indexInto takes a map and walks the path and returns a string
// value and succeess bool, if path not found string value is
// an empty string and bool is false.
func indexInto(m map[string]interface{}, parts ...string) (interface{}, bool) {
	switch len(parts) {
	case 0:
		return "", false
	case 1:
		if val, ok := m[parts[0]]; ok == true {
			return val, true
		}
		return "", false
	default:
		if val, ok := m[parts[0]]; ok == true {
			switch val.(type) {
			case map[string]interface{}:
				return indexInto(val.(map[string]interface{}), parts[1:]...)
			}
		}
		return "", false
	}
}

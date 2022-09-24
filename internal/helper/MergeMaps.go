package helper

func MergeMap(a, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		// Test if the key exists in the first map
		if _, ok := a[k]; ok {
			// Test if the value is a map
			if _, ok = v.(map[string]interface{}); ok {
				// Test if the value of the first map is a map
				if _, ok = a[k].(map[string]interface{}); ok {
					// Because both values are maps, merge them
					a[k] = MergeMap(a[k].(map[string]interface{}), v.(map[string]interface{}))
					continue
				}
			}
		}
		// If the key doesn't exist in the first map, or the value is not a map, just set the value
		if v == nil {
			delete(a, k)
		} else {
			a[k] = v
		}
	}
	return a
}

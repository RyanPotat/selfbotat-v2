package bot

type Predicate func(any) bool

func Filter(slice []any, predicate Predicate) []any {
	var result []any
	for _, v := range slice {
			if predicate(v) {
					result = append(result, v)
			}
	}
	return result
}

func Find(slice []interface{}, predicate Predicate) interface{} {
	for _, v := range slice {
		if predicate(v) {
			return v
		}
	}
	return nil
}
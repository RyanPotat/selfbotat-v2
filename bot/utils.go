package bot

type Predicate[T any] func(T) bool

func Filter[T any](slice []T, predicate Predicate[T]) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func Find[T any](slice []T, predicate Predicate[T]) (T, bool) {
	for _, v := range slice {
		if predicate(v) {
			return v, true
		}
	}
	var notFound T
	return notFound, false
}
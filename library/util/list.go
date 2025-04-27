package util

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	ret = make([]T, 0, len(ss))

	for _, s := range ss {
		if !test(s) {
			continue
		}

		ret = append(ret, s)
	}

	return
}

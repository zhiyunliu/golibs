package v2

type SortPatterns[T Pattern] []T

func (s SortPatterns[T]) Len() int { return len(s) }

func (s SortPatterns[T]) Swap(i, j int) {
	if i >= len(s) || j >= len(s) {
		return
	}
	s[i], s[j] = s[j], s[i]
}

func (s SortPatterns[T]) Less(i, j int) bool {
	iv := s[i].Pattern()
	jv := s[j].Pattern()
	il := len(iv)
	jl := len(jv)

	for x := 0; x < jl && x < il; x++ {
		if iv[x] == jv[x] {
			continue
		}
		if iv[x] == byte('*') {
			return false
		}
		if jv[x] == byte('*') {
			return true
		}
		return iv[x] < jv[x]
	}
	return iv < jv
}

package collections

func Values[K comparable, V any](in map[K]V) []V {
	out := make([]V, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}

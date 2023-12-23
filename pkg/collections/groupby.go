package collections

func GroupBy[K comparable, V any, R any](in []V, f func(V) (K, R)) map[K][]R {
	out := map[K][]R{}
	for _, v := range in {
		k, r := f(v)
		out[k] = append(out[k], r)
	}
	return out
}

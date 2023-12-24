package collections

func Map[V any, R any](in []V, f func(V) (R, error)) ([]R, error) {
	out := make([]R, 0, len(in))
	for _, v := range in {
		r, err := f(v)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

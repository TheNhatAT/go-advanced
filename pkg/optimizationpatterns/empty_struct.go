package optimizationpatterns

//TODO: pprof these func to see why the allocs in HasDuplicates_Better is 19 while HasDuplicates is 3

func HasDuplicates[T comparable](slice ...T) bool {
	dup := make(map[T]any, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = "whatever, I don't use this value"
	}
	return false
}

func HasDuplicates_Better[T comparable](slice ...T) bool {
	dup := make(map[T]struct{}, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = struct{}{}
	}
	return false
}

func HasDuplicates_NonGeneric(slice ...float64) bool {
	dup := make(map[float64]struct{}, len(slice))
	for _, s := range slice {
		if _, ok := dup[s]; ok {
			return true
		}
		dup[s] = struct{}{}
	}
	return false
}

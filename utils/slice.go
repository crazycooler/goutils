package utils

func ArrayIndex[T comparable](needle T, hystack []T) int {
	for i, item := range hystack {
		if needle == item {
			return i
		}
	}

	return -1
}

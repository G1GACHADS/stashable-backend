package main

import "golang.org/x/exp/constraints"

func MaxSlice[T constraints.Ordered](slice []T) T {
	max := T(0)
	for _, v := range slice {
		if v > max {
			max = v
		}
	}
	return max
}

func MinSlice[T constraints.Ordered](slice []T) T {
	min := T(0)
	for _, v := range slice {
		if v < min {
			min = v
		}
	}
	return min
}

func RemoveDuplicates[T constraints.Ordered](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

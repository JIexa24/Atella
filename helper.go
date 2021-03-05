package atella

import "time"

// Pause implement interruptable pause.
func Pause(interval int64, interrupt *bool) {
	c := int64(0)
	one := time.Duration(1) * time.Second
	for {
		if *interrupt {
			break
		}
		if c >= interval {
			break
		}
		c = c + 1
		time.Sleep(one)
	}
}

// ElExistsString check string array and return true if item exist.
func ElExistsString(array []string, item string) bool {
	for i := 0; i < len(array); i = i + 1 {
		if array[i] == item {
			return true
		}
	}
	return false
}

// ElExistsInt64 check string array and return true if item exist.
func ElExistsInt64(array []int64, item int64) bool {
	for i := 0; i < len(array); i = i + 1 {
		if array[i] == item {
			return true
		}
	}
	return false
}

// SubsetInt64 return whether first is a subset of second
func SubsetInt64(first, second []int64) bool {

	set := make(map[int64]bool)
	for _, value := range second {
		set[value] = true
	}

	for _, value := range first {
		if !set[value] {
			return false
		}
	}

	return true
}

package utime

import "time"

// DateBefore t1是否小于t2 只比较到天
func DateBefore(t1, t2 time.Time) bool {
	if t1.Year() != t2.Year() {
		return t1.Year() < t2.Year()
	}
	if t1.Month() != t2.Month() {
		return t1.Month() < t2.Month()
	}
	return t1.Day() < t2.Day()
}

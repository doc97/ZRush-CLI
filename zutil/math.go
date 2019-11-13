/*
Package zutil contains utility methods used in ZRush.
*/
package zutil

// Min returns the smaller of the two arguments.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

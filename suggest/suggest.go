package suggest

// Completion is the signature of a function supporting simple 1..n suggestions
type Completion func(value string) []string

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}

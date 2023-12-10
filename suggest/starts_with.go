package suggest

import "strings"

type swOpts struct {
	minLen int
}

// StartsWithOpt represents a function type that manipulates an internal options configuration.
type StartsWithOpt func(o *swOpts)

// StartsWithMin returns a StartsWithOpt function that modifies the minimum string length of the internal options configuration.
func StartsWithMin(minimum int) StartsWithOpt {
	return func(o *swOpts) {
		if minimum >= 0 {
			o.minLen = minimum
		}
	}
}

// StartsWith takes a slice of strings (data) and returns a function (Completion).
// The Completion function takes a string (search query) and returns all strings from data that start with the search query.
// The search is configured by StartsWithOpt-functional options. By default, the minimum length of the search query is 3.
func StartsWith(data []string, options ...StartsWithOpt) Completion {
	opts := swOpts{
		minLen: 3,
	}
	for _, opt := range options {
		opt(&opts)
	}

	input := data[:]
	// update/view are done via goroutines, so we need to synchronize shared data between threads
	persistent := safeResults{results: make([]string, 0)}
	return func(value string) []string {
		persistent.mux.Lock()
		defer persistent.mux.Unlock()

		length := len(value)
		if length == 0 || length < opts.minLen {
			persistent.current = value
			return persistent.results
		}

		// when length equals minLength, we must fall through to the initialization branch.
		// otherwise, we will reuse the initialized search results to filter the subset of data.
		if length > opts.minLen && strings.HasPrefix(value, persistent.current) && len(persistent.current) > 0 {
			// filter on previous results
			newResults := make([]string, 0, len(persistent.results))
			for _, r := range persistent.results {
				if strings.HasPrefix(r, value) {
					newResults = append(newResults, r)
				}
			}
			persistent.results = newResults
		} else {
			// keep allocated slice memory and find initial filter
			persistent.results = persistent.results[:0]
			for _, s := range input {
				if strings.HasPrefix(s, value) {
					persistent.results = append(persistent.results, s)
				}
			}
		}

		persistent.current = value
		return persistent.results
	}
}

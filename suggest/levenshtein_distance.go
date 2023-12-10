package suggest

import (
	"strings"
	"sync"
)

type ldOpts struct {
	ignoreCase  bool
	minDistance int
	maxDistance int
}

type safeResults struct {
	mux     sync.Mutex
	current string
	results []string
}

// LevenshteinDistanceOpt is a set of options for use with LevenshteinDistance
type LevenshteinDistanceOpt func(o *ldOpts)

// LevenshteinDistanceMin returns a LevenshteinDistanceOpt that
// sets the minimum Levenshtein distance on an ldOpts instance.
//
// Parameters:
//   - minimum (int): is the minimum distance value.
//
// Returns: a function that accepts an internal option instance and sets its minDistance.
func LevenshteinDistanceMin(minimum int) LevenshteinDistanceOpt {
	return func(o *ldOpts) {
		if minimum >= 0 {
			o.minDistance = minimum
		}
	}
}

// LevenshteinDistanceMax returns a LevenshteinDistanceOpt that
// sets the maximum Levenshtein distance on an ldOpts instance.
//
// Parameters:
//   - maximum (int): is the maximum distance value.
//
// Returns: a function that accepts an internal option instance and sets its maxDistance.
func LevenshteinDistanceMax(maximum int) LevenshteinDistanceOpt {
	return func(o *ldOpts) {
		o.maxDistance = maximum
	}
}

// LevenshteinDistanceIgnoreCase returns a LevenshteinDistanceOpt that
// sets the ignoreCase property on an ldOpts instance.
//
// Parameters:
//   - ignoreCase (bool): is the flag to determine whether to ignore case.
//
// Returns: a function that accepts an internal option instance and sets its ignoreCase.
func LevenshteinDistanceIgnoreCase(ignoreCase bool) LevenshteinDistanceOpt {
	return func(o *ldOpts) {
		o.ignoreCase = ignoreCase
	}
}

// LevenshteinDistance is a function
func LevenshteinDistance(data []string, options ...LevenshteinDistanceOpt) Completion {
	opts := ldOpts{
		ignoreCase:  true,
		minDistance: 0,
		maxDistance: 5,
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

		if len(value) > 0 && persistent.current != value {
			// keep allocated slice memory
			persistent.results = persistent.results[:0]
			for _, s := range input {
				distance := calculateLevenshteinDistance(s, value, opts.ignoreCase)
				if opts.minDistance <= distance && distance <= opts.maxDistance {
					persistent.results = append(persistent.results, s)
				}
			}
			persistent.current = value
			return persistent.results
		} else if persistent.current == value {
			return persistent.results
		} else {
			persistent.results = persistent.results[:0]
		}

		persistent.current = value
		return persistent.results
	}
}

// calculateLevenshteinDistance is a Levenshtein Distance implementation for finding the edit distance between two strings.
// For more details, see: https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Levenshtein_distance#Go
// For more of an explanation, see: https://www.baeldung.com/cs/levenshtein-distance-computation
func calculateLevenshteinDistance(source, target string, ignoreCase bool) int {
	sourceRunes := stringToRunes(source, ignoreCase)
	targetRunes := stringToRunes(target, ignoreCase)
	column := make([]int, len(sourceRunes)+1)

	for index := 1; index <= len(sourceRunes); index++ {
		column[index] = index
	}

	for i := 1; i <= len(targetRunes); i++ {
		previousDiagonal := i - 1
		column[0] = i
		for j := 1; j <= len(sourceRunes); j++ {
			currentDiagonal := column[j]
			cost := 0
			if sourceRunes[j-1] != targetRunes[i-1] {
				cost = 1
			}
			column[j] = minimum(minimum(column[j]+1, column[j-1]+1), previousDiagonal+cost)
			previousDiagonal = currentDiagonal
		}
	}

	return column[len(sourceRunes)]
}

func stringToRunes(str string, ignoreCase bool) []rune {
	if ignoreCase {
		str = strings.ToLower(str)
	}
	return []rune(str)
}

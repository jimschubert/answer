package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func errMessage(args ...any) string {
	switch a := args; {
	case a == nil:
		return ""
	case len(a) == 0:
		return ""
	case len(a) == 1:
		switch msg := a[0].(type) {
		case string:
			return msg
		default:
			return fmt.Sprintf("%+v", msg)
		}
	default:
		return fmt.Sprintf(a[0].(string), a[1:]...)
	}
}

// Func determines if the input string is valid, returning nil if valid or an error if invalid
type Func func(input string) error

// MinLength defines the minimum allowed length (in runes)
func (fn Func) MinLength(length int, msgAndArgs ...any) Func {
	return func(input string) error {
		actual := len([]rune(input))
		if actual < length {
			if msg := errMessage(msgAndArgs...); msg != "" {
				return errors.New(msg)
			}
			return fmt.Errorf("minimum length required=%d actual=%d", length, actual)
		}

		return fn(input)
	}
}

// MaxLength defines the maximum allowed length (in runes)
func (fn Func) MaxLength(length int, msgAndArgs ...any) Func {
	return func(input string) error {
		actual := len([]rune(input))
		if actual > length {
			if msg := errMessage(msgAndArgs...); msg != "" {
				return errors.New(msg)
			}
			return fmt.Errorf("maximum length allowed=%d actual=%d", length, actual)
		}

		return fn(input)
	}
}

// Matches defines the pattern required for the targeted input
func (fn Func) Matches(pattern string, msgAndArgs ...any) Func {
	re := regexp.MustCompile(pattern)
	return func(input string) error {
		if !re.MatchString(input) {
			if msg := errMessage(msgAndArgs...); msg != "" {
				return errors.New(msg)
			}
			return errors.New("❌invalid input")
		}
		return fn(input)
	}
}

// And allows chaining another *required* validation function to the end of other functions in the chain
func (fn Func) And(other Func) Func {
	return func(input string) error {
		return errors.Join(fn(input), other(input))
	}
}

// Or allows defining a different validation to invoke when the first validation was successful.
func (fn Func) Or(other Func) Func {
	return func(input string) error {
		err := fn(input)
		if err == nil {
			err = other(input)
		}
		return err
	}
}

// Contains defines a substring which is required in the target input
func (fn Func) Contains(value string, msgAndArgs ...any) Func {
	return func(input string) error {
		if !strings.Contains(input, value) {
			if msg := errMessage(msgAndArgs...); msg != "" {
				return errors.New(msg)
			}
			return fmt.Errorf("input does not contain %q", value)
		}
		return fn(input)
	}
}

// Build returns the raw underlying functional type
func (fn Func) Build() func(string) error {
	return fn
}

func emptyValidateFunc(_ string) error {
	return nil
}

// NewValidation creates the initial chain for validations
func NewValidation() Func {
	return emptyValidateFunc
}

package validate

import (
	"errors"
	"strings"
	"testing"
	"unicode"

	_ "github.com/charmbracelet/x/exp/teatest"
)

func TestNewValidation_all(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		validationFn Func
		want         error
	}{
		{
			name:         "NewValidation() returns no errors for empty input",
			input:        "",
			validationFn: NewValidation(),
			want:         nil,
		},
		{
			name:         "NewValidation() returns no errors for non-empty input",
			input:        "asdf",
			validationFn: NewValidation(),
			want:         nil,
		},
		{
			name:         "MinLength() returns no errors for valid input",
			input:        "asdf",
			validationFn: NewValidation().MinLength(3),
			want:         nil,
		},
		{
			name:         "MinLength() returns error based on rune count",
			input:        "簡",
			validationFn: NewValidation().MinLength(2),
			want:         errors.New("minimum length required=2 actual=1"),
		},
		{
			name:         "MinLength() returns error for invalid input",
			input:        "asdf",
			validationFn: NewValidation().MinLength(5),
			want:         errors.New("minimum length required=5 actual=4"),
		},
		{
			name:         "MinLength() returns custom error for invalid input",
			input:        "asdf",
			validationFn: NewValidation().MinLength(5, "Input must be at least 5 characters"),
			want:         errors.New("Input must be at least 5 characters"),
		},
		{
			name:         "MaxLength() returns no errors for valid input",
			input:        "asdf",
			validationFn: NewValidation().MaxLength(10),
			want:         nil,
		},
		{
			name:         "MaxLength() returns error based on rune count",
			input:        "シンプルさ",
			validationFn: NewValidation().MaxLength(2),
			want:         errors.New("maximum length allowed=2 actual=5"),
		},
		{
			name:         "MaxLength() returns error for invalid input",
			input:        "asdfasdf",
			validationFn: NewValidation().MaxLength(5),
			want:         errors.New("maximum length allowed=5 actual=8"),
		},
		{
			name:         "MaxLength() returns custom error for invalid input",
			input:        "asdfasdf",
			validationFn: NewValidation().MaxLength(5, "Input must be at most 5 characters"),
			want:         errors.New("Input must be at most 5 characters"),
		},
		{
			name:         "Matches() returns no errors for valid input",
			input:        "asdf",
			validationFn: NewValidation().Matches(".*"),
			want:         nil,
		},
		{
			name:         "Matches() returns error for invalid input",
			input:        "asdfasdf",
			validationFn: NewValidation().Matches(`\d{1,}`),
			want:         errors.New("❌invalid input"),
		},
		{
			name:         "Matches() returns custom error for invalid input",
			input:        "asdfasdf",
			validationFn: NewValidation().Matches(`\d{1,}`, "Numbers only!"),
			want:         errors.New("Numbers only!"),
		},
		{
			name:         "Contains() returns no errors for valid input",
			input:        "The quick brown fox",
			validationFn: NewValidation().Contains("quick"),
			want:         nil,
		},
		{
			name:         "Contains() returns error for invalid input",
			input:        "asdf",
			validationFn: NewValidation().Contains("quick"),
			want:         errors.New(`input does not contain "quick"`),
		},
		{
			name:         "Contains() returns custom error for invalid input",
			input:        "asdf",
			validationFn: NewValidation().Contains("quick", "Use standard font test string"),
			want:         errors.New("Use standard font test string"),
		},
		{
			name:         "MinLength().MaxLength() returns no errors for valid input",
			input:        "asdf",
			validationFn: NewValidation().MinLength(1).MaxLength(10),
			want:         nil,
		},
		{
			name:         "MinLength().MaxLength() returns error for invalid input exceeding maximum",
			input:        "asdfasdf",
			validationFn: NewValidation().MinLength(3).MaxLength(5),
			want:         errors.New("maximum length allowed=5 actual=8"),
		},
		{
			name:         "MinLength().MaxLength() returns error for invalid input below minimum",
			input:        "as",
			validationFn: NewValidation().MinLength(3).MaxLength(5),
			want:         errors.New("minimum length required=3 actual=2"),
		},
		{
			name:         "MinLength().MaxLength() returns custom error for invalid input on minimum",
			input:        "asd",
			validationFn: NewValidation().MinLength(5, "Min length is 5").MaxLength(10, "Max length is 10"),
			want:         errors.New("Min length is 5"),
		},
		{
			name:         "MinLength().MaxLength() returns custom error for invalid input on maximum",
			input:        "asdfasdfasdf",
			validationFn: NewValidation().MinLength(5, "Min length is 5").MaxLength(10, "Max length is 10"),
			want:         errors.New("Max length is 10"),
		},
		{
			name:  "And() returns no errors for valid input",
			input: "asdf",
			validationFn: NewValidation().Matches(".*").And(func(input string) error {
				if input == "asdf" {
					return nil
				}
				return errors.New("invalid")
			}),
			want: nil,
		},
		{
			name:  "And() returns single error for invalid input failing one condition",
			input: "12334",
			validationFn: NewValidation().Matches(`\d{1,}`).And(func(input string) error {
				if input == "12334" {
					return errors.New("Don't just slap the keyboard for test input.")
				}
				return nil
			}),
			want: errors.New("Don't just slap the keyboard for test input."),
		},
		{
			name:  "And() returns combined error for invalid input failing both conditions",
			input: "asdfasdf",
			validationFn: NewValidation().Matches(`\d{1,}`).And(func(input string) error {
				if input == "asdfasdf" {
					return errors.New("Don't just slap the keyboard for test input.")
				}
				return nil
			}),
			want: errors.Join(errors.New("❌invalid input"), errors.New("Don't just slap the keyboard for test input.")),
		},
		{
			name:  "And() returns custom errors for invalid input",
			input: "as",
			validationFn: NewValidation().Matches(`\d{1,}`, "Numbers only!").And(func(input string) error {
				return errors.New("Always error")
			}).And(NewValidation().MinLength(3, "100 or more")),
			want: errors.Join(
				errors.New("Numbers only!"),
				errors.New("Always error"),
				errors.New("100 or more"),
			),
		},
		{
			name:  "Or() does not invoke when first condition fails",
			input: "asdf",
			validationFn: NewValidation().Matches(`\d{1,}`).Or(func(input string) error {
				return errors.New("Should never happen")
			}),
			want: errors.New("❌invalid input"),
		},
		{
			name:  "Or() invokes only when first condition is successful",
			input: "1234",
			validationFn: NewValidation().Matches(`\d{1,}`).Or(func(input string) error {
				if input == "1234" {
					return errors.New("Don't just slap the keyboard for test input.")
				}
				return nil
			}),
			want: errors.New("Don't just slap the keyboard for test input."),
		},
		{
			name:  "complex scenario passing",
			input: "A man a plan a canal Panama",
			validationFn: NewValidation().
				MinLength(5).
				MaxLength(30).
				Contains("Panama", "Uppercase proper nouns").
				And(func(input string) error {
					forward := strings.Builder{}
					for _, r := range input {
						if unicode.IsSpace(r) || !unicode.IsLetter(r) {
							continue
						}
						if unicode.IsLower(r) {
							forward.WriteRune(r)
						} else {
							forward.WriteRune(unicode.ToLower(r))
						}
					}
					reversed := strings.Builder{}
					rs := []rune(input)
					maxLen := len(rs) - 1
					for i := range rs {
						r := rs[maxLen-i]
						if unicode.IsSpace(r) || !unicode.IsLetter(r) {
							continue
						}
						if unicode.IsLower(r) {
							reversed.WriteRune(r)
						} else {
							reversed.WriteRune(unicode.ToLower(r))
						}
					}

					if reversed.String() != forward.String() {
						return errors.New("not a palindrome")
					}
					return nil
				}),
			want: nil,
		},
		{
			name: "complex scenario failing",
			// note the typo on canal
			input: "A man a plan a canale Panama",
			validationFn: NewValidation().
				MinLength(5).
				MaxLength(30).
				Contains("Panama", "Uppercase proper nouns").
				And(func(input string) error {
					forward := strings.Builder{}
					for _, r := range input {
						if unicode.IsSpace(r) || !unicode.IsLetter(r) {
							continue
						}
						if unicode.IsLower(r) {
							forward.WriteRune(r)
						} else {
							forward.WriteRune(unicode.ToLower(r))
						}
					}
					reversed := strings.Builder{}
					rs := []rune(input)
					maxLen := len(rs) - 1
					for i := range rs {
						r := rs[maxLen-i]
						if unicode.IsSpace(r) || !unicode.IsLetter(r) {
							continue
						}
						if unicode.IsLower(r) {
							reversed.WriteRune(r)
						} else {
							reversed.WriteRune(unicode.ToLower(r))
						}
					}

					if reversed.String() != forward.String() {
						return errors.New("not a palindrome")
					}
					return nil
				}),
			want: errors.New("not a palindrome"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.validationFn(tt.input)
			// we can't use reflect.DeepEqual for any errors returning via errors.Join because it makes the test case uglier.
			// so, instead we make the conditional check a little uglier.
			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) || (got != nil && tt.want != nil && got.Error() != tt.want.Error()) {
				t.Errorf("validations: got %v, want %v", got, tt.want)
			}
		})
	}
}

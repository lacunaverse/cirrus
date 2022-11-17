package cirrus

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	_ "embed"

	"github.com/araddon/dateparse"
	"github.com/hvlck/txt"
)

var (
	//go:embed data/cardinality.txt
	CARDINALITY_DICT_EMBED string
	CARDINALITY_DICT       = LoadDictionary(CARDINALITY_DICT_EMBED)
)

// Loads a dictionary from an embeded newline-delimited dictionary file.
// The dictionary words are placed in a trie.
func LoadDictionary(data string) *txt.Node {
	t := txt.NewTrie()

	for _, v := range strings.Split(data, "\n") {
		t.Insert(v, nil)
	}

	return t
}

type Unit int

// Recognized units.
var Units = [][]string{
	NO_UNIT: {"none"},
	// metric
	METERS: {"meters", "m"},

	// imperial
	INCHES: {"inches", "in"},
	FEET:   {"feet", "ft", "foot"},
	MILES:  {"mile", "mi"},

	TIME: {"minute", "second", "hour", "day"},
}

func (u Unit) String() string {
	return Units[u][0]
}

const (
	NO_UNIT = iota
	// metric
	METERS

	INCHES
	FEET
	MILES
	// todo: other units

	TIME
)

type ResultType int

var results = []string{
	LINK:     "link",
	QUANTITY: "quantity",
	DATE:     "date",
	ORG:      "org",
	CARDINAL: "cardinal",
	MONEY:    "monetary",
	EVENT:    "event",
}

func (r ResultType) String() string {
	return results[r]
}

const (
	NONE = iota
	LINK
	QUANTITY
	DATE
	ORG
	CARDINAL
	MONEY
	EVENT
)

type dataType interface {
	dataType()
}

type UnitValue struct {
	Value Unit
}

func (u *UnitValue) dataType() {}

type Result struct {
	// type of result
	ResultType ResultType `json:"type"`
	// title/label of the unit, as in with graphs
	// e.g. "Number of Queries Per Second"
	Label string `json:"label"`
	// Custom data
	Data dataType `json:"data"`
	// value
	Value      string `json:"value"`
	Start, End uint
}

func matchProper(s string) {

}

// Determines whether a token is a quantity or not
func hasUnit(token string) (Unit, bool) {
	token = strings.ToLower(token)
	for i, abbreviations := range Units {
		if i == NONE {
			continue
		}

		for _, name := range abbreviations {
			if len(name) <= 2 {
				if token == name {
					return Unit(i), true
				}

				continue
			}

			if strings.HasPrefix(token, name) {
				return Unit(i), true
			}
		}
	}

	return NONE, false
}

// Matches any sequence of one or more numbers
var NUM_REGEXP = regexp.MustCompile(`\d+`)

// Matches a single number
var SINGLE_NUMBER_REGEXP = regexp.MustCompile(`\d`)

// Matches one or more letters, used as a filter for determining the unit
var UNIT_TYPE_REGEXP = regexp.MustCompile(`[a-zA-Z]+`)

// Matches one or more numbers followed by optional whitespace and a sequence of one or more letters, used for
// determining whether a string may be a quantity or not.
var UNIT_EXTRACT_REGEXP = regexp.MustCompile(`\d+\s?[A-Z..a-z]+`)

// Attempts to extract a unit from a given string.
// Matches take the form of `{num}{unit}`, with optional whitespace in between.
// Examples: `10 feet` `1ft` `1200inches`
func extractUnit(s string) (Unit, string) {
	if UNIT_EXTRACT_REGEXP.MatchString(s) {
		v := UNIT_EXTRACT_REGEXP.FindString(s)

		if NUM_REGEXP.MatchString(v) && UNIT_TYPE_REGEXP.MatchString(v) {
			unit := UNIT_TYPE_REGEXP.FindString(v)
			value := NUM_REGEXP.FindString(v)

			unitInt, err := strconv.Atoi(unit)
			if err != nil {
				return Unit(unitInt), value
			}
		}
	}

	return NONE, ""
}

func isUrl(token string) *Result {
	if strings.HasPrefix(token, "http") {
		if u, ok := url.Parse(token); ok == nil {
			r := &Result{
				ResultType: LINK,
				Value:      u.String(),
			}
			return r
		}
	}

	return nil
}

var currencies = map[string]bool{
	"$": true,
	"¥": true,
	"£": true,
	"€": true,
}

var (
	ErrNoExtract = errors.New("couldn't determine meaning")
)

func splitter(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsControl(r) || unicode.IsSpace(r) || r == ';' || r == ',' || r == '!'
	})
}

// Does nothing, just gets around a bug in txt with passing a nil tokenizer to txt.Tokenize.
var noTokenizers txt.Tokenizer = func(tokens []string) []string { return tokens }
var noNormalizer txt.Normalizer = func(s string) string { return s }

// Recognizes entities within a piece of natural text.
func Recognize(text string) ([]*Result, error) {
	results := []*Result{}

	tokens := txt.Tokenize(text, splitter, []txt.Normalizer{noNormalizer}, noTokenizers)
	for idx, v := range tokens {
		// checks if token is URL
		url := isUrl(v)
		if url != nil {
			results = append(results, url)
			continue
		}

		// check for dates
		// todo: fix/check timezone implementation
		// != 1 test case for decimal numbers. These take precedence over dates.
		if t, ok := dateparse.ParseAny(v); ok == nil && strings.Count(v, ".") != 1 {
			r := &Result{
				Label: "",
				Data: &UnitValue{
					Value: 0,
				},
				ResultType: DATE,
				Value:      t.String(),
			}

			results = append(results, r)

			continue
		}

		ch := rune(v[0])
		if unicode.IsUpper(ch) {
			// if prop, ok := MatchProper(v); ok {

			// }
		} else if unicode.IsNumber(ch) {
			i := 0
			for _, char := range v {
				if unicode.IsNumber(char) {
					i++
					continue
				}

				if char == '-' || unicode.ToLower(char) == 'e' {
					continue
				}

				break
			}

			if len(v) >= i+1 {
				i = i + 1
			}

			r := &Result{
				ResultType: QUANTITY,
				Value:      v[:i],
			}

			if len(tokens) > idx+1 {
				if un, ok := hasUnit(tokens[idx+1]); ok {
					r := &Result{
						ResultType: QUANTITY,
						Data: &UnitValue{
							Value: un,
						},
						Value: v[:i],
					}

					results = append(results, r)

					continue
				}
			}

			results = append(results, r)

			continue
		}

		if CARDINALITY_DICT.Contains(v) {
			r := &Result{
				ResultType: CARDINAL,
				Value:      v,
			}

			results = append(results, r)
			continue
		}

		if _, ok := currencies[string(v[0])]; ok {
			i := 0
			for _, char := range v[1:] {
				if unicode.IsNumber(char) {
					i++
					continue
				}

				break
			}

			if len(v) >= i+1 {
				i = i + 1
			}

			r := &Result{
				ResultType: MONEY,
				Value:      string(v[0]),
			}

			quantity := &Result{
				ResultType: QUANTITY,
				Value:      v[0:i],
			}

			results = append(results, r, quantity)
			continue
		}

		r := &Result{
			ResultType: NONE,
			Value:      v,
		}
		results = append(results, r)
	}

	return results, nil
}

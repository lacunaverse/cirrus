package cirrus

import (
	"bytes"
	"errors"
	"io/ioutil"
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

var Units = [][]string{
	NO_UNIT: {"none"},
	// metric
	METERS: {"meters", "m"},

	// imperial
	FEET: {"feet", "ft"},
}

func (u Unit) String() string {
	return Units[u][0]
}

const (
	NO_UNIT = iota
	// metric
	METERS

	FEET
	// todo: other units
)

type ResultType int

var results = []string{
	LINK:     "link",
	QUANTITY: "quantity",
	DATE:     "date",
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

type Result struct {
	// type of result
	ResultType ResultType `json:"type"`
	// title/label of the unit, as in with graphs
	// e.g. "Number of Queries Per Second"
	Label string `json:"label"`
	// unit, if applicable
	Unit Unit `json:"unit"`
	// value
	Value string `json:"value"`
	Start uint
	End   uint
}

func hasUnit(s string) (Unit, bool) {
	s = strings.ToLower(s)

// Determines whether a token is a quantity or not
func hasUnit(token string) (Unit, bool) {
	token = strings.ToLower(token)
	for i, v := range Units {
		if i == NONE {
			continue
		}

		for _, name := range v {
			if strings.HasSuffix(token, name) {
				return Unit(i), true
			}
		}
	}

	return NO_UNIT, false
}

var NUM_REGEXP = regexp.MustCompile(`\d+`)
var SINGLE_NUMBER_REGEXP = regexp.MustCompile(`\d`)
var UNIT_TYPE_REGEXP = regexp.MustCompile(`[a-zA-Z]+`)
var UNIT_EXTRACT_REGEXP = regexp.MustCompile(`\d+\s?[A-Z..a-z]+`)

func extractUnit(s string) (Unit, string) {
	if UNIT_EXTRACT_REGEXP.Match([]byte(s)) {
		v := UNIT_EXTRACT_REGEXP.FindString(s)
		vb := []byte(v)

		if NUM_REGEXP.Match(vb) && UNIT_TYPE_REGEXP.Match(vb) {
			unit := UNIT_TYPE_REGEXP.FindString(v)
			value := NUM_REGEXP.FindString(v)

			uInt, err := strconv.Atoi(unit)
			if err != nil {
				return Unit(uInt), value
			}
		}
	}

	return 0, ""
}

var (
	ErrNoExtract = errors.New("couldn't determine meaning")
)
func Recognize(text string) (*Result, error) {
	if strings.HasPrefix(text, "http") {
		if u, ok := url.Parse(text); ok == nil {
			return &Result{
				ResultType: LINK,
				Value:      u.String(),
			}, nil
		}
	}

	// check for dates
	// todo: fix/check timezone implementation
	if t, ok := dateparse.ParseAny(text); ok == nil {
		return &Result{
			Label:      "",
			Unit:       0,
			ResultType: DATE,
			Value:      t.String(),
		}, nil
	}

	if u, ok := hasUnit(text); ok {
		_, val := extractUnit(text)

		return &Result{
			ResultType: QUANTITY,
			Unit:       u,
			Value:      val,
		}, nil
	}

	return &Result{}, ErrNoExtract
}
}

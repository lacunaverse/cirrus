package synapse

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Unit int

var units = []string{
	METERS: "m",
}

func (u Unit) String() string {
	return units[u]
}

const (
	METERS = iota
	// todo: other units
)

type ResultType int

var results = []string{}

func (r ResultType) String() string {
	return results[r]
}

const (
	LINK = iota
	QUANTITY
	DATE
)

//
type Result struct {
	// type of result
	ResultType ResultType
	// title/label of the unit, as in with graphs
	// e.g. "Number of Queries Per Second"
	Label string
	// unit, if applicable
	Unit Unit
	// value
	Value string
}

func hasUnit(s string) (Unit, bool) {
	s = strings.ToLower(s)
	for i, v := range units {
		if strings.Contains(s, v) {
			return Unit(i), true
		}
	}

	return -1, false
}

var NUM_REGEXP = regexp.MustCompile(`\d+`)
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
		fmt.Println(v)
	}

	return 0, ""
}

func Determine(s string) (Result, error) {
	if u, ok := url.Parse(s); ok == nil {
		return Result{
			ResultType: LINK,
			Value:      u.String(),
		}, nil
	}

	if u, ok := hasUnit(s); ok {
		_, val := extractUnit(s)

		return Result{
			ResultType: QUANTITY,
			Unit:       u,
			Value:      val,
		}, nil
	}

	return Result{}, errors.New("couldn't determine meaning")
}

package cirrus

import (
	"testing"
	"time"

	"github.com/hvlck/txt"
)

func TestTokenizer(t *testing.T) {
	examples := []string{
		"https://google.com is a search engine",
		"length of 20m",
		"on 2/11/2015, something happened",
	}

	expected := [][]string{
		{"https://google.com", "is", "a", "search", "engine"},
		{"length", "of", "20m"},
		{"on", "2/11/2015", "something", "happened"},
	}

	for index, v := range examples {
		tokenized := txt.Tokenize(v, splitter, []txt.Normalizer{noNormalizer}, noTokenizers)
		for idx, tok := range tokenized {
			if expected[index][idx] != tok {
				t.Fatalf("expected %v, got %v", expected[index][idx], tok)
			}
		}
	}
}

func TestRecognize(t *testing.T) {
	examples := []string{
		"https://google.com extra tokens",
		"length of 20m",
		"on 2/11/2015, something happened",
		"two dozen",
		"$20",
		"10.22",
		"10e12",
		"5e-3",
	}

	output := [][]Result{
		{
			{
				ResultType: LINK,
				Value:      "https://google.com/",
			},
			{
				ResultType: NONE,
				Value:      "extra",
			},
			{
				ResultType: NONE,
				Value:      "tokens",
			},
		},
		{
			{
				ResultType: NONE,
				Value:      "length",
			},
			{
				ResultType: NONE,
				Value:      "of",
			},
			{
				ResultType: QUANTITY,
				Value:      "20",
				Data:       &UnitValue{Value: METERS},
			},
		},
		{
			{
				ResultType: NONE,
				Value:      "on",
			},
			{
				ResultType: DATE,
				Value:      time.Date(2015, 2, 11, 0, 0, 0, 0, time.UTC).String(),
			},
			{
				ResultType: NONE,
				Value:      "something",
			},
			{
				ResultType: NONE,
				Value:      "happened",
			},
		},
		{
			{
				ResultType: CARDINAL,
				Value:      "two",
			},
			{
				ResultType: CARDINAL,
				Value:      "dozen",
			},
		},
		{
			{
				ResultType: MONEY,
				Value:      "$",
			},
			{
				ResultType: QUANTITY,
				Value:      "20",
			},
		},
		{
			{
				ResultType: QUANTITY,
				Value:      "10.22",
			},
		},
		{
			{
				ResultType: QUANTITY,
				Value:      "10e12",
			},
		},
		{
			{
				ResultType: QUANTITY,
				Value:      "5e-3",
			},
		},
	}

	for index, v := range examples {
		expected := output[index]
		result, err := Recognize(v)
		if err != nil {
			t.Fatal(err)
		}

		if len(result) != len(expected) {
			t.Fatalf("expected length of result to be %v, got %v for example '%v'", len(expected), len(result), v)
		}

		for idx, token := range result {
			expTok := expected[idx]
			if token.ResultType != expTok.ResultType || token.Label != expTok.Label {
				t.Fatalf("expected %v, got %v for example '%v'", expTok, token, v)
			}
		}
	}
}

func TestURL(t *testing.T) {
	v, err := Recognize("https://github.com/")
	if err != nil {
		t.Fail()
	}

	if v[0].Value != "https://github.com/" && v[0].ResultType != LINK {
		t.Fail()
	}
}

func TestHasUnit(t *testing.T) {
	units := []struct {
		string
		Unit
	}{
		{"m", METERS},
		{"meters", METERS},
		{"ft", FEET},
	}

	for _, v := range units {
		unit, exists := hasUnit(v.string)
		if !exists || v.Unit != unit {
			t.Fail()
		}
	}
}

func TestExtractUnit(t *testing.T) {
	un, value := extractUnit("length of 10m")
	if un != METERS && value != "10" {
		t.Fail()
	}
}

func TestDate(t *testing.T) {
	d, err := Recognize("3/12/21")

	if err != nil || d[0].ResultType != DATE {
		t.Fail()
	}
}

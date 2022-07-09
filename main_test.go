package synapse

import (
	"testing"
)

func TestURL(t *testing.T) {
	v, err := Determine("https://github.com/")
	if err != nil {
		t.Fail()
	}

	if v.Value != "https://github.com/" && v.ResultType != LINK {
		t.Fail()
	}
}

func TestHasUnit(t *testing.T) {
	v, exists := hasUnit("length of 10m")

	if !exists || v != METERS {
		t.Fail()
	}
}

func TestExtractUnit(t *testing.T) {
	un, value := extractUnit("length of 10m")
	if un != METERS && value != "10" {
		t.Fail()
	}
}

func TestDate(t *testing.T) {
	d, err := Determine("3/12/21")

	if err != nil || d.ResultType != DATE {
		t.Fail()
	}
}

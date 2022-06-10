package utype

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBooleanJson = []struct {
	input    string
	expected bool
	isErr    bool
}{
	{`true`, true, false},
	{`false`, false, false},
	{`"TRue"`, true, false},
	{`"False"`, false, false},
	{`1`, true, false},
	{`0`, false, false},
	{`"1"`, true, false},
	{`"0"`, false, false},
	{`2`, false, true},
	{`"satu"`, false, true},
}

var testBoleanScan = []struct {
	input    interface{}
	expected bool
	isErr    bool
}{
	{true, true, false},
	{false, false, false},
	{"TRue", true, false},
	{"False", false, false},
	{1, true, false},
	{0, false, false},
	{"1", true, false},
	{"0", false, false},
	{2, false, true},
	{"satu", false, true},
}

func TestBoolean_UnmarshalJSON(t *testing.T) {
	for _, test := range testBooleanJson {
		t.Logf("testing with input: %v", test.input)
		expected := BooleanFrom(test.expected)
		if test.isErr {
			expected.Valid = false
		}
		var b Boolean
		err := json.Unmarshal([]byte(test.input), &b)
		if test.isErr && err == nil {
			t.Errorf("expecting an error from unmarshal of string %q. %s", test.input, err)
			continue
		}

		if !test.isErr && err != nil {
			t.Errorf(err.Error())
			continue
		}

		assert.Equal(t, expected, b, "expecting %v, got %v", expected, b)
	}
}

func TestBoolean_Scan(t *testing.T) {
	for _, test := range testBoleanScan {
		t.Logf("testing scan with input: %v", test.input)
		expected := BooleanFrom(test.expected)
		if test.isErr {
			expected.Valid = false
		}

		var b Boolean
		err := b.Scan(test.input)
		if test.isErr && err == nil {
			t.Errorf("expecting an error from scan of value %q. %s", test.input, err)
			continue
		}

		if !test.isErr && err != nil {
			t.Errorf(err.Error())
			continue
		}

		assert.Equal(t, expected, b, "expecting %v, got %v", expected, b)
	}
}

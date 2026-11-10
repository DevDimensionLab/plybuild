package canonicaljson

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestNumberPolicyIsExplicitScopedAndStrict(t *testing.T) {
	allowDuration := func(path []string, number string) (float64, error) {
		if strings.Join(path, "/") != "steps/0/duration_seconds" {
			return 0, fmt.Errorf("number outside selected field")
		}
		return strconv.ParseFloat(number, 64)
	}
	input := []byte(`{"steps":[{"duration_seconds":0.125}],"sequence":1}`)
	if _, err := DecodeStrict(input); err == nil {
		t.Fatal("default contract accepted a fraction")
	}
	value, err := DecodeStrictWithNumberPolicy(input, allowDuration)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Marshal(value)
	if err != nil || string(b) != `{"sequence":1,"steps":[{"duration_seconds":0.125}]}` {
		t.Fatalf("selected finite number or integer changed: %s %v", b, err)
	}
	for _, input := range []string{
		`{"steps":[{"duration_seconds":0.125}],"sequence":1.5}`,
		`{"steps":[{"duration_seconds":0.125,"duration_seconds":0.25}]}`,
		`{"steps":[{"duration_seconds":0.125}],"text":"\uD800"}`,
		`{"steps":[{"duration_seconds":0.125}]} {}`,
		`{"steps":[{"duration_seconds":0.125}],"sequence":9223372036854775808}`,
	} {
		if _, err := DecodeStrictWithNumberPolicy([]byte(input), allowDuration); err == nil {
			t.Errorf("scoped policy bypassed strict validation: %s", input)
		}
	}
	for _, number := range []float64{math.Inf(1), math.NaN()} {
		if _, err := DecodeStrictWithNumberPolicy([]byte(`0.125`), func([]string, string) (float64, error) { return number, nil }); err == nil {
			t.Fatal("non-finite policy output accepted")
		}
	}
}

func TestNumberPolicyCanonicalizesNegativeZero(t *testing.T) {
	value, err := DecodeStrictWithNumberPolicy([]byte(`-0.0`), func(_ []string, number string) (float64, error) { return strconv.ParseFloat(number, 64) })
	if err != nil {
		t.Fatal(err)
	}
	b, err := Marshal(value)
	if err != nil || string(b) != "0" {
		t.Fatalf("explicit numeric policy did not canonicalize zero: %s %v", b, err)
	}
}

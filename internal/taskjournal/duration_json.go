package taskjournal

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// Fractional numbers belong only to the existing snapshot duration field.
// Contributions, native evidence and every other snapshot field retain the
// strict integer domain. The preliminary kind lookup changes no source bytes;
// duplicate keys, UTF-8, escapes and trailing values are still strictly checked.
func decodeJournalJSON(raw []byte) (canonicaljson.Value, error) {
	var header struct {
		Kind  string            `json:"kind"`
		Steps []json.RawMessage `json:"steps"`
	}
	if json.Unmarshal(raw, &header) != nil || header.Kind != "ply.workspace.task-journal-snapshot" {
		return canonicaljson.DecodeStrict(raw)
	}
	return canonicaljson.DecodeStrictWithNumberPolicy(raw, func(path []string, number string) (float64, error) {
		if len(path) != 3 || path[0] != "steps" || path[2] != "duration_seconds" {
			return 0, fmt.Errorf("unsupported JSON number %q outside snapshot step duration_seconds", number)
		}
		if _, err := strconv.ParseUint(path[1], 10, 64); err != nil {
			return 0, fmt.Errorf("invalid snapshot duration path")
		}
		// Measured time.Duration values need at most a short float64 spelling.
		// Bound work before exact decimal comparison of untrusted offline input.
		if len(number) > 64 {
			return 0, fmt.Errorf("snapshot duration_seconds number is too long")
		}
		value, err := strconv.ParseFloat(number, 64)
		if err != nil || math.IsInf(value, 0) || math.IsNaN(value) || value < 0 || value == 0 && math.Signbit(value) {
			return 0, fmt.Errorf("invalid snapshot duration_seconds number %q", number)
		}
		if value == 0 {
			// Avoid constructing huge powers for zero with an arbitrary exponent,
			// and reject a nonzero mantissa that underflowed during ParseFloat.
			mantissa := strings.FieldsFunc(number, func(r rune) bool { return r == 'e' || r == 'E' })[0]
			if strings.Trim(mantissa, "0.") != "" {
				return 0, fmt.Errorf("snapshot duration_seconds loses precision as float64: %q", number)
			}
			return 0, nil
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			return 0, err
		}
		// Do not silently round an overprecise or underflowing input to a valid
		// measured duration. Equivalent spellings (including exponents) remain
		// canonical, while the exact decimal value must survive serialization.
		given, givenOK := new(big.Rat).SetString(number)
		canonical, canonicalOK := new(big.Rat).SetString(string(encoded))
		if !givenOK || !canonicalOK || given.Cmp(canonical) != 0 {
			return 0, fmt.Errorf("snapshot duration_seconds loses precision as float64: %q", number)
		}
		return value, nil
	})
}

// Package canonicaljson implements the deliberately small RFC 8785 JSON
// domain used by workflow handoffs. DecodeStrict restricts numbers to signed
// integers. A caller may explicitly provide a field-specific number policy.
package canonicaljson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

type Value any

type Member struct {
	Name  string
	Value Value
}

type Object []Member

// NumberPolicy handles numbers outside the default signed-integer domain.
// Paths contain object member names and decimal array indices. Returning an
// error rejects the number; the policy must explicitly validate its own scope
// and any precision constraints. Other decoding checks remain unchanged.
type NumberPolicy func(path []string, number string) (float64, error)

// policyNumber cannot be constructed outside this package; only an explicit
// number policy can introduce a finite fractional number into a decoded value.
type policyNumber float64

func DecodeStrict(input []byte) (Value, error) {
	return DecodeStrictWithNumberPolicy(input, nil)
}

func DecodeStrictWithNumberPolicy(input []byte, policy NumberPolicy) (Value, error) {
	if len(input) >= 3 && bytes.Equal(input[:3], []byte{0xef, 0xbb, 0xbf}) {
		return nil, errors.New("JSON BOM is not allowed")
	}
	if !utf8.Valid(input) {
		return nil, errors.New("JSON must be valid UTF-8")
	}
	if err := validateEscapedUnicode(input); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	value, err := decodeValue(decoder, nil, policy)
	if err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return nil, errors.New("JSON must contain exactly one value")
		}
		return nil, fmt.Errorf("trailing JSON: %w", err)
	}
	return value, nil
}

func decodeValue(decoder *json.Decoder, path []string, policy NumberPolicy) (Value, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	switch typed := token.(type) {
	case nil, bool, string:
		return typed, nil
	case json.Number:
		text := typed.String()
		if text == "-0" || bytes.ContainsAny([]byte(text), ".eE") {
			if policy != nil {
				value, err := policy(append([]string(nil), path...), text)
				if err != nil {
					return nil, err
				}
				if math.IsNaN(value) || math.IsInf(value, 0) {
					return nil, fmt.Errorf("number policy returned a non-finite value for %q", text)
				}
				if value == 0 {
					value = 0 // RFC 8785 renders either sign of zero as 0.
				}
				return policyNumber(value), nil
			}
			return nil, fmt.Errorf("unsupported JSON number %q", text)
		}
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unsupported JSON integer %q", text)
		}
		return value, nil
	case json.Delim:
		switch typed {
		case '{':
			object := Object{}
			seen := make(map[string]struct{})
			for decoder.More() {
				nameToken, err := decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("decode object member: %w", err)
				}
				name, ok := nameToken.(string)
				if !ok {
					return nil, errors.New("object member name must be a string")
				}
				if _, duplicate := seen[name]; duplicate {
					return nil, fmt.Errorf("duplicate object member %q", name)
				}
				seen[name] = struct{}{}
				child := path
				if policy != nil {
					child = append(path, name)
				}
				value, err := decodeValue(decoder, child, policy)
				if err != nil {
					return nil, err
				}
				object = append(object, Member{Name: name, Value: value})
			}
			closing, err := decoder.Token()
			if err != nil || closing != json.Delim('}') {
				return nil, errors.New("unterminated JSON object")
			}
			return object, nil
		case '[':
			array := []Value{}
			for decoder.More() {
				child := path
				if policy != nil {
					child = append(path, strconv.Itoa(len(array)))
				}
				value, err := decodeValue(decoder, child, policy)
				if err != nil {
					return nil, err
				}
				array = append(array, value)
			}
			closing, err := decoder.Token()
			if err != nil || closing != json.Delim(']') {
				return nil, errors.New("unterminated JSON array")
			}
			return array, nil
		}
	}
	return nil, fmt.Errorf("unsupported JSON token %T", token)
}

func Marshal(value Value) ([]byte, error) {
	var output bytes.Buffer
	if err := appendValue(&output, value); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func appendValue(output *bytes.Buffer, value Value) error {
	switch typed := value.(type) {
	case nil:
		output.WriteString("null")
	case bool:
		if typed {
			output.WriteString("true")
		} else {
			output.WriteString("false")
		}
	case string:
		if !utf8.ValidString(typed) {
			return errors.New("string value must be valid UTF-8")
		}
		appendString(output, typed)
	case int64:
		output.WriteString(strconv.FormatInt(typed, 10))
	case policyNumber:
		encoded, err := json.Marshal(float64(typed))
		if err != nil {
			return err
		}
		output.Write(encoded)
	case []Value:
		output.WriteByte('[')
		for index, element := range typed {
			if index > 0 {
				output.WriteByte(',')
			}
			if err := appendValue(output, element); err != nil {
				return err
			}
		}
		output.WriteByte(']')
	case Object:
		members := append(Object(nil), typed...)
		sort.SliceStable(members, func(i, j int) bool {
			return compareUTF16(members[i].Name, members[j].Name) < 0
		})
		output.WriteByte('{')
		for index, member := range members {
			if !utf8.ValidString(member.Name) {
				return errors.New("object member name must be valid UTF-8")
			}
			if index > 0 && members[index-1].Name == member.Name {
				return fmt.Errorf("duplicate object member %q", member.Name)
			}
			if index > 0 {
				output.WriteByte(',')
			}
			appendString(output, member.Name)
			output.WriteByte(':')
			if err := appendValue(output, member.Value); err != nil {
				return err
			}
		}
		output.WriteByte('}')
	default:
		return fmt.Errorf("unsupported canonical JSON value %T", value)
	}
	return nil
}

func appendString(output *bytes.Buffer, value string) {
	output.WriteByte('"')
	for _, character := range value {
		switch character {
		case '"', '\\':
			output.WriteByte('\\')
			output.WriteRune(character)
		case '\b':
			output.WriteString(`\b`)
		case '\t':
			output.WriteString(`\t`)
		case '\n':
			output.WriteString(`\n`)
		case '\f':
			output.WriteString(`\f`)
		case '\r':
			output.WriteString(`\r`)
		default:
			if character < 0x20 {
				fmt.Fprintf(output, `\u%04x`, character)
			} else {
				output.WriteRune(character)
			}
		}
	}
	output.WriteByte('"')
}

func compareUTF16(left, right string) int {
	a := utf16.Encode([]rune(left))
	b := utf16.Encode([]rune(right))
	for index := 0; index < len(a) && index < len(b); index++ {
		if a[index] < b[index] {
			return -1
		}
		if a[index] > b[index] {
			return 1
		}
	}
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	default:
		return 0
	}
}

func validateEscapedUnicode(input []byte) error {
	for index := 0; index < len(input); index++ {
		if input[index] != '"' {
			continue
		}
		index++
		for ; index < len(input) && input[index] != '"'; index++ {
			if input[index] < 0x20 {
				return errors.New("unescaped control character in JSON string")
			}
			if input[index] != '\\' {
				continue
			}
			index++
			if index >= len(input) {
				return errors.New("unterminated JSON escape")
			}
			if input[index] != 'u' {
				continue
			}
			unit, ok := parseHex16(input, index+1)
			if !ok {
				return errors.New("invalid Unicode escape")
			}
			index += 4
			switch {
			case unit >= 0xd800 && unit <= 0xdbff:
				if index+6 >= len(input) || input[index+1] != '\\' || input[index+2] != 'u' {
					return errors.New("lone high surrogate")
				}
				low, ok := parseHex16(input, index+3)
				if !ok || low < 0xdc00 || low > 0xdfff {
					return errors.New("invalid surrogate pair")
				}
				index += 6
			case unit >= 0xdc00 && unit <= 0xdfff:
				return errors.New("lone low surrogate")
			}
		}
	}
	return nil
}

func parseHex16(input []byte, start int) (uint16, bool) {
	if start+4 > len(input) {
		return 0, false
	}
	var value uint16
	for _, character := range input[start : start+4] {
		value <<= 4
		switch {
		case character >= '0' && character <= '9':
			value |= uint16(character - '0')
		case character >= 'a' && character <= 'f':
			value |= uint16(character-'a') + 10
		case character >= 'A' && character <= 'F':
			value |= uint16(character-'A') + 10
		default:
			return 0, false
		}
	}
	return value, true
}

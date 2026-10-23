package agentwork

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

// A focused wire-shape check keeps the published schema aligned with actual Go
// decoding/serialization. Semantic limits and lifecycle behavior are exercised
// by the executable service/CLI tests, without a validator dependency.
func TestAgentWorkSchemaFieldParity(t *testing.T) {
	b, e := os.ReadFile("../../schemas/agent-work-v1/contract.schema.json")
	if e != nil {
		t.Fatal(e)
	}
	var schema struct {
		Defs map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
			Required   []string                   `json:"required"`
		} `json:"$defs"`
	}
	if e = json.Unmarshal(b, &schema); e != nil {
		t.Fatal(e)
	}
	for name, v := range map[string]any{"RegisterRequest": Draft{}, "Mutation": Mutation{}, "Registrar": Registrar{}, "CycleHistory": CycleHistory{}, "Followup": Followup{}, "Work": Work{}, "Subject": Subject{}, "Capability": Capability{}, "Transport": Transport{}, "Recheck": Recheck{}, "Source": Source{}, "Agent": Agent{}, "Overview": Overview{}, "List": List{}} {
		typ := reflect.TypeOf(v)
		def, ok := schema.Defs[name]
		if !ok {
			t.Fatal("missing schema", name)
		}
		required := map[string]bool{}
		for _, k := range def.Required {
			required[k] = true
		}
		count := 0
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			tag := strings.Split(f.Tag.Get("json"), ",")
			if tag[0] == "-" {
				continue
			}
			count++
			if _, ok := def.Properties[tag[0]]; !ok {
				t.Errorf("%s.%s absent from schema", name, tag[0])
			}
			optional := len(tag) > 1 && tag[1] == "omitempty"
			if optional == required[tag[0]] {
				t.Errorf("%s.%s schema required mismatch", name, tag[0])
			}
		}
		if count != len(def.Properties) {
			t.Errorf("%s has schema-only fields", name)
		}
	}
}

package gitworktree

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseInventoryRecords(t *testing.T) {
	tests := []struct {
		name, raw string
		want      []Entry
	}{
		{"empty", "", []Entry{}},
		{"separators", "\x00\x00", []Entry{}},
		{"branch", "worktree /z path\nwith newline\x00HEAD opaque\x00branch refs/heads/main\x00locked reason with spaces\x00prunable reason\x00\x00", []Entry{{Locator: "/z path\nwith newline", OID: "opaque", Ref: "refs/heads/main", Locked: true, Prunable: true}}},
		{"bare without head", "worktree /bare\x00bare\x00\x00", []Entry{{Locator: "/bare", Bare: true}}},
		{"detached", "worktree /detached\x00HEAD abc\x00detached\x00\x00", []Entry{{Locator: "/detached", OID: "abc", Detached: true}}},
		{"wire order", "worktree /z\x00bare\x00\x00worktree /a\x00bare\x00\x00", []Entry{{Locator: "/z", Bare: true}, {Locator: "/a", Bare: true}}},
		{"adjacent starts and eof", "worktree /z\x00bare\x00worktree /a\x00bare", []Entry{{Locator: "/z", Bare: true}, {Locator: "/a", Bare: true}}},
		{"values belong to adapters", "worktree relative/../path\x00HEAD invalid\x00branch invalid", []Entry{{Locator: "relative/../path", OID: "invalid", Ref: "invalid"}}},
		{"missing values belong to adapters", "worktree /a\x00", []Entry{{Locator: "/a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInventory([]byte(tt.raw))
			if err != nil || !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, %v; want %#v", got, err, tt.want)
			}
		})
	}
}

func TestParseInventoryRejectsMalformedFields(t *testing.T) {
	bad := []string{"HEAD abc\x00", "unknown value\x00", "worktree \x00bare\x00", "worktree\x00", "worktree /a\x00unknown value\x00", "worktree /a\x00bare\x00\x00HEAD abc\x00"}
	for _, key := range []string{"HEAD", "branch", "locked", "prunable", "bare", "detached"} {
		bad = append(bad, key+" value\x00", "worktree /a\x00"+key+" value\x00"+key+" value\x00")
	}
	for _, raw := range bad {
		t.Run(strings.ReplaceAll(raw, "\x00", "|"), func(t *testing.T) {
			if got, err := ParseInventory([]byte(raw)); err == nil || got != nil {
				t.Fatalf("accepted %q: %#v, %v", raw, got, err)
			}
		})
	}
}

func TestParseInventoryOwnsResults(t *testing.T) {
	raw := []byte("worktree /a\x00HEAD abc\x00branch refs/heads/main\x00\x00")
	original := append([]byte(nil), raw...)
	first, err := ParseInventory(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(raw, original) {
		t.Fatal("input mutated")
	}
	second, err := ParseInventory(raw)
	if err != nil {
		t.Fatal(err)
	}
	for i := range raw {
		raw[i] = 'x'
	}
	want := Entry{Locator: "/a", OID: "abc", Ref: "refs/heads/main"}
	if first[0] != want || second[0] != want {
		t.Fatal("result aliases input")
	}
	first[0].Locator = "changed"
	if second[0] != want {
		t.Fatal("results share storage")
	}
}

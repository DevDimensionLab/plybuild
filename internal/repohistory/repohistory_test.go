package repohistory

import (
	"errors"
	"reflect"
	"testing"
)

func TestCommonDirs(t *testing.T) {
	for name, tc := range map[string]struct {
		in   []Binding
		want map[string]string
	}{
		"empty":            {nil, map[string]string{}},
		"same common dir":  {[]Binding{{"ply", "/r/ply/.git"}, {"ply", "/r/ply/.git"}, {"app", "/r/app/.git"}}, map[string]string{"ply": "/r/ply/.git", "app": "/r/app/.git"}},
		"empty is skipped": {[]Binding{{"ply", ""}, {"ply", "/r/ply/.git"}, {"Not valid", ""}, {"ply", ""}}, map[string]string{"ply": "/r/ply/.git"}},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := CommonDirs(tc.in)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v", got, err)
			}
		})
	}
}

func TestCommonDirsConflict(t *testing.T) {
	_, err := CommonDirs([]Binding{{"ply", "/r/ply/.git"}, {"app", "/r/app/.git"}, {"ply", "/r/other/.git"}})
	var conflict *ConflictError
	if !errors.As(err, &conflict) || conflict.RepoID != "ply" {
		t.Fatalf("error %v", err)
	}
	if err.Error() != "conflicting historical repository identity for ply" {
		t.Fatalf("message %q", err.Error())
	}
}
